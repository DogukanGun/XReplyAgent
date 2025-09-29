import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"

import { registerSwapSameChainTool } from "../tools/swap/same-chain"
import { registerSwapCrossChainTool } from "../tools/swap/cross-chain"
import Logger from "../utils/logger.js"

// Create and start the MCP server
export const startServer = () => {
  try {
    // Create a new MCP server instance
    const server = new McpServer({
      name: "Kanalabs Swap MCP Server",
      version: "1.0.0"
    })

    // Register all resources, tools, and prompts
    registerSwapSameChainTool(server)
    registerSwapCrossChainTool(server)
    return server
  } catch (error) {
    Logger.error("Failed to initialize server:", error)
    process.exit(1)
  }
}
