package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
	mcpnats "github.com/sinadarbouy/mcp-nats"
	"github.com/sinadarbouy/mcp-nats/internal/logger"
	"github.com/sinadarbouy/mcp-nats/tools"
)

func newServer() (*server.MCPServer, error) {
	s := server.NewMCPServer(
		"mcp-nats",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Initialize NATS server tools
	natsTools, err := tools.NewNATSServerTools()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NATS tools: %v", err)
	}

	// Register all NATS server tools
	tools.RegisterTools(s, natsTools)

	return s, nil
}

func run(transport, addr string) error {
	s, err := newServer()
	if err != nil {
		return err
	}

	switch transport {
	case "stdio":
		srv := server.NewStdioServer(s)
		srv.SetContextFunc(mcpnats.ComposedStdioContextFunc())
		logger.Info("Starting NATS MCP server using stdio transport")
		return srv.Listen(context.Background(), os.Stdin, os.Stdout)
	case "sse":
		srv := server.NewSSEServer(s, server.WithSSEContextFunc(mcpnats.ComposedSSEContextFunc()))
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

	logger.Info("Starting MCP NATS server",
		"transport", *transport,
		"version", "0.1.0",
	)

	if err := run(*transport, *sseAddr); err != nil {
		logger.Error("Server failed",
			"error", err,
		)
		os.Exit(1)
	}
}
