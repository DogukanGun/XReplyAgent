import { Server } from "@modelcontextprotocol/sdk/server/index.js"
import { CallToolRequestSchema, ListToolsRequestSchema } from "@modelcontextprotocol/sdk/types.js"
import { runAptosAgent } from "../index.js"

// Create and start the MCP server
export const startServer = () => {
  try {
    // Create a new MCP server instance
    const server = new Server(
      { name: "Aptos MCP Server", version: "1.0.0" },
      { capabilities: { tools: {} } }
    )

    // Register Aptos tools
    registerAptosTools(server)
    
    return server
  } catch (error) {
    console.error("Failed to initialize server:", error)
    process.exit(1)
  }
}

// Register Aptos-specific tools
const registerAptosTools = (server: Server) => {
  // List available tools
  server.setRequestHandler(ListToolsRequestSchema, async () => {
    return {
      tools: [
        {
          name: "aptos_handle",
          description: "Handle Aptos on-chain operations. Provide message and twitter_id.",
          inputSchema: {
            type: "object",
            properties: {
              message: { type: "string", description: "Natural language instruction for Aptos agent" },
              twitter_id: { type: "string", description: "User's twitter_id for wallet lookup" },
            },
            required: ["message"],
          },
        },
      ],
    }
  })

  // Handle tool calls
  server.setRequestHandler(CallToolRequestSchema, async (request) => {
    if (request.params.name !== "aptos_handle") {
      return {
        isError: true,
        content: [{ type: "text", text: `Unknown tool: ${request.params.name}` }],
      }
    }
    
    const args = (request.params.arguments as any) || {}
    const message = String(args.message || "")
    const twitterId = String(args.twitter_id)
    
    if (!message) {
      return {
        isError: true,
        content: [{ type: "text", text: "'message' is required" }],
      }
    }
    
    try {
      const result = await runAptosAgent(message, twitterId)
      return {
        content: [{ type: "text", text: result.final }],
      }
    } catch (error: any) {
      return {
        isError: true,
        content: [{ type: "text", text: `Error: ${error?.message || String(error)}` }],
      }
    }
  })
}
