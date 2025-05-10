package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// StreamTools represents all NATS stream-related tools
type StreamTools struct {
	nats *NATSServerTools
}

// NewStreamTools creates a new StreamTools instance
func NewStreamTools(nats *NATSServerTools) *StreamTools {
	return &StreamTools{
		nats: nats,
	}
}

// GetTools implements the ToolCategory interface
func (s *StreamTools) GetTools() []Tool {
	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "stream_info",
				Description: "Get information about a NATS stream",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"stream": map[string]interface{}{
							"type":        "string",
							"description": "The name of the stream to get information about",
						},
					},
					Required: []string{"account_name", "stream"},
				},
			},
			Handler: s.streamInfoHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "stream_list",
				Description: "List all known Streams",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: s.streamListHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "stream_report",
				Description: "Reports on Stream statistics",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: s.streamReportHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "stream_find",
				Description: "Finds streams matching certain criteria",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: s.streamFindHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "stream_state",
				Description: "Get the state of a NATS stream",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"stream": map[string]interface{}{
							"type":        "string",
							"description": "The name of the stream to get state for",
						},
					},
					Required: []string{"account_name", "stream"},
				},
			},
			Handler: s.streamStateHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "stream_subjects",
				Description: "Query subjects held in a stream",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"stream": map[string]interface{}{
							"type":        "string",
							"description": "Stream name",
						},
					},
					Required: []string{"account_name", "stream"},
				},
			},
			Handler: s.streamSubjectsHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "stream_view",
				Description: "View messages in a stream",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"stream": map[string]interface{}{
							"type":        "string",
							"description": "Stream name",
						},
						"size": map[string]interface{}{
							"type":        "integer",
							"description": "Page size",
						},
					},
					Required: []string{"account_name", "stream", "size"},
				},
			},
			Handler: s.streamViewHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "stream_get",
				Description: "Retrieves a specific message from a Stream",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"stream": map[string]interface{}{
							"type":        "string",
							"description": "Stream name",
						},
						"id": map[string]interface{}{
							"type":        "string",
							"description": "Message Sequence to retrieve",
						},
					},
					Required: []string{"account_name", "stream", "id"},
				},
			},
			Handler: s.streamGetHandler(),
		},
	}
}

func (s *StreamTools) streamInfoHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		stream, ok := request.Params.Arguments["stream"].(string)
		if !ok {
			return nil, fmt.Errorf("missing stream")
		}

		executor, err := s.nats.GetExecutor(accountName)
		if err != nil {
			return nil, err
		}

		output, err := executor.ExecuteCommand("stream", "info", stream)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (s *StreamTools) streamListHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := s.nats.GetExecutor(accountName)
		if err != nil {
			return nil, err
		}

		output, err := executor.ExecuteCommand("stream", "list")
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats stream report
func (s *StreamTools) streamReportHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := s.nats.GetExecutor(accountName)
		if err != nil {
			return nil, err
		}

		output, err := executor.ExecuteCommand("stream", "report")
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats stream find
func (s *StreamTools) streamFindHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := s.nats.GetExecutor(accountName)
		if err != nil {
			return nil, err
		}

		output, err := executor.ExecuteCommand("stream", "find")
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats stream state
// Args:
//
//	[<stream>]  Stream to retrieve state information for
func (s *StreamTools) streamStateHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		stream, ok := request.Params.Arguments["stream"].(string)
		if !ok {
			return nil, fmt.Errorf("missing stream")
		}

		executor, err := s.nats.GetExecutor(accountName)
		if err != nil {
			return nil, err
		}

		output, err := executor.ExecuteCommand("stream", "state", stream)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats stream subjects
// Args:
//
//	[<stream>]  Stream name
func (s *StreamTools) streamSubjectsHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		stream, ok := request.Params.Arguments["stream"].(string)
		if !ok {
			return nil, fmt.Errorf("missing stream")
		}

		executor, err := s.nats.GetExecutor(accountName)
		if err != nil {
			return nil, err
		}

		output, err := executor.ExecuteCommand("stream", "subjects", stream)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats stream view
// Args:
//
//		[<stream>]  Stream name
//	 [<size>]    Page size
func (s *StreamTools) streamViewHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		stream, ok := request.Params.Arguments["stream"].(string)
		if !ok {
			return nil, fmt.Errorf("missing stream")
		}

		size, ok := request.Params.Arguments["size"].(int)
		if !ok {
			return nil, fmt.Errorf("missing size")
		}

		executor, err := s.nats.GetExecutor(accountName)
		if err != nil {
			return nil, err
		}

		output, err := executor.ExecuteCommand("stream", "view", stream, strconv.Itoa(size))
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

// nats stream get
// Args:
//
//	[<stream>]  Stream name
//	[<id>]      Message Sequence to retrieve
func (s *StreamTools) streamGetHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		stream, ok := request.Params.Arguments["stream"].(string)
		if !ok {
			return nil, fmt.Errorf("missing stream")
		}

		id, ok := request.Params.Arguments["id"].(string)
		if !ok {
			return nil, fmt.Errorf("missing id")
		}

		executor, err := s.nats.GetExecutor(accountName)
		if err != nil {
			return nil, err
		}

		output, err := executor.ExecuteCommand("stream", "get", stream, id)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}
