# mcp-nats

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for [NATS](https://nats.io/) messaging system integration

## Overview

This project provides a Model Context Protocol (MCP) server for NATS, enabling AI models and applications to interact with NATS messaging systems through a standardized interface. It exposes a comprehensive set of tools for interacting with NATS servers, making it ideal for AI-powered applications that need to work with messaging systems.

## What is MCP?

The Model Context Protocol (MCP) is an open protocol that standardizes how applications provide context to Large Language Models (LLMs). This server implements the MCP specification to provide NATS messaging capabilities to LLMs and AI applications, allowing them to:

- Interact with NATS messaging systems in a standardized way
- Safely inspect and monitor NATS servers and streams
- Perform read-only operations through a secure interface
- Integrate with other MCP-compatible clients and hosts

## Features
- Server Management (Read-only Operations)
  - List and inspect NATS servers
  - Server health monitoring and ping
  - Server information retrieval
- Stream Operations (Read-only Operations)
  - View and inspect NATS streams
  - Stream state and information queries
  - Message viewing and retrieval
  - Subject inspection
- Multi-Account Support
  - Handle multiple NATS accounts simultaneously
  - Secure credential management
- MCP Integration
  - Implements MCP server specification
  - Compatible with MCP clients like Claude Desktop
  - Standardized tool definitions for LLM interaction
  - Safe, read-only operations for AI interaction with NATS

## Requirements
- Go 1.24 or later
- NATS server (accessible via URL)
- NATS credentials for authentication
- MCP-compatible client (e.g., Claude Desktop, or other MCP clients)

## Installation

### Using Go
```sh
go install github.com/sinadarbouy/mcp-nats/cmd/mcp-nats@latest
```

### Building from Source
```sh
git clone https://github.com/sinadarbouy/mcp-nats.git
cd mcp-nats
go build -o mcp-nats ./cmd/mcp-nats
```

## Configuration

### Environment Variables
- `NATS_URL`: The URL of your NATS server (e.g., `nats://localhost:4222`)
- `NATS_<ACCOUNT>_CREDS`: Base64 encoded NATS credentials for each account
  - Example: `NATS_SYS_CREDS`, `NATS_A_CREDS`

### Command Line Flags
- `--transport`: Transport type (stdio or sse), default: stdio
- `--sse-address`: Address for SSE server to listen on, default: 0.0.0.0:8000
- `--log-level`: Log level (debug, info, warn, error), default: info
- `--json-logs`: Output logs in JSON format, default: false

### Example Usage
```sh
# Run with SSE transport and debug logging
./mcp-nats --transport sse --log-level debug

# Run with JSON logging
./mcp-nats --json-logs

# Run with custom SSE address
./mcp-nats --transport sse --sse-address localhost:9000
```

### Using VSCode with remote MCP server
Make sure your .vscode/settings.json includes:
```json
"mcp": {
  "servers": {
    "nats": {
      "type": "sse",
      "url": "http://localhost:8000/sse"
    }
  }
}
```
or 
cursor
```json
{
  "mcpServers": {
    "MCP_NATS_3": {
      "env": {
        "NATS_URL": "localhost:42222",
        "NATS_SYS_CREDS": "<base64 of SYS account creds>"
        "NATS_A_CREDS": "<base64 of SYS account creds>"
      },
      "url": "http://localhost:8000/sse"
    }
  }
}
```
If using the binary:
```json
{
  "mcpServers": {
    "nats": {
      "command": "mcp-nats",
      "args": [
        "--transport",
        "stdio"
      ],
      "env": {
        "NATS_URL": "localhost:42222",
        "NATS_SYS_CREDS": "<base64 of SYS account creds>"
        "NATS_A_CREDS": "<base64 of SYS account creds>"
      }
    }
  }
}
```

## Development

### Prerequisites
- Go 1.24+
- Docker (optional)
- NATS CLI
- Understanding of MCP specification

### Available Make Commands
```sh
make help      # Print help message
make build     # Build the binary
make run       # Run in stdio mode
make run-sse   # Run with SSE transport
make lint      # Run linters
```

## Resources
- [Model Context Protocol Documentation](https://modelcontextprotocol.io/introduction)
- [MCP Specification](https://modelcontextprotocol.io)
- [Example MCP Servers](https://modelcontextprotocol.io/example-servers)
