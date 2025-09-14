import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"
import { z } from "zod"
import { withTwitterAuth } from "../middleware/twitter.js"
import { mcpToolRes } from "../utils/helper.js"

// Import all same-chain swap services
// Aptos services
import { kanaswap as aptosRecipientSwap } from "../services/same-chain/aptos/recipientSwap.js"
import { setIntegratorFee as aptosSetIntegratorFee } from "../services/same-chain/aptos/setIntegratorFee.js"
import { kanaswap as aptosSwap } from "../services/same-chain/aptos/swap.js"
import { kanaswap as aptosSwapUsingCore } from "../services/same-chain/aptos/swapUsingCore.js"
import { kanaswap as aptosSwapWithFee } from "../services/same-chain/aptos/swapWithFee.js"
import { kanaswap as aptosSwapWithFeeCoinIn } from "../services/same-chain/aptos/swapWith_fee_coin_in.js"

// EVM services
import { kanaswap as evmRecipientSwap } from "../services/same-chain/evm/recipientSwap.js"
import { kanaswap as evmSwap } from "../services/same-chain/evm/swap.js"

// Solana services
import { kanaswap as solanaRecipientSwap } from "../services/same-chain/solana/recipientSwap.js"
import { kanaswap as solanaSwap } from "../services/same-chain/solana/swap.js"

// Common parameter schemas
const baseSwapParams = {
  twitter_id: z.string().describe("Twitter ID of the user"),
  amountIn: z.string().describe("Amount to swap (without decimals)")
}

const recipientSwapParams = {
  ...baseSwapParams,
  recipient: z.string().describe("Recipient wallet address")
}

const integratorFeeParams = {
  twitter_id: z.string().describe("Twitter ID of the user"),
  feeBps: z.number().describe("Fee in basis points (10 bps = 0.1%)")
}

export function registerSameChainSwapTools(server: McpServer) {
  // Aptos same-chain swaps
  server.tool(
    "aptos_swap",
    "Perform a same-chain token swap on Aptos using Kana Labs aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await aptosSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Aptos same-chain swap")
      }
    }, "aptos_swap")
  )

  server.tool(
    "aptos_recipient_swap",
    "Perform a same-chain token swap on Aptos with custom recipient using Kana Labs aggregator",
    recipientSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn, recipient }) => {
      try {
        const result = await aptosRecipientSwap({
          privateKey,
          address,
          amountIn,
          recipient
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Aptos recipient swap")
      }
    }, "aptos_recipient_swap")
  )

  server.tool(
    "aptos_swap_with_fee",
    "Perform a same-chain token swap on Aptos with integrator fee using Kana Labs aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await aptosSwapWithFee({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Aptos swap with fee")
      }
    }, "aptos_swap_with_fee")
  )

  server.tool(
    "aptos_swap_with_fee_coin_in",
    "Perform a same-chain token swap on Aptos with fee coin in using Kana Labs aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await aptosSwapWithFeeCoinIn({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Aptos swap with fee coin in")
      }
    }, "aptos_swap_with_fee_coin_in")
  )

  server.tool(
    "aptos_swap_using_core",
    "Get swap instruction for Aptos using Kana Labs core functions",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await aptosSwapUsingCore({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Aptos swap using core")
      }
    }, "aptos_swap_using_core")
  )

  server.tool(
    "aptos_set_integrator_fee",
    "Set integrator referral fee for Kana Labs router on Aptos",
    integratorFeeParams,
    withTwitterAuth(async ({ privateKey, address, feeBps }) => {
      try {
        const result = await aptosSetIntegratorFee({
          privateKey,
          address,
          feeBps
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "setting Aptos integrator fee")
      }
    }, "aptos_set_integrator_fee")
  )

  // EVM same-chain swaps
  server.tool(
    "evm_swap",
    "Perform a same-chain token swap on EVM (Polygon) using Kana Labs aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await evmSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing EVM same-chain swap")
      }
    }, "evm_swap")
  )

  server.tool(
    "evm_recipient_swap",
    "Perform a same-chain token swap on EVM (Polygon) with custom recipient using Kana Labs aggregator",
    recipientSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn, recipient }) => {
      try {
        const result = await evmRecipientSwap({
          privateKey,
          address,
          amountIn,
          recipient
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing EVM recipient swap")
      }
    }, "evm_recipient_swap")
  )

  // Solana same-chain swaps
  server.tool(
    "solana_swap",
    "Perform a same-chain token swap on Solana using Kana Labs aggregator",
    baseSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn }) => {
      try {
        const result = await solanaSwap({
          privateKey,
          address,
          amountIn
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Solana same-chain swap")
      }
    }, "solana_swap")
  )

  server.tool(
    "solana_recipient_swap",
    "Perform a same-chain token swap on Solana with custom recipient using Kana Labs aggregator",
    recipientSwapParams,
    withTwitterAuth(async ({ privateKey, address, amountIn, recipient }) => {
      try {
        const result = await solanaRecipientSwap({
          privateKey,
          address,
          amountIn,
          recipient
        })
        return mcpToolRes.success(result)
      } catch (error) {
        return mcpToolRes.error(error, "executing Solana recipient swap")
      }
    }, "solana_recipient_swap")
  )
} 