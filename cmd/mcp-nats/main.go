package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
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
	Address          string
	EndpointPath     string
	LogLevel         string
	JSONLogs         bool
	NoAuthentication bool
	NATSUser         string
	NATSPassword     string
}

// validateConfig ensures all config values are valid
func validateConfig(cfg *Config) error {
	if cfg.Transport != "stdio" && cfg.Transport != "sse" && cfg.Transport != "streamable-http" {
		return fmt.Errorf("invalid transport type: %s (must be 'stdio', 'sse' or 'streamable-http')", cfg.Transport)
	}
	if (cfg.Transport == "sse" || cfg.Transport == "streamable-http") && cfg.Address == "" {
		return fmt.Errorf("address cannot be empty when using %s transport", cfg.Transport)
	}
	if cfg.Transport == "streamable-http" {
		if cfg.EndpointPath == "" {
			return fmt.Errorf("endpoint-path cannot be empty when using streamable-http transport")
		}
		if !strings.HasPrefix(cfg.EndpointPath, "/") {
			return fmt.Errorf("endpoint-path must start with '/'")
		}
	}
	return nil
}

type httpServer interface {
	Start(addr string) error
	Shutdown(ctx context.Context) error
}

const readinessTimeout = 2 * time.Second

func writeOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func handleLivez(w http.ResponseWriter, _ *http.Request) {
	writeOK(w)
}

func handleReadyz(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), readinessTimeout)
	defer cancel()

	if err := checkNATSConnectivity(ctx, os.Getenv("NATS_URL")); err != nil {
		logger.Warn("Readiness check failed", "error", err)
		http.Error(w, "nats unavailable", http.StatusServiceUnavailable)
		return
	}

	writeOK(w)
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	// Keep /healthz as a stable alias for liveness checks.
	handleLivez(w, nil)
}

func checkNATSConnectivity(ctx context.Context, rawURL string) error {
	natsURL := strings.TrimSpace(rawURL)
	if natsURL == "" {
		natsURL = "localhost:4222"
	}

	address := natsURL
	if strings.Contains(natsURL, "://") {
		parsedURL, err := url.Parse(natsURL)
		if err != nil {
			return fmt.Errorf("invalid NATS URL %q: %w", natsURL, err)
		}
		if parsedURL.Host == "" {
			return fmt.Errorf("invalid NATS URL %q: missing host", natsURL)
		}
		address = parsedURL.Host
	}

	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("dial NATS at %s: %w", address, err)
	}

	if closeErr := conn.Close(); closeErr != nil {
		logger.Debug("Failed to close readiness probe connection", "error", closeErr)
	}
	return nil
}

func newHTTPMux(mcpPath string, mcpHandler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(mcpPath, mcpHandler)
	mux.HandleFunc("/livez", handleLivez)
	mux.HandleFunc("/readyz", handleReadyz)
	mux.HandleFunc("/healthz", handleHealthz)
	return mux
}

func newSSEHTTPMux(sseSrv *server.SSEServer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/sse", sseSrv.SSEHandler())
	mux.Handle("/message", sseSrv.MessageHandler())
	mux.HandleFunc("/livez", handleLivez)
	mux.HandleFunc("/readyz", handleReadyz)
	mux.HandleFunc("/healthz", handleHealthz)
	return mux
}

func runHTTPServer(ctx context.Context, srv httpServer, addr, transportName string) error {
	serverErr := make(chan error, 1)
	go func() {
		if err := srv.Start(addr); err != nil {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logger.Info("Shutting down HTTP server", "transport", transportName)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown %s server: %w", transportName, err)
		}
		if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
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
		httpSrv := &http.Server{Addr: cfg.Address}
		srv := server.NewSSEServer(
			s,
			server.WithSSEContextFunc(mcpnats.ComposedSSEContextFunc()),
			server.WithHTTPServer(httpSrv),
		)
		httpSrv.Handler = newSSEHTTPMux(srv)
		logger.Info("Starting NATS MCP server using SSE transport",
			"address", cfg.Address,
		)
		return runHTTPServer(ctx, srv, cfg.Address, "sse")

	case "streamable-http":
		httpSrv := &http.Server{Addr: cfg.Address}
		srv := server.NewStreamableHTTPServer(s,
			server.WithHTTPContextFunc(server.HTTPContextFunc(mcpnats.ComposedSSEContextFunc())),
			server.WithEndpointPath(cfg.EndpointPath),
			server.WithStreamableHTTPServer(httpSrv),
		)
		httpSrv.Handler = newHTTPMux(cfg.EndpointPath, srv)

		logger.Info("Starting NATS MCP server using Streamable HTTP transport",
			"address", cfg.Address,
			"endpointPath", cfg.EndpointPath,
		)
		return runHTTPServer(ctx, srv, cfg.Address, "streamable-http")
	}

	return nil
}

func main() {
	cfg := &Config{}

	// Parse command line flags
	flag.StringVar(&cfg.Transport, "transport", "streamable-http", "Transport type (stdio, sse or streamable-http)")
	flag.StringVar(&cfg.Address, "address", "0.0.0.0:8000", "Address for HTTP server to listen on")
	flag.StringVar(&cfg.Address, "sse-address", "0.0.0.0:8000", "Deprecated: use --address instead")
	flag.StringVar(&cfg.EndpointPath, "endpoint-path", "/mcp", "Endpoint path for streamable-http server")
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
