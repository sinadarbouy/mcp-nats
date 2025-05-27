# Testing MCP Server using stdio Transport

This guide explains how to test the MCP server locally using stdio transport format.

## Prerequisites

- Go (for building the server)
- Python 3.x (for running test client)
- NATS credentials

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/sinadarbouy/mcp-nats.git
cd mcp-nats
```

2. Build the server:
```bash
go build -o mcp-nats ./cmd/mcp-nats
chmod +x mcp-nats
```

## Testing Methods

### Using Cursor IDE

If you're using Cursor IDE, you can configure the MCP server in your settings:

```json
{
  "mcpServers": {
    "MCP_NATS_STDIO": {
      "command": "./mcp-nats",
      "args": ["--transport", "stdio"],
      "env": {
        "NATS_URL": "Your NATS server URL",
        "NATS_A_CREDS": "<base64 of A account creds>"
      }
    }
  }
}
```

After configuration, you can interact with the server directly through Cursor using natural language commands. For example:

1. You can tell Cursor:
```
"publish dummy message to A account"
```

2. Cursor will automatically translate this to an MCP tool call:
```json
{
  "account_name": "A",
  "subject": "test.message",
  "body": "Hello from test message!"
}
```

3. The result will be displayed:
```
"Published 1 message(s) to test.message"
```

This natural language interface makes it easy to interact with the MCP server without needing to remember specific command formats or tool names.

### Using Python Test Client

If you're not using Cursor, you can test the server using the provided Python client:

1. Set up Python environment:
```bash
python -m venv venv
source venv/bin/activate
pip install --upgrade pip
pip install -r docs/stdio/requirements.txt
```

2. Run the example client (publish_example.py):
```python
from mcp_client import StdioServerParameters, stdio_client, ClientSession

async def main():
    server_params = StdioServerParameters(
        command="./mcp-nats",
        args=["--transport", "stdio"],
        env={
            "NATS_URL": "Your NATS server URL",
            "NATS_A_CREDS": "Your base64 encoded credentials",
        }
    )
    
    async with stdio_client(server_params) as streams:
        async with ClientSession(*streams) as session:
            await session.initialize()
            
            # List available tools
            response = await session.list_tools()
            print("Available tools:", [tool.name for tool in response.tools])
            
            # Example: Publish a message
            result = await session.call_tool("publish", {
                "account_name": "A",
                "subject": "test.message",
                "body": "Hello from test message!"
            })
            print("Publish result:", result.content)

if __name__ == "__main__":
    asyncio.run(main())
```

3. Execute the example:
```bash
python docs/stdio/publish_example.py
```

## Expected Output

When running the example client, you should see output similar to:

```
Available tools: ['account_backup', 'account_info', 'account_report_connections', ...]
Publish result: [TextContent(type='text', text='Published 1 message(s) to test.message')]
```

## Environment Variables

- `NATS_URL`: The URL of your NATS server
- `NATS_A_CREDS`: Base64 encoded credentials for the A account

## Notes

- Make sure to replace placeholder values (NATS_URL, NATS_A_CREDS) with your actual configuration
- The server must be built and executable before running any tests
- All credentials should be properly base64 encoded
