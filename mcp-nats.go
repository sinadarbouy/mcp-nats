// Package mcpnats provides NATS integration for MCP (Model Context Protocol).
// It handles NATS URL configuration, credentials management, and context propagation
// for both SSE (Server-Sent Events) and stdio-based communication. This package
// facilitates the communication layer for the Model Context Protocol over NATS,
// providing a robust and scalable messaging infrastructure.
package mcpnats

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sinadarbouy/mcp-nats/tools/common"
)

const (
	defaultNatsURL = "localhost:4222"

	natsURLEnvVar = "NATS_URL"

	natsURLHeader = "X-Nats-URL"
)

type natsURLKey struct{}
type natsCredsKey struct{}

func urlFromHeaders(req *http.Request) (string, string) {
	if req == nil {
		return "", ""
	}
	slog.Debug("Processing NATS headers", "headers", req.Header)
	return req.Header.Get(natsURLHeader), ""
}

func urlFromEnv() (string, string) {
	slog.Debug("Processing NATS environment variables")
	u := strings.TrimRight(os.Getenv(natsURLEnvVar), "/")
	return u, ""
}

// WithNatsURL adds the Nats URL to the context.
func WithNatsURL(ctx context.Context, url string) context.Context {
	return context.WithValue(ctx, natsURLKey{}, url)
}

// WithNatsCreds adds NATS credentials to the context
func WithNatsCreds(ctx context.Context, creds map[string]common.NATSCreds) context.Context {
	return context.WithValue(ctx, natsCredsKey{}, creds)
}

// natsURLFromContext extracts the nats url from the context.
// This can be used by tools to extract the url regardless of the
// transport being used by the server.
func natsURLFromContext(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", fmt.Errorf("context is nil")
	}

	url, ok := ctx.Value(natsURLKey{}).(string)
	if !ok {
		return "", fmt.Errorf("nats url not found in context")
	}
	return url, nil
}

// natsCredsFromContext extracts the NATS credentials from the context
func natsCredsFromContext(ctx context.Context) (map[string]common.NATSCreds, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	creds, ok := ctx.Value(natsCredsKey{}).(map[string]common.NATSCreds)
	if !ok {
		return nil, fmt.Errorf("nats credentials not found in context")
	}
	return creds, nil
}

// ExtractNatsInfoFromHeaders is a SSEContextFunc that extracts nats configuration
// from request headers and injects a configured client into the context.
var ExtractNatsInfoFromHeaders server.SSEContextFunc = func(ctx context.Context, req *http.Request) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	u, _ := urlFromHeaders(req)
	uEnv, _ := urlFromEnv()

	// Determine final URL with fallbacks
	u = determineNatsURL(u, uEnv)

	// Validate URL
	if err := validateNatsURL(u); err != nil {
		slog.Error("Invalid NATS URL", "url", u, "error", err)
		// Use default URL as fallback
		u = defaultNatsURL
	}

	// Get NATS credentials
	creds, err := common.GetCredsFromEnv()
	if err != nil {
		slog.Error("Failed to get NATS credentials", "error", err)
		creds = make(map[string]common.NATSCreds)
	}

	return WithNatsCreds(WithNatsURL(ctx, u), creds)
}

// ExtractNatsInfoFromEnv is a StdioContextFunc that extracts NATS configuration
// from environment variables and injects a configured client into the context.
var ExtractNatsInfoFromEnv server.StdioContextFunc = func(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	u, _ := urlFromEnv()
	if u == "" {
		u = defaultNatsURL
	}

	if err := validateNatsURL(u); err != nil {
		slog.Error("Invalid NATS URL from environment", "url", u, "error", err)
		u = defaultNatsURL
	}

	// Get NATS credentials
	creds, err := common.GetCredsFromEnv()
	if err != nil {
		slog.Error("Failed to get NATS credentials from environment", "error", err)
		creds = make(map[string]common.NATSCreds)
	}

	return WithNatsCreds(WithNatsURL(ctx, u), creds)
}

// ComposeSSEContextFuncs composes multiple SSEContextFuncs into a single function.
// This allows for chaining multiple context modifiers together in a clean way.
func ComposeSSEContextFuncs(funcs ...server.SSEContextFunc) server.SSEContextFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		for _, f := range funcs {
			ctx = f(ctx, req)
		}
		return ctx
	}
}

// ComposedSSEContextFunc returns a composed SSEContextFunc that includes all
// predefined context functions for SSE handling. Currently, this includes
// ExtractNatsInfoFromHeaders.
func ComposedSSEContextFunc() server.SSEContextFunc {
	return ComposeSSEContextFuncs(
		ExtractNatsInfoFromHeaders,
	)
}

// ComposeStdioContextFuncs composes multiple StdioContextFuncs into a single function.
// This allows for chaining multiple context modifiers together in a clean way.
func ComposeStdioContextFuncs(funcs ...server.StdioContextFunc) server.StdioContextFunc {
	return func(ctx context.Context) context.Context {
		for _, f := range funcs {
			ctx = f(ctx)
		}
		return ctx
	}
}

// ComposedStdioContextFunc returns a composed StdioContextFunc that includes all
// predefined context functions for stdio handling. Currently, this includes
// ExtractNatsInfoFromEnv.
func ComposedStdioContextFunc() server.StdioContextFunc {
	return ComposeStdioContextFuncs(
		ExtractNatsInfoFromEnv,
	)
}

// GetCredsFromContext retrieves NATS credentials for a specific account from the context.
// It returns an error if:
// - The NATS credentials are not found in the context
// - No credentials are found for the specified account
func GetCredsFromContext(ctx context.Context, accountName string) (common.NATSCreds, error) {

	creds, err := natsCredsFromContext(ctx)
	if err != nil {
		return common.NATSCreds{}, fmt.Errorf("failed to get NATS credentials: %w", err)
	}

	if cred, ok := creds[accountName]; ok {
		return cred, nil
	}

	return common.NATSCreds{}, fmt.Errorf("no credentials found for account %s", accountName)
}

// GetNatsURLFromContext retrieves the NATS URL from the context.
// It returns an error if:
// - The NATS URL is not found in the context
// - The NATS URL is not valid
func GetNatsURLFromContext(ctx context.Context) (string, error) {
	url, err := natsURLFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get NATS URL: %w", err)
	}
	return url, nil
}

// Helper functions

// determineNatsURL returns the appropriate NATS URL based on the provided values
func determineNatsURL(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	if fallback != "" {
		return fallback
	}
	return defaultNatsURL
}

// validateNatsURL ensures the provided URL is valid for NATS
func validateNatsURL(u string) error {
	if u == "" {
		return fmt.Errorf("empty NATS URL")
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("invalid NATS URL format: %w", err)
	}

	// Basic validation for NATS URL
	if parsedURL.Scheme != "" && parsedURL.Scheme != "nats" {
		return fmt.Errorf("invalid NATS URL scheme: %s", parsedURL.Scheme)
	}

	return nil
}
