package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sinadarbouy/mcp-nats/internal/logger"
)

// ObjectTools represents all NATS object-related tools
type ObjectTools struct {
	nats *NATSServerTools
}

// NewObjectTools creates a new ObjectTools instance
func NewObjectTools(nats *NATSServerTools) *ObjectTools {
	return &ObjectTools{nats: nats}
}

// GetTools implements the ToolCategory interface
func (o *ObjectTools) GetTools() []Tool {
	return []Tool{
		{
			Tool: mcp.Tool{
				Name:        "object_add",
				Description: "Adds a new Object Store Bucket",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "A description for the bucket",
						},
						"ttl": map[string]interface{}{
							"type":        "string",
							"description": "How long to keep objects for",
						},
						"storage": map[string]interface{}{
							"type":        "string",
							"description": "Storage backend to use (file, memory)",
						},
						"replicas": map[string]interface{}{
							"type":        "integer",
							"description": "How many replicas of the data to store",
							"default":     1,
						},
						"max_bucket_size": map[string]interface{}{
							"type":        "string",
							"description": "Maximum size for the bucket",
						},
						"tags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Place the store on servers that has specific tags",
						},
						"cluster": map[string]interface{}{
							"type":        "string",
							"description": "Place the store on a specific cluster",
						},
						"metadata": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Adds metadata to the bucket",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: o.objectAddHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "object_put",
				Description: "Puts a file into the store",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"file": map[string]interface{}{
							"type":        "string",
							"description": "The file to put",
						},
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Override the name supplied to the object store",
						},
						"description": map[string]interface{}{
							"type":        "string",
							"description": "Sets an optional description for the object",
						},
						"header": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Adds headers to the object",
						},
						"progress": map[string]interface{}{
							"type":        "boolean",
							"description": "Enable/disable progress bars",
							"default":     true,
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Act without confirmation",
						},
						"data": map[string]interface{}{
							"type":        "string",
							"description": "Data to store, when empty reads from file",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: o.objectPutHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "object_get",
				Description: "Retrieves a file from the store",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"file": map[string]interface{}{
							"type":        "string",
							"description": "The file to retrieve",
						},
						"output": map[string]interface{}{
							"type":        "string",
							"description": "Override the output file name",
						},
						"progress": map[string]interface{}{
							"type":        "boolean",
							"description": "Enable/disable progress bars",
							"default":     true,
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Act without confirmation",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket", "file"},
				},
			},
			Handler: o.objectGetHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "object_del",
				Description: "Deletes a file or bucket from the store",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"file": map[string]interface{}{
							"type":        "string",
							"description": "The file to retrieve",
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Act without confirmation",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: o.objectDelHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "object_info",
				Description: "Get information about a bucket or object",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"file": map[string]interface{}{
							"type":        "string",
							"description": "The file to retrieve",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: o.objectInfoHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "object_ls",
				Description: "List buckets or contents of a specific bucket",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"names": map[string]interface{}{
							"type":        "boolean",
							"description": "When listing buckets, show just the bucket names",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name"},
				},
			},
			Handler: o.objectLsHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "object_seal",
				Description: "Seals a bucket preventing further updates",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
						"force": map[string]interface{}{
							"type":        "boolean",
							"description": "Force sealing without prompting",
						},
						"flags": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Optional flags to pass to the command",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: o.objectSealHandler(),
		},
		{
			Tool: mcp.Tool{
				Name:        "object_watch",
				Description: "Watch a bucket for changes",
				InputSchema: mcp.ToolInputSchema{
					Type: "object",
					Properties: map[string]interface{}{
						"account_name": map[string]interface{}{
							"type":        "string",
							"description": "The NATS account to use",
						},
						"bucket": map[string]interface{}{
							"type":        "string",
							"description": "The bucket to act on",
						},
					},
					Required: []string{"account_name", "bucket"},
				},
			},
			Handler: o.objectWatchHandler(),
		},
	}
}

// Helper function to get flags from arguments
func getObjectFlags(args map[string]interface{}) []string {
	if flags, ok := args["flags"].([]interface{}); ok {
		strFlags := make([]string, len(flags))
		for i, flag := range flags {
			if strFlag, ok := flag.(string); ok {
				strFlags[i] = strFlag
			}
		}
		return strFlags
	}
	return nil
}

