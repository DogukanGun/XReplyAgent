import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"

import { registerWaletPrompts } from "./prompts.ts"
import { registerWalletTools } from "./tools.ts"

export function registerWallet(server: McpServer) {
  registerWalletTools(server)
  registerWaletPrompts(server)
}
