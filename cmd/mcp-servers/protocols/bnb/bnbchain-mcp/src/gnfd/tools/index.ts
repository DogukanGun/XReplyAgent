import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"

import { registerAccountTools } from "./account.ts"
import { registerPaymentTools } from "./payment.ts"
import { registerStorageTools } from "./storage.ts"

export * from "./common.ts"

/**
 * Register all Greenfield-related tools
 */
export function registerGnfdTools(server: McpServer) {
  registerAccountTools(server)
  registerStorageTools(server)
  registerPaymentTools(server)
}
