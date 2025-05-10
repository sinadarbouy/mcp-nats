package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/nats-io/nats.go"
	"github.com/sinadarbouy/mcp-nats/tools"
)

func newServer(natsURL string, natsCredsPath string) *server.MCPServer {
	s := server.NewMCPServer(
		"mcp-nats",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Initialize NATS server tools
	natsTools := tools.NewNATSServerTools(natsURL, natsCredsPath)

	// Register all NATS server tools
	tools.RegisterTools(s, natsTools)

	return s
}

func run(transport, addr, natsURL, natsCredsPath string) error {
	s := newServer(natsURL, natsCredsPath)

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
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		log.Fatal("NATS_URL environment variable is required")
	}
	natsCredsPath := os.Getenv("NATS_CREDS_PATH")
	if natsCredsPath == "" {
		log.Fatal("NATS_CREDS_PATH environment variable is required")
	}

	// Connect to NATS
	nc, err := nats.Connect(natsURL, nats.UserCredentials(natsCredsPath))
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	fmt.Printf("Starting MCP NATS server (stdio)...\n")
	if err := run("stdio", "0.0.0.0:8002", natsURL, natsCredsPath); err != nil {
		panic(err)
	}
}
