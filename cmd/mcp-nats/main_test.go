package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/sinadarbouy/mcp-nats/internal/logger"
)

func TestMain(m *testing.M) {
	logger.Initialize(logger.Config{Level: logger.LevelError})
	os.Exit(m.Run())
}

func TestHandleLivez(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/livez", nil)
	rec := httptest.NewRecorder()

	handleLivez(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "ok" {
		t.Fatalf("expected body ok, got %q", rec.Body.String())
	}
}

func TestHandleReadyzUnavailable(t *testing.T) {
	t.Setenv("NATS_URL", "nats://127.0.0.1:1")

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	handleReadyz(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, rec.Code)
	}
}

func TestCheckNATSConnectivityReachable(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer func() { _ = listener.Close() }()

	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr == nil {
			_ = conn.Close()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := checkNATSConnectivity(ctx, listener.Addr().String()); err != nil {
		t.Fatalf("expected reachable listener to pass readiness check: %v", err)
	}
}
