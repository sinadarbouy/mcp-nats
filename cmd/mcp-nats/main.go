package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
	mcpnats "github.com/sinadarbouy/mcp-nats"
	"github.com/sinadarbouy/mcp-nats/internal/logger"
	"github.com/sinadarbouy/mcp-nats/tools"
)

const (
	// Version of the application
	Version = "0.1.0"
	// AppName is the name of the application
	AppName = "mcp-nats"
)

// Config holds all configuration for the server
type Config struct {
	Transport        string
	SSEAddr          string
	LogLevel         string
	JSONLogs         bool
	NoAuthentication bool
	NATSUser         string
	NATSPassword     string
}

// validateConfig ensures all config values are valid
func validateConfig(cfg *Config) error {
	if cfg.Transport != "stdio" && cfg.Transport != "sse" {
		return fmt.Errorf("invalid transport type: %s (must be 'stdio' or 'sse')", cfg.Transport)
	}
	if cfg.Transport == "sse" && cfg.SSEAddr == "" {
		return fmt.Errorf("sse-address cannot be empty when using sse transport")
	}
	return nil
}

func newServer() (*server.MCPServer, error) {
	s := server.NewMCPServer(
		AppName,
		Version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Initialize NATS server tools
	natsTools, err := tools.NewNATSServerTools()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NATS tools: %w", err)
	}

	// Register all NATS server tools
	tools.RegisterTools(s, natsTools)

	return s, nil
}

func run(ctx context.Context, cfg *Config) error {
	// Set environment variables for authentication configuration
	if cfg.NoAuthentication {
		if err := os.Setenv("NATS_NO_AUTHENTICATION", "true"); err != nil {
			return fmt.Errorf("failed to set NATS_NO_AUTHENTICATION env var: %w", err)
		}
	}
	if cfg.NATSUser != "" {
		if err := os.Setenv("NATS_USER", cfg.NATSUser); err != nil {
			return fmt.Errorf("failed to set NATS_USER env var: %w", err)
		}
	}
	if cfg.NATSPassword != "" {
		if err := os.Setenv("NATS_PASSWORD", cfg.NATSPassword); err != nil {
			return fmt.Errorf("failed to set NATS_PASSWORD env var: %w", err)
		}
	}

	s, err := newServer()
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	switch cfg.Transport {
	case "stdio":
		srv := server.NewStdioServer(s)
		srv.SetContextFunc(mcpnats.ComposedStdioContextFunc())
		logger.Info("Starting NATS MCP server using stdio transport")
		return srv.Listen(ctx, os.Stdin, os.Stdout)

	case "sse":
		srv := server.NewSSEServer(s, server.WithSSEContextFunc(mcpnats.ComposedSSEContextFunc()))
		logger.Info("Starting NATS MCP server using SSE transport",
			"address", cfg.SSEAddr,
		)

		errChan := make(chan error, 1)
		go func() {
			if err := srv.Start(cfg.SSEAddr); err != nil {
				errChan <- fmt.Errorf("server error: %w", err)
			}
		}()

		// Wait for either context cancellation or server error
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			// Give the server some time to shutdown gracefully
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return srv.Shutdown(shutdownCtx)
		}
	}

	return nil
}

func main() {
	cfg := &Config{}

	// Parse command line flags
	flag.StringVar(&cfg.Transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&cfg.SSEAddr, "sse-address", "0.0.0.0:8000", "Address for SSE server to listen on")
	flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.BoolVar(&cfg.JSONLogs, "json-logs", false, "Output logs in JSON format")
	flag.BoolVar(&cfg.NoAuthentication, "no-authentication", false, "Allow anonymous connections without credentials")
	flag.StringVar(&cfg.NATSUser, "user", "", "NATS username or token (can also be set via NATS_USER env var)")
	flag.StringVar(&cfg.NATSPassword, "password", "", "NATS password (can also be set via NATS_PASSWORD env var)")
	flag.Parse()

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Initialize(logger.Config{
		Level:      logger.GetLevel(cfg.LogLevel),
		JSONFormat: cfg.JSONLogs,
	})

	logger.Info("Starting MCP NATS server",
		"transport", cfg.Transport,
		"version", Version,
	)

	// Setup context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal", "signal", sig)
		cancel()
	}()

	if err := run(ctx, cfg); err != nil {
		logger.Error("Server failed",
			"error", err,
		)
		os.Exit(1)
	}
}
