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
	// print all headers
	for k, v := range req.Header {
		fmt.Println(k, v)
	}
	u := req.Header.Get(natsURLHeader)
	return u, ""
}

func urlFromEnv() (string, string) {
	// print all env vars
	for _, e := range os.Environ() {
		fmt.Println(e)
	}
	u := strings.TrimRight(os.Getenv(natsURLEnvVar), "/")
	return u, ""
}

// WithNatsURL adds the Grafana URL to the context.
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
	url, ok := ctx.Value(natsURLKey{}).(string)
	if !ok {
		return "", fmt.Errorf("missing nats url")
	}
	return url, nil
}

// natsCredsFromContext extracts the NATS credentials from the context
func natsCredsFromContext(ctx context.Context) (map[string]common.NATSCreds, error) {
	creds, ok := ctx.Value(natsCredsKey{}).(map[string]common.NATSCreds)
	if !ok {
		return nil, fmt.Errorf("missing nats credentials")
	}
	return creds, nil
}

// ExtractNatsInfoFromHeaders is a SSEContextFunc that extracts Grafana configuration
// from request headers and injects a configured client into the context.
var ExtractNatsInfoFromHeaders server.SSEContextFunc = func(ctx context.Context, req *http.Request) context.Context {
	u, apiKey := urlFromHeaders(req)
	uEnv, apiKeyEnv := urlFromEnv()
	if u == "" {
		u = uEnv
	}
	if u == "" {
		u = defaultNatsURL
	}
	if apiKey == "" {
		apiKey = apiKeyEnv
	}

	// Get NATS credentials from environment
	creds, err := common.GetCredsFromEnv()
	if err != nil {
		// Log error or handle it appropriately
		creds = make(map[string]common.NATSCreds)
	}

	ctx = WithNatsURL(ctx, u)
	return WithNatsCreds(ctx, creds)
}

// ExtractNatsInfoFromEnv is a StdioContextFunc that extracts NATS configuration
// from environment variables and injects a configured client into the context.
var ExtractNatsInfoFromEnv server.StdioContextFunc = func(ctx context.Context) context.Context {
	u, apiKey := urlFromEnv()
	if u == "" {
		u = defaultNatsURL
	}
	parsedURL, err := url.Parse(u)
	if err != nil {
		panic(fmt.Errorf("invalid NATS URL %s: %w", u, err))
	}
	// Get NATS credentials from environment
	creds, err := common.GetCredsFromEnv()
	if err != nil {
		// Log error or handle it appropriately
		creds = make(map[string]common.NATSCreds)
	}
	slog.Info("Using NATS configuration", "url", parsedURL.Redacted(), "api_key_set", apiKey != "")
	return WithNatsURL(WithNatsCreds(ctx, creds), u)
}

// ComposeSSEContextFuncs composes multiple SSEContextFuncs into a single one.
func ComposeSSEContextFuncs(funcs ...server.SSEContextFunc) server.SSEContextFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		for _, f := range funcs {
			ctx = f(ctx, req)
		}
		return ctx
	}
}

// ComposedSSEContextFunc is a SSEContextFunc that comprises all predefined SSEContextFuncs.
func ComposedSSEContextFunc() server.SSEContextFunc {
	return ComposeSSEContextFuncs(
		ExtractNatsInfoFromHeaders,
	)
}

// ComposeStdioContextFuncs composes multiple StdioContextFuncs into a single one.
func ComposeStdioContextFuncs(funcs ...server.StdioContextFunc) server.StdioContextFunc {
	return func(ctx context.Context) context.Context {
		for _, f := range funcs {
			ctx = f(ctx)
		}
		return ctx
	}
}

// ComposedStdioContextFunc is a StdioContextFunc that comprises all predefined StdioContextFuncs.
func ComposedStdioContextFunc() server.StdioContextFunc {
	return ComposeStdioContextFuncs(
		ExtractNatsInfoFromEnv,
	)
}

// GetCredsFromContext gets credentials for a specific account from the context
func GetCredsFromContext(ctx context.Context, accountName string) (common.NATSCreds, error) {
	_, err := natsURLFromContext(ctx)
	if err != nil {
		return common.NATSCreds{}, err
	}

	creds, err := natsCredsFromContext(ctx)
	if err != nil {
		return common.NATSCreds{}, err
	}

	if cred, ok := creds[accountName]; ok {
		return cred, nil
	}

	return common.NATSCreds{}, fmt.Errorf("no credentials found for account %s", accountName)
}