func (o *ObjectTools) objectAddHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := o.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		args := []string{"object", "add", bucket}

		// Add optional parameters
		if description, ok := request.Params.Arguments["description"].(string); ok {
			args = append(args, fmt.Sprintf("--description=%s", description))
		}
		if ttl, ok := request.Params.Arguments["ttl"].(string); ok {
			args = append(args, fmt.Sprintf("--ttl=%s", ttl))
		}
		if storage, ok := request.Params.Arguments["storage"].(string); ok {
			args = append(args, fmt.Sprintf("--storage=%s", storage))
		}
		if replicas, ok := request.Params.Arguments["replicas"].(float64); ok {
			args = append(args, fmt.Sprintf("--replicas=%d", int(replicas)))
		}
		if maxBucketSize, ok := request.Params.Arguments["max_bucket_size"].(string); ok {
			args = append(args, fmt.Sprintf("--max-bucket-size=%s", maxBucketSize))
		}
		if tags, ok := request.Params.Arguments["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if strTag, ok := tag.(string); ok {
					args = append(args, fmt.Sprintf("--tags=%s", strTag))
				}
			}
		}
		if cluster, ok := request.Params.Arguments["cluster"].(string); ok {
			args = append(args, fmt.Sprintf("--cluster=%s", cluster))
		}
		if metadata, ok := request.Params.Arguments["metadata"].([]interface{}); ok {
			for _, meta := range metadata {
				if strMeta, ok := meta.(string); ok {
					args = append(args, fmt.Sprintf("--metadata=%s", strMeta))
				}
			}
		}

		// Add any additional flags
		if flags := getObjectFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		logger.Debug("Executing NATS object add command",
			"account", accountName,
			"bucket", bucket,
		)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (o *ObjectTools) objectPutHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := o.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		file, ok := request.Params.Arguments["file"].(string)
		if !ok {
			return nil, fmt.Errorf("missing file")
		}

		args := []string{"object", "put", bucket, file}

		// If data is provided, use it as stdin
		if data, ok := request.Params.Arguments["data"].(string); ok {
			executor.SetStdin(data)
		}

		// Add optional parameters
		if name, ok := request.Params.Arguments["name"].(string); ok {
			args = append(args, fmt.Sprintf("--name=%s", name))
		}
		if description, ok := request.Params.Arguments["description"].(string); ok {
			args = append(args, fmt.Sprintf("--description=%s", description))
		}
		if headers, ok := request.Params.Arguments["header"].([]interface{}); ok {
			for _, header := range headers {
				if strHeader, ok := header.(string); ok {
					args = append(args, "-H", strHeader)
				}
			}
		}
		if progress, ok := request.Params.Arguments["progress"].(bool); ok && !progress {
			args = append(args, "--no-progress")
		}
		if force, ok := request.Params.Arguments["force"].(bool); ok && force {
			args = append(args, "-f")
		}

		// Add any additional flags
		if flags := getObjectFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		logger.Debug("Executing NATS object put command",
			"account", accountName,
			"bucket", bucket,
			"file", file,
		)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (o *ObjectTools) objectGetHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := o.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		file, ok := request.Params.Arguments["file"].(string)
		if !ok {
			return nil, fmt.Errorf("missing file")
		}

		args := []string{"object", "get", bucket, file}

		// Add optional parameters
		if output, ok := request.Params.Arguments["output"].(string); ok {
			args = append(args, "-O", output)
		}
		if progress, ok := request.Params.Arguments["progress"].(bool); ok && !progress {
			args = append(args, "--no-progress")
		}
		if force, ok := request.Params.Arguments["force"].(bool); ok && force {
			args = append(args, "-f")
		}

		// Add any additional flags
		if flags := getObjectFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		logger.Debug("Executing NATS object get command",
			"account", accountName,
			"bucket", bucket,
			"file", file,
		)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (o *ObjectTools) objectDelHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := o.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		args := []string{"object", "del", bucket}

		// Add file if provided
		if file, ok := request.Params.Arguments["file"].(string); ok {
			args = append(args, file)
		}

		// Add force flag if true
		if force, ok := request.Params.Arguments["force"].(bool); ok && force {
			args = append(args, "-f")
		}

		// Add any additional flags
		if flags := getObjectFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		logger.Debug("Executing NATS object delete command",
			"account", accountName,
			"bucket", bucket,
		)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (o *ObjectTools) objectInfoHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := o.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"object", "info"}

		// Add bucket if provided
		var bucket string
		if b, ok := request.Params.Arguments["bucket"].(string); ok {
			bucket = b
			args = append(args, bucket)

			// Add file if provided and bucket is specified
			if file, ok := request.Params.Arguments["file"].(string); ok {
				args = append(args, file)
			}
		}

		// Add any additional flags
		if flags := getObjectFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		logger.Debug("Executing NATS object info command",
			"account", accountName,
			"bucket", bucket,
		)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (o *ObjectTools) objectLsHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := o.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		args := []string{"object", "ls"}

		// Add bucket if provided
		var bucket string
		if b, ok := request.Params.Arguments["bucket"].(string); ok {
			bucket = b
			args = append(args, bucket)
		}

		// Add names flag if true
		if names, ok := request.Params.Arguments["names"].(bool); ok && names {
			args = append(args, "-n")
		}

		// Add any additional flags
		if flags := getObjectFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		logger.Debug("Executing NATS object list command",
			"account", accountName,
			"bucket", bucket,
		)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (o *ObjectTools) objectSealHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := o.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		args := []string{"object", "seal", bucket}

		// Add force flag if true
		if force, ok := request.Params.Arguments["force"].(bool); ok && force {
			args = append(args, "-f")
		}

		// Add any additional flags
		if flags := getObjectFlags(request.Params.Arguments); flags != nil {
			args = append(args, flags...)
		}

		logger.Debug("Executing NATS object seal command",
			"account", accountName,
			"bucket", bucket,
		)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}

func (o *ObjectTools) objectWatchHandler() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		accountName, ok := request.Params.Arguments["account_name"].(string)
		if !ok {
			return nil, fmt.Errorf("missing account_name")
		}

		executor, err := o.nats.GetExecutor(ctx, accountName)
		if err != nil {
			return nil, err
		}

		bucket, ok := request.Params.Arguments["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("missing bucket")
		}

		args := []string{"object", "watch", bucket}

		logger.Debug("Executing NATS object watch command",
			"account", accountName,
			"bucket", bucket,
		)

		output, err := executor.ExecuteCommand(args...)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(output), nil
	}
}
