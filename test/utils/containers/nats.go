package containers

import (
	"context"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NatsContainer represents a Nats test container
type NatsContainer struct {
	Container testcontainers.Container
	Port      nat.Port
	Host      string
}

// NewNatsContainer creates and starts a new Nats container for testing
func NewNatsContainer(ctx context.Context, t *testing.T) *NatsContainer {
	t.Helper()
	NatsPort := "4222/tcp"
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "nats:2.10.12",
			ExposedPorts: []string{NatsPort},
			WaitingFor: wait.ForAll(
				wait.ForLog("Server is ready"),
				wait.ForListeningPort(nat.Port(NatsPort)),
			),
			AutoRemove: true,
			Cmd:        []string{"-js"}, // Enable JetStream
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to create Nats container: %v", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(NatsPort))
	if err != nil {
		t.Fatalf("failed to get mapped port: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	return &NatsContainer{
		Container: container,
		Port:      mappedPort,
		Host:      host,
	}
}
