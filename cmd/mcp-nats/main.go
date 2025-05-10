package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sinadarbouy/mcp-nats/internal/logger"
	"github.com/sinadarbouy/mcp-nats/tools"
)

func newServer(natsURL string) (*server.MCPServer, error) {
	s := server.NewMCPServer(
		"mcp-nats",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Initialize NATS server tools
	natsTools, err := tools.NewNATSServerTools(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NATS tools: %v", err)
	}

	// Register all NATS server tools
	tools.RegisterTools(s, natsTools)

	return s, nil
}

func run(transport, addr, natsURL string) error {
	s, err := newServer(natsURL)
	if err != nil {
		return err
	}

	switch transport {
	case "stdio":
		srv := server.NewStdioServer(s)
		logger.Info("Starting NATS MCP server using stdio transport")
		return srv.Listen(context.Background(), os.Stdin, os.Stdout)
	case "sse":
		srv := server.NewSSEServer(s)
		logger.Info("Starting NATS MCP server using SSE transport",
			"address", addr,
		)
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
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	jsonLogs := flag.Bool("json-logs", false, "Output logs in JSON format")
	flag.Parse()

	// Initialize logger
	logger.Initialize(logger.Config{
		Level:      logger.GetLevel(*logLevel),
		JSONFormat: *jsonLogs,
	})

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		logger.Error("NATS_URL environment variable is required")
		os.Exit(1)
	}

	logger.Info("Starting MCP NATS server",
		"transport", *transport,
		"version", "0.1.0",
	)

	if err := run(*transport, *sseAddr, natsURL); err != nil {
		logger.Error("Server failed",
			"error", err,
		)
		os.Exit(1)
	}
}
