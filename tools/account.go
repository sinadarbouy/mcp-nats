package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// AccountTools represents all NATS account-related tools
type AccountTools struct {
	nats *NATSServerTools
}

// NewAccountTools creates a new AccountTools instance
func NewAccountTools(nats *NATSServerTools) *AccountTools {
	return &AccountTools{
		nats: nats,
	}
}

// GetTools implements the ToolCategory interface
func (a *AccountTools) GetTools() []Tool {
	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "account_info",
				Description: "Get information about a NATS account",
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
			Handler: a.accountInfoHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "account_report_connections",
				Description: "Report on connections",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"sort": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"subs", "in-bytes", "out-bytes", "in-msgs", "out-msgs", "uptime", "cid"},
							"description": "Sort by a specific property",
						},
						"top": map[string]interface{}{
							"type":        "integer",
							"description": "Limit results to the top results",
							"default":     1000,
						},
						"subject": map[string]interface{}{
							"type":        "string",
							"description": "Limits responses only to those connections with matching subscription interest",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: a.accountReportConnectionsHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "account_report_statistics",
				Description: "Report on server statistics",
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
			Handler: a.accountReportStatisticsHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "account_backup",
				Description: "Creates a backup of all JetStream Streams over the NATS network",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"target": map[string]interface{}{
							"type":        "string",
							"description": "Directory to create the backup in",
						},
						"check": map[string]interface{}{
							"type":        "boolean",
							"description": "Checks the Stream for health prior to backup",
							"default":     false,
						},
						"consumers": map[string]interface{}{
							"type":        "boolean",
							"description": "Enable or disable consumer backups",
							"default":     true,
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Perform backup without prompting",
							"default":     false,
						},
						"critical_warnings": map[string]interface{}{
							"type":        "boolean",
							"description": "Treat warnings as failures",
							"default":     false,
						},
					},
					Required: []string{"account_name", "target"},
				},
			},
			Handler: a.accountBackupHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "account_restore",
				Description: "Restore an account backup over the NATS network",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"directory": map[string]interface{}{
							"type":        "string",
							"description": "The directory holding the account backup to restore",
						},
						"cluster": map[string]interface{}{
							"type":        "string",
							"description": "Place the stream in a specific cluster",
						},
						"tags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Place the stream on servers that has specific tags (can specify multiple)",
						},
					},
					Required: []string{"account_name", "directory"},
				},
			},
			Handler: a.accountRestoreHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "account_tls",
				Description: "Report TLS chain for connected server",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"expire_warn": map[string]interface{}{
							"type":        "string",
							"description": "Warn about certs expiring this soon (e.g. '1w'; 0 to disable)",
						},
						"ocsp": map[string]interface{}{
							"type":        "boolean",
							"description": "Report OCSP information, if any",
							"default":     false,
						},
						"pem": map[string]interface{}{
							"type":        "boolean",
							"description": "Show PEM Certificate blocks",
							"default":     true,
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: a.accountTLSHandler(),
		},
	}
}

func (a *AccountTools) accountInfoHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := a.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"account", "info"}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (a *AccountTools) accountReportConnectionsHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := a.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"account", "report", "connections"}

		// Add sort flag if provided
		if sort, ok := request.Params.Arguments["sort"].(string); ok {
			args = append(args, fmt.Sprintf("--sort=%s", sort))
		}

		// Add top flag if provided
		if top, ok := request.Params.Arguments["top"].(float64); ok {
			args = append(args, fmt.Sprintf("--top=%d", int(top)))
		}

		// Add subject flag if provided
		if subject, ok := request.Params.Arguments["subject"].(string); ok {
			args = append(args, fmt.Sprintf("--subject=%s", subject))
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (a *AccountTools) accountReportStatisticsHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := a.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"account", "report", "statistics"}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (a *AccountTools) accountBackupHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		target, ok := request.Params.Arguments["target"].(string)
		if !ok {
			return nil, fmt.Errorf("missing target directory")
		}

		executor, err := a.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"account", "backup"}

		// Add check flag if true
		if check, ok := request.Params.Arguments["check"].(bool); ok && check {
			args = append(args, "--check")
		}

		// Add consumers flag (handle both true and false cases)
		if consumers, ok := request.Params.Arguments["consumers"].(bool); ok {
			if consumers {
				args = append(args, "--consumers")
			} else {
				args = append(args, "--no-consumers")
			}
		}

		// Add force flag if true
		if force, ok := request.Params.Arguments["force"].(bool); ok && force {
			args = append(args, "--force")
		}

		// Add critical-warnings flag if true
		if criticalWarnings, ok := request.Params.Arguments["critical_warnings"].(bool); ok && criticalWarnings {
			args = append(args, "--critical-warnings")
		}

		// Add target directory as the final argument
		args = append(args, target)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (a *AccountTools) accountRestoreHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		directory, ok := request.Params.Arguments["directory"].(string)
		if !ok {
			return nil, fmt.Errorf("missing directory")
		}

		executor, err := a.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"account", "restore"}

		// Add cluster flag if provided
		if cluster, ok := request.Params.Arguments["cluster"].(string); ok && cluster != "" {
			args = append(args, fmt.Sprintf("--cluster=%s", cluster))
		}

		// Add tags if provided
		if tags, ok := request.Params.Arguments["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if strTag, ok := tag.(string); ok && strTag != "" {
					args = append(args, fmt.Sprintf("--tag=%s", strTag))
				}
			}
		}

		// Add directory as the final argument
		args = append(args, directory)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (a *AccountTools) accountTLSHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := a.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"account", "tls"}

		// Add expire-warn flag if provided
		if expireWarn, ok := request.Params.Arguments["expire_warn"].(string); ok && expireWarn != "" {
			args = append(args, fmt.Sprintf("--expire-warn=%s", expireWarn))
		}

		// Add ocsp flag if true
		if ocsp, ok := request.Params.Arguments["ocsp"].(bool); ok && ocsp {
			args = append(args, "--ocsp")
		}

		// Add pem flag (handle both true and false cases)
		if pem, ok := request.Params.Arguments["pem"].(bool); ok {
			if pem {
				args = append(args, "--pem")
			} else {
				args = append(args, "--no-pem")
			}
		}

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}
