import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"

import { registerEVM } from "../evm"
import { registerGnfd } from "../gnfd"
import Logger from "../utils/logger.js"

// Create and start the MCP server
export const startServer = () => {
  try {
    // Create a new MCP server instance
    const server = new McpServer({
      name: "BNBChain MCP Server",
      version: "1.0.0"
    })

    // Register all resources, tools, and prompts
    registerEVM(server)
    registerGnfd(server)
    return server
  } catch (error) {
    Logger.error("Failed to initialize server:", error)
    process.exit(1)
  }
}
