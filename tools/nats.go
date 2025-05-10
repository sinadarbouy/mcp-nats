package tools

import "github.com/sinadarbouy/mcp-nats/tools/common"

// NATSServerTools contains all NATS server-related tool definitions
type NATSServerTools struct {
	accExecutor *common.NATSExecutor
	sysExecutor *common.NATSExecutor
	serverTools *ServerTools
	streamTools *StreamTools
}

// NewNATSServerTools creates a new instance of NATSServerTools
func NewNATSServerTools(natsURL, AccNatsCredsPath, SysNatsCredsPath string) *NATSServerTools {
	accExecutor := &common.NATSExecutor{
		URL:       natsURL,
		CredsPath: AccNatsCredsPath,
	}
	sysExecutor := &common.NATSExecutor{
		URL:       natsURL,
		CredsPath: SysNatsCredsPath,
	}

	n := &NATSServerTools{
		accExecutor: accExecutor,
		sysExecutor: sysExecutor,
	}

	// Initialize tool categories
	n.serverTools = NewServerTools(n)
	n.streamTools = NewStreamTools(n)

	return n
}

// ServerTools returns the server tools category
func (n *NATSServerTools) ServerTools() ToolCategory {
	return n.serverTools
}

// StreamTools returns the stream tools category
func (n *NATSServerTools) StreamTools() ToolCategory {
	return n.streamTools
}

// GetSysExecutor returns the system executor
func (n *NATSServerTools) GetSysExecutor() *common.NATSExecutor {
	return n.sysExecutor
}

// GetAccExecutor returns the account executor
func (n *NATSServerTools) GetAccExecutor() *common.NATSExecutor {
	return n.accExecutor
}
