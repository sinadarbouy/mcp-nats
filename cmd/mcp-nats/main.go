package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sinadarbouy/mcp-nats/tools"
)

func newServer(natsURL string, AccNatsCredsPath string, SysNatsCredsPath string) *server.MCPServer {
	s := server.NewMCPServer(
		"mcp-nats",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Initialize NATS server tools
	natsTools := tools.NewNATSServerTools(natsURL, AccNatsCredsPath, SysNatsCredsPath)

	// Register all NATS server tools
	tools.RegisterTools(s, natsTools)

	return s
}

func run(transport, addr, natsURL, AccNatsCredsPath, SysNatsCredsPath string) error {
	s := newServer(natsURL, AccNatsCredsPath, SysNatsCredsPath)

	switch transport {
	case "stdio":
		srv := server.NewStdioServer(s)
		return srv.Listen(context.Background(), os.Stdin, os.Stdout)
	case "sse":
		srv := server.NewSSEServer(s)
		slog.Info("Starting NATS MCP server using SSE transport", "address", addr)
		if err := srv.Start(addr); err != nil {
			return fmt.Errorf("Server error: %v", err)
		}
	default:
		return fmt.Errorf(
			"Invalid transport type: %s. Must be 'stdio' or 'sse'",
			transport,
		)
	}
	return nil
}

func main() {
	// Parse command line flags
	transport := flag.String("transport", "stdio", "Transport type (stdio or sse)")
	sseAddr := flag.String("sse-address", "0.0.0.0:8000", "Address for SSE server to listen on")
	flag.Parse()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL environment variable is required")
	}
	AccNatsCredsPath := os.Getenv("NATS_CREDS_PATH")
	if AccNatsCredsPath == "" {
		log.Fatal("NATS_CREDS_PATH environment variable is required")
	}
	SysNatsCredsPath := os.Getenv("NATS_CREDS_PATH_SYS")
	if SysNatsCredsPath == "" {
		log.Fatal("NATS_CREDS_PATH_SYS environment variable is required")
	}

	fmt.Printf("Starting MCP NATS server (%s)...\n", *transport)
	if err := run(*transport, *sseAddr, natsURL, AccNatsCredsPath, SysNatsCredsPath); err != nil {
		panic(err)
	}
}
