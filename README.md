# mcp-nats
[![Install MCP Server](https://cursor.com/deeplink/mcp-install-dark.svg)](https://cursor.com/en/install-mcp?name=NATS&config=eyJjb21tYW5kIjoiZG9ja2VyIHJ1biAtaSAtLXJtIC0taW5pdCAtZSBOQVRTX1VSTCAtZSBOQVRTX05PX0FVVEhFTlRJQ0FUSU9OIGNuYWRiL21jcC1uYXRzIC0tdHJhbnNwb3J0IHN0ZGlvIiwiZW52Ijp7Ik5BVFNfVVJMIjoibmF0czovL2xvY2FsaG9zdDo0MjIyIiwiTkFUU19OT19BVVRIRU5USUNBVElPTiI6InRydWUifX0%3D)
[<img alt="Install in VS Code" src="https://img.shields.io/badge/Install%20in%20VS%20Code-007ACC?style=for-the-badge&logo=visualstudiocode&logoColor=white&labelColor=007ACC&color=white&logoWidth=24&logoPadding=8">](https://insiders.vscode.dev/redirect?url=vscode%3Amcp%2Finstall%3F%7B%22name%22%3A%22NATS%22%2C%22command%22%3A%22docker%22%2C%22args%22%3A%5B%22run%22%2C%22-i%22%2C%22--rm%22%2C%22--init%22%2C%22-e%22%2C%22NATS_URL%22%2C%22-e%22%2C%22NATS_NO_AUTHENTICATION%22%2C%22cnadb%2Fmcp-nats%22%2C%22--transport%22%2C%22stdio%22%5D%2C%22env%22%3A%7B%22NATS_URL%22%3A%22nats%3A%2F%2Flocalhost%3A4222%22%2C%22NATS_NO_AUTHENTICATION%22%3A%22true%22%7D%7D)


A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for [NATS](https://nats.io/) messaging system integration

[![MCP Review Certified](https://img.shields.io/badge/MCP%20Review-Certified-brightgreen)](https://mcpreview.com/mcp-servers/sinadarbouy/mcp-nats)

**This MCP server is certified by [MCP Review](https://mcpreview.com/mcp-servers/sinadarbouy/mcp-nats).**

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
  - Round-trip time (RTT) measurement
- Stream Operations (Read-only Operations)
  - View and inspect NATS streams
  - Stream state and information queries
  - Message viewing and retrieval
  - Subject inspection
- Object Store Operations
  - Create and manage object store buckets
  - Put and get files from object stores
  - List buckets and their contents
  - Delete objects and buckets
  - Watch buckets for changes
  - Seal buckets to prevent updates
- Key-Value Operations
  - Create and manage KV buckets
  - Store and retrieve key-value pairs
  - Watch for KV updates
  - Delete keys and buckets
- Publish Operations
  - Publish messages to NATS subjects
  - Support for different message formats
  - Asynchronous message publishing
- Account Operations
  - View account information and metrics
  - Generate account reports (connections and statistics)
  - Create and restore account backups
  - Inspect TLS chain for connected servers
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
- `NATS_URL`: The URL of your NATS server (e.g., `localhost:4222`)
- `NATS_<ACCOUNT>_CREDS`: Base64 encoded NATS credentials for each account
  - Example: `NATS_SYS_CREDS`, `NATS_A_CREDS`
- `NATS_NO_AUTHENTICATION`: Set to "true" to enable anonymous connections (no credentials required)
- `NATS_USER`: Username or token for user/password authentication
- `NATS_PASSWORD`: Password for user/password authentication

### Command Line Flags
- `--transport`: Transport type (stdio or sse), default: stdio
- `--sse-address`: Address for SSE server to listen on, default: 0.0.0.0:8000
- `--log-level`: Log level (debug, info, warn, error), default: info
- `--json-logs`: Output logs in JSON format, default: false
- `--no-authentication`: Allow anonymous connections without credentials
- `--user`: NATS username or token (can also be set via NATS_USER env var)
- `--password`: NATS password (can also be set via NATS_PASSWORD env var)

### Authentication Methods

The MCP NATS server supports three authentication methods:

1. **Credentials-based Authentication** (default): Uses NATS credentials files
   - Set `NATS_<ACCOUNT>_CREDS` environment variables
   - Requires `account_name` parameter in all tools

2. **User/Password Authentication**: Uses username and password
   - Set `NATS_USER` and `NATS_PASSWORD` environment variables or use `--user` and `--password` flags

3. **Anonymous Authentication**: No authentication required
   - Set `NATS_NO_AUTHENTICATION=true` environment variable or use `--no-authentication` flag

### Example Usage
```sh
# Run with SSE transport and debug logging
./mcp-nats --transport sse --log-level debug

# Run with JSON logging
./mcp-nats --json-logs

# Run with custom SSE address
./mcp-nats --transport sse --sse-address localhost:9000

# Run with anonymous authentication
./mcp-nats --no-authentication

# Run with user/password authentication
./mcp-nats --user myuser --password mypass

# Run with environment variables for authentication
NATS_NO_AUTHENTICATION=true ./mcp-nats
NATS_USER=myuser NATS_PASSWORD=mypass ./mcp-nats
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
    "nats": {
      "env": {
        "NATS_URL": "nats://localhost:42222",
        "NATS_SYS_CREDS": "<base64 of SYS account creds>"
        "NATS_A_CREDS": "<base64 of A account creds>"
      },
      "url": "http://localhost:8000/sse"
    }
  }
}
```

**Anonymous Authentication:**
```json
{
  "mcpServers": {
    "nats": {
      "env": {
        "NATS_URL": "nats://localhost:42222",
        "NATS_NO_AUTHENTICATION": "true"
      },
      "url": "http://localhost:8000/sse"
    }
  }
}
```

**User/Password Authentication:**
```json
{
  "mcpServers": {
    "nats": {
      "env": {
        "NATS_URL": "nats://localhost:42222",
        "NATS_USER": "myuser",
        "NATS_PASSWORD": "mypass"
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
        "NATS_URL": "nats://localhost:42222",
        "NATS_SYS_CREDS": "<base64 of SYS account creds>",
        "NATS_A_CREDS": "<base64 of A account creds>"
      }
    }
  }
}
```

**Anonymous Authentication with Binary:**
```json
{
  "mcpServers": {
    "nats": {
      "command": "mcp-nats",
      "args": [
        "--transport",
        "stdio",
        "--no-authentication"
      ],
      "env": {
        "NATS_URL": "nats://localhost:4222"
      }
    }
  }
}
```

**User/Password Authentication with Binary:**
```json
{
  "mcpServers": {
    "nats": {
      "command": "mcp-nats",
      "args": [
        "--transport",
        "stdio",
        "--user",
        "myuser"
      ],
      "env": {
        "NATS_URL": "nats://localhost:4222",
        "NATS_PASSWORD": "mypass"
      }
    }
  }
}
```

**Docker Configuration:**
```json
{
  "mcpServers": {
    "nats": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "--init",
        "-e",
        "NATS_URL",
        "-e",
        "NATS_SYS_CREDS",
        "cnadb/mcp-nats",
        "--transport",
        "stdio"
      ],
      "env": {
        "NATS_SYS_CREDS": "<base64 of SYS account creds>",
        "NATS_URL": "<nats url>"
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

## Testing with stdio Transport

For detailed instructions on how to test the MCP server using stdio transport, please refer to our [Stdio Example Guide](docs/stdio/stdio_example.md).

## Resources
- [Model Context Protocol Documentation](https://modelcontextprotocol.io/introduction)
- [MCP Specification](https://modelcontextprotocol.io)
- [Example MCP Servers](https://modelcontextprotocol.io/example-servers)
