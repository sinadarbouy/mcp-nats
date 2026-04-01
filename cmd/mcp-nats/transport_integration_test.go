package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sinadarbouy/mcp-nats/test/utils/containers"
)

func buildMCPBinary(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	cmdDir := filepath.Dir(file)
	out := filepath.Join(t.TempDir(), "mcp-nats"+exeSuffix())
	build := exec.Command("go", "build", "-o", out, ".")
	build.Dir = cmdDir
	build.Env = os.Environ()
	if outBytes, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, outBytes)
	}
	return out
}

func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func natsURLFromContainer(c *containers.NatsContainer) string {
	// c.Port is docker/nat.Port; String() includes "/tcp" which breaks nats URLs.
	return fmt.Sprintf("nats://%s:%s", c.Host, c.Port.Port())
}

func waitHTTPStatus(t *testing.T, url string, want int, d time.Duration) {
	t.Helper()
	deadline := time.Now().Add(d)
	var lastErr error
	for time.Now().Before(deadline) {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if resp.StatusCode == want {
			return
		}
		lastErr = fmt.Errorf("status %d", resp.StatusCode)
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("GET %s: want status %d: %v", url, want, lastErr)
}

func mcpInitRequest() mcp.InitializeRequest {
	req := mcp.InitializeRequest{}
	req.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	req.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp-nats-transport-integration",
		Version: "1.0.0",
	}
	req.Params.Capabilities = mcp.ClientCapabilities{}
	return req
}

func assertMCPSmoke(ctx context.Context, t *testing.T, c *client.Client) {
	t.Helper()
	if _, err := c.Initialize(ctx, mcpInitRequest()); err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	tools, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	if tools == nil || len(tools.Tools) == 0 {
		t.Fatal("expected non-empty tool list")
	}
	// Server tools use $SYS and need a system account; publish works with anonymous auth.
	callReq := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "publish",
			Arguments: map[string]any{
				"subject": "mcp.integration.smoke",
				"body":    "ok",
				"count":   1,
			},
		},
	}
	res, err := c.CallTool(ctx, callReq)
	if err != nil {
		t.Fatalf("CallTool publish: %v", err)
	}
	if res == nil || res.IsError {
		t.Fatalf("publish: IsError=%v result=%+v", res != nil && res.IsError, res)
	}
	if len(res.Content) == 0 {
		t.Fatal("publish: empty content")
	}
}

// TestMCPTransports_Integration runs MCP initialize, tools/list, and server_list over
// streamable-http, SSE, and stdio against a real NATS testcontainer and a subprocess
// mcp-nats binary (requires Docker).
func TestMCPTransports_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	natsC := containers.NewNatsContainer(ctx, t)
	t.Cleanup(func() {
		_ = natsC.Container.Terminate(context.Background())
	})
	nURL := natsURLFromContainer(natsC)
	bin := buildMCPBinary(t)

	t.Run("streamable-http", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
		defer cancel()

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		addr := ln.Addr().String()
		_ = ln.Close()

		cmdCtx, cmdCancel := context.WithCancel(ctx)
		defer cmdCancel()
		cmd := exec.CommandContext(cmdCtx, bin,
			"--transport", "streamable-http",
			"--address", addr,
			"--endpoint-path", "/mcp",
			"--log-level", "error",
		)
		cmd.Env = append(os.Environ(),
			"NATS_URL="+nURL,
			"NATS_NO_AUTHENTICATION=true",
		)
		var stderr []byte
		cmd.Stderr = &stderrWriter{&stderr}
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			cmdCancel()
			_ = cmd.Wait()
			if t.Failed() && len(stderr) > 0 {
				t.Logf("mcp-nats stderr:\n%s", stderr)
			}
		})

		base := "http://" + addr
		waitHTTPStatus(t, base+"/livez", http.StatusOK, 30*time.Second)
		waitHTTPStatus(t, base+"/readyz", http.StatusOK, 30*time.Second)

		httpClient, err := client.NewStreamableHttpClient(base + "/mcp")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = httpClient.Close() }()
		if err := httpClient.Start(ctx); err != nil {
			t.Fatalf("client Start: %v", err)
		}
		assertMCPSmoke(ctx, t, httpClient)
	})

	t.Run("sse", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
		defer cancel()

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		addr := ln.Addr().String()
		_ = ln.Close()

		cmdCtx, cmdCancel := context.WithCancel(ctx)
		defer cmdCancel()
		cmd := exec.CommandContext(cmdCtx, bin,
			"--transport", "sse",
			"--address", addr,
			"--log-level", "error",
		)
		cmd.Env = append(os.Environ(),
			"NATS_URL="+nURL,
			"NATS_NO_AUTHENTICATION=true",
		)
		var stderr []byte
		cmd.Stderr = &stderrWriter{&stderr}
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			cmdCancel()
			_ = cmd.Wait()
			if t.Failed() && len(stderr) > 0 {
				t.Logf("mcp-nats stderr:\n%s", stderr)
			}
		})

		base := "http://" + addr
		waitHTTPStatus(t, base+"/livez", http.StatusOK, 30*time.Second)
		waitHTTPStatus(t, base+"/readyz", http.StatusOK, 30*time.Second)

		sseClient, err := client.NewSSEMCPClient(base + "/sse")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = sseClient.Close() }()
		if err := sseClient.Start(ctx); err != nil {
			t.Fatalf("client Start: %v", err)
		}
		assertMCPSmoke(ctx, t, sseClient)
	})

	t.Run("stdio", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
		defer cancel()

		env := append(os.Environ(),
			"NATS_URL="+nURL,
			"NATS_NO_AUTHENTICATION=true",
		)
		stdioClient, err := client.NewStdioMCPClient(bin, env,
			"--transport", "stdio",
			"--log-level", "error",
		)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = stdioClient.Close() }()
		if err := stdioClient.Start(ctx); err != nil {
			t.Fatalf("client Start: %v", err)
		}
		assertMCPSmoke(ctx, t, stdioClient)
	})
}

type stderrWriter struct {
	b *[]byte
}

func (w *stderrWriter) Write(p []byte) (int, error) {
	*w.b = append(*w.b, p...)
	return len(p), nil
}
