package tools

import (
	"context"
	"fmt"

	mcpnats "github.com/sinadarbouy/mcp-nats"
	"github.com/sinadarbouy/mcp-nats/internal/logger"
	"github.com/sinadarbouy/mcp-nats/tools/common"
)

// NATSServerTools contains all NATS server-related tool definitions
type NATSServerTools struct {
	executors    map[string]*common.NATSExecutor
	serverTools  *ServerTools
	streamTools  *StreamTools
	kvTools      *KVTools
	publishTools *PublishTools
	accountTools *AccountTools
	rttTools     *RTTTools
}

// NewNATSServerTools creates a new instance of NATSServerTools
func NewNATSServerTools() (*NATSServerTools, error) {
	n := &NATSServerTools{
		executors: make(map[string]*common.NATSExecutor),
	}

	// Initialize tool categories
	n.serverTools = NewServerTools(n)
	n.streamTools = NewStreamTools(n)
	n.kvTools = NewKVTools(n)
	n.publishTools = NewPublishTools(n)
	n.accountTools = NewAccountTools(n)
	n.rttTools = NewRTTTools(n)

	logger.Info("Initialized NATS server tools")

	return n, nil
}

// GetExecutor returns the executor for the specified account
func (n *NATSServerTools) GetExecutor(ctx context.Context, accountName string) (*common.NATSExecutor, error) {
	// Try to get existing executor
	if executor, ok := n.executors[accountName]; ok {
		return executor, nil
	}

	// Get credentials from context
	creds, err := mcpnats.GetCredsFromContext(ctx, accountName)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials for account %s: %v", accountName, err)
	}

	natsURL, err := mcpnats.GetNatsURLFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get NATS URL: %w", err)
	}

	// Create new executor
	executor, err := common.NewNATSExecutor(natsURL, creds)
	if err != nil {
		return nil, fmt.Errorf("failed to create executor for account %s: %v", accountName, err)
	}

	// Cache the executor
	n.executors[accountName] = executor
	return executor, nil
}

// Cleanup removes all temporary credential files
func (n *NATSServerTools) Cleanup() {
	for _, executor := range n.executors {
		if err := executor.Cleanup(); err != nil {
			logger.Error("Failed to cleanup executor",
				"error", err,
				"account", executor.Creds.AccountName,
			)
		}
	}
}

// ServerTools returns the server tools category
func (n *NATSServerTools) ServerTools() ToolCategory {
	return n.serverTools
}

// StreamTools returns the stream tools category
func (n *NATSServerTools) StreamTools() ToolCategory {
	return n.streamTools
}

// KVTools returns the KV tools category
func (n *NATSServerTools) KVTools() ToolCategory {
	return n.kvTools
}

// PublishTools returns the publish tools category
func (n *NATSServerTools) PublishTools() ToolCategory {
	return n.publishTools
}

// AccountTools returns the account tools category
func (n *NATSServerTools) AccountTools() ToolCategory {
	return n.accountTools
}

// RTTTools returns the RTT tools category
func (n *NATSServerTools) RTTTools() ToolCategory {
	return n.rttTools
}
