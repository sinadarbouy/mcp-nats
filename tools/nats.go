package tools

import (
	"fmt"

	"github.com/sinadarbouy/mcp-nats/internal/logger"
	"github.com/sinadarbouy/mcp-nats/tools/common"
)

// NATSServerTools contains all NATS server-related tool definitions
type NATSServerTools struct {
	executors   map[string]*common.NATSExecutor
	serverTools *ServerTools
	streamTools *StreamTools
}

// NewNATSServerTools creates a new instance of NATSServerTools
func NewNATSServerTools(natsURL string) (*NATSServerTools, error) {
	// Get credentials from environment variables
	creds, err := common.GetCredsFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials from env: %v", err)
	}

	if len(creds) == 0 {
		return nil, fmt.Errorf("no NATS credentials found in environment variables")
	}

	// Create executors for each account
	executors := make(map[string]*common.NATSExecutor)
	for _, cred := range creds {
		executor, err := common.NewNATSExecutor(natsURL, cred)
		if err != nil {
			// Cleanup already created executors
			for _, e := range executors {
				if cleanupErr := e.Cleanup(); cleanupErr != nil {
					logger.Error("Failed to cleanup executor",
						"error", cleanupErr,
						"account", e.Creds.AccountName,
					)
				}
			}
			return nil, fmt.Errorf("failed to create executor for account %s: %v", cred.AccountName, err)
		}
		logger.Debug("Created NATS executor",
			"account", cred.AccountName,
			"url", natsURL,
		)
		executors[cred.AccountName] = executor
	}

	n := &NATSServerTools{
		executors: executors,
	}

	// Initialize tool categories
	n.serverTools = NewServerTools(n)
	n.streamTools = NewStreamTools(n)

	logger.Info("Initialized NATS server tools",
		"num_accounts", len(executors),
	)

	return n, nil
}

// GetExecutor returns the executor for the specified account
func (n *NATSServerTools) GetExecutor(accountName string) (*common.NATSExecutor, error) {
	executor, ok := n.executors[accountName]
	if !ok {
		return nil, fmt.Errorf("no executor found for account %s", accountName)
	}
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
