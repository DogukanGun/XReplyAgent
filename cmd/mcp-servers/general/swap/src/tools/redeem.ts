import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"
import { z } from "zod"
import { withTwitterAuth } from "../middleware/twitter.js"
import { mcpToolRes } from "../utils/helper.js"

// Import redeem service
import { kanaswap as redeem } from "../services/redeem/redeem"

// Parameter schema for redeem
const redeemParams = {
  twitter_id: z.string().describe("Twitter ID of the user"),
  txHash: z.string().describe("Transaction hash to redeem"),
  sourceChain: z.string().describe("Source blockchain network"),
  targetChain: z.string().describe("Target blockchain network")
}

export function registerRedeemTools(server: McpServer) {
  server.tool(
    "redeem_cross_chain_transaction",
    "Redeem a cross-chain transaction using Kana Labs aggregator",
    redeemParams,
    withTwitterAuth(async ({ privateKey, address, txHash, sourceChain, targetChain }) => {
      try {
        const result = await redeem({
          privateKey,
          address,
          txHash,
          sourceChain,
          targetChain
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "redeeming cross-chain transaction")
      }
    }, "redeem_cross_chain_transaction")
  )
} 