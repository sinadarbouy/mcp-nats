# mcp-nats

MCP server for NATS

## Overview

This project provides an MCP (Multi-Cloud Proxy) server for NATS, It exposes a set of tools for interacting with a NATS server, such as publishing and subscribing to messages, via a unified API.

## Features
- Publish messages to NATS subjects
- Subscribe to NATS subjects
- Request/reply pattern support
- Designed for extensibility and integration

## Usage

### 1. Build and Run

#### Using Go
```sh
go build -o mcp-nats ./cmd/mcp-nats
./mcp-nats
```

#### Using Docker
```sh
docker build -t mcp-nats .
docker run -p 8000:8000 -e NATS_URL=nats://localhost:4222 mcp-nats
```

### 2. Configuration

Set the following environment variables:
- `NATS_URL`: The URL of your NATS server (e.g., `nats://localhost:4222`)
- (Optional) `NATS_CREDS`: Path to NATS credentials file

### 3. Tools

The MCP NATS server exposes the following tools:
- `publish_message`: Publish a message to a subject
- `subscribe_subject`: Subscribe to a subject
- `request_reply`: Send a request and await a reply

## Development

- Written in Go
- Contributions welcome!

## License

Apache-2.0 

## Server Info

The MCP NATS server exposes the following HTTP endpoints:

- `GET /healthz` — Health check endpoint. Returns `200 OK` if the server is running.
- `POST /api/tools` — Main API endpoint for invoking NATS tools (e.g., publish, subscribe, request). (To be implemented)

The server listens on port `8000` by default.

## API Usage

Clients interact with the MCP NATS server via HTTP requests. Example usage for each tool will be documented as the implementation progresses.

### Example: Health Check

```
GET http://localhost:8000/healthz
Response: 200 OK
ok
```

### Example: Publish Message (Planned)

```
POST http://localhost:8000/api/tools
Content-Type: application/json

{
  "tool": "publish_message",
  "params": {
    "subject": "foo.bar",
    "message": "hello world"
  }
}
```

Response:
```
{
  "status": "ok"
}
``` 
