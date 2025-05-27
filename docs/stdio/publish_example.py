# publish_example.py
import asyncio
from mcp_client import StdioServerParameters, stdio_client, ClientSession

async def main():
    # Configure the server parameters
    server_params = StdioServerParameters(
        command="./mcp-nats",
        args=["--transport", "stdio"],
        env={
            "NATS_URL": "Your NATS server URL",  # Replace with your NATS server URL
            "NATS_A_CREDS": "Your base64 encoded credentials",  # Replace with your credentials
        }
    )
    
    # Create the connection via stdio transport
    async with stdio_client(server_params) as streams:
        # Create the client session with the streams
        async with ClientSession(*streams) as session:
            # Initialize the session
            await session.initialize()
            
            # List available tools to verify publish tool is available
            response = await session.list_tools()
            print("Available tools:", [tool.name for tool in response.tools])
            
            # Example 1: Simple message publish
            result = await session.call_tool("publish", {
                "account_name": "A",
                "subject": "test.message",
                "body": "Hello from test message!"
            })
            print("Publish result:", result.content)

if __name__ == "__main__":
    asyncio.run(main()) 
