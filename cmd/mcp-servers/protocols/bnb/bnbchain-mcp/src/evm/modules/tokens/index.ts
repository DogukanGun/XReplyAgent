import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"

import { registerTokenPrompts } from "./prompts.ts"
import { registerTokenTools } from "./tools.ts"

export function registerTokens(server: McpServer) {
  registerTokenTools(server)
  registerTokenPrompts(server)
}
