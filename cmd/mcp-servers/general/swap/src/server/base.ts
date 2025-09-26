import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"

import Logger from "../utils/logger"
import { registerSameChainSwapTools } from "../tools/same-chain-swap"
import { registerCrossChainSwapTools } from "../tools/cross-chain-swap"
import { registerRedeemTools } from "../tools/redeem"

// Set process timeout to handle long-running blockchain operations
process.env.UV_THREADPOOL_SIZE = "128" // Increase thread pool for concurrent operations

// Create and start the MCP server
export const startServer = () => {
  try {
    // Create a new MCP server instance
    const server = new McpServer({
      name: "BNBChain MCP Server",
      version: "1.0.0"
    })

    // Register all resources, tools, and prompts
    registerSameChainSwapTools(server)
    registerCrossChainSwapTools(server)
    registerRedeemTools(server)
    
    return server
  } catch (error) {
    Logger.error("Failed to initialize server:", error)
    process.exit(1)
  }
}
