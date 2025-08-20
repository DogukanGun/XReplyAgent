import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"

import { registerTransactionPrompts } from "./prompts.ts"
import { registerTransactionTools } from "./tools.ts"

export function registerTransactions(server: McpServer) {
  registerTransactionTools(server)
  registerTransactionPrompts(server)
}
