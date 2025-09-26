import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"
import { z } from "zod"
import { withTwitterAuth } from "../middleware/twitter.js"
import { mcpToolRes } from "../utils/helper.js"

// Import all cross-chain swap services
import { kanaswap as aptosArbitrumSwap } from "../services/cross-chain/aptos-arbitrum-swap.js"
import { kanaswap as aptosSolanaSwap } from "../services/cross-chain/aptos-solana-swap.js"
import { kanaswap as arbitrumAptosSwap } from "../services/cross-chain/arbitrum-aptos-swap.js"
import { kanaswap as arbitrumBaseSwap } from "../services/cross-chain/arbitrum-base-swap.js"
import { kanaswap as arbitrumSolanaSwap } from "../services/cross-chain/arbitrum-solana-swap.js"
import { kanaswap as solanaAptosSwap } from "../services/cross-chain/solana-aptos-swap.js"
import { kanaswap as solanaArbitrumSwap } from "../services/cross-chain/solana-arbitrum-swap.js"

// Common parameter schemas
const baseSwapParams = {
  twitter_id: z.string().describe("Twitter ID of the user"),
  amountIn: z.string().describe("Amount to swap (without decimals)")
}

export function registerCrossChainSwapTools(server: McpServer) {
  // Aptos to Arbitrum swap
  server.tool(
    "aptos_arbitrum_swap",
    "Swap tokens from Aptos to Arbitrum using Kana Labs cross-chain aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await aptosArbitrumSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Aptos to Arbitrum swap")
      }
    }, "aptos_arbitrum_swap")
  )

  // Aptos to Solana swap
  server.tool(
    "aptos_solana_swap",
    "Swap tokens from Aptos to Solana using Kana Labs cross-chain aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await aptosSolanaSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Aptos to Solana swap")
      }
    }, "aptos_solana_swap")
  )

  // Arbitrum to Aptos swap
  server.tool(
    "arbitrum_aptos_swap",
    "Swap tokens from Arbitrum to Aptos using Kana Labs cross-chain aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await arbitrumAptosSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Arbitrum to Aptos swap")
      }
    }, "arbitrum_aptos_swap")
  )

  // Arbitrum to Base swap
  server.tool(
    "arbitrum_base_swap",
    "Swap tokens from Arbitrum to Base using Kana Labs cross-chain aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await arbitrumBaseSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Arbitrum to Base swap")
      }
    }, "arbitrum_base_swap")
  )

  // Arbitrum to Solana swap
  server.tool(
    "arbitrum_solana_swap",
    "Swap tokens from Arbitrum to Solana using Kana Labs cross-chain aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await arbitrumSolanaSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Arbitrum to Solana swap")
      }
    }, "arbitrum_solana_swap")
  )

  // Solana to Aptos swap
  server.tool(
    "solana_aptos_swap",
    "Swap tokens from Solana to Aptos using Kana Labs cross-chain aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await solanaAptosSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Solana to Aptos swap")
      }
    }, "solana_aptos_swap")
  )

  // Solana to Arbitrum swap
  server.tool(
    "solana_arbitrum_swap",
    "Swap tokens from Solana to Arbitrum using Kana Labs cross-chain aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await solanaArbitrumSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Solana to Arbitrum swap")
      }
    }, "solana_arbitrum_swap")
  )
} 