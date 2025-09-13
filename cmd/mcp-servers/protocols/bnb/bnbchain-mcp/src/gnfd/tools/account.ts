import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"
import type { Hex } from "viem"
import { z } from "zod"

import * as services from "../services"
import { getAddressFromPrivateKey } from "../services"
import { mcpToolRes } from "../../utils/helper.ts"
import { networkParam } from "./common.ts"
import { withTwitterAuth } from "../../middleware/twitter.ts"

export function registerAccountTools(server: McpServer) {
  // Get account balance
  server.tool(
    "gnfd_get_account_balance",
    "Get the balance for an account",
    {
      network: networkParam,
      address: z
        .string()
        .optional()
        .describe("The address of the account to get balance for"),
      twitter_id: z.string().describe("The Twitter id of the user")
    },
    withTwitterAuth(async ({ network, address, privateKey }) => {
      try {
        const balance = await services.getAccountBalance(network, {
          address: address || getAddressFromPrivateKey(privateKey as Hex)
        })
        return mcpToolRes.success(balance)
      } catch (error) {
        return mcpToolRes.error(error, "fetching account balance")
      }
    }, "gnfd_get_account_balance")
  )

  // Get all storage providers
  server.tool(
    "gnfd_get_all_sps",
    "Get a list of all storage providers in the Greenfield network",
    {
      network: networkParam
    },
    async ({ network }) => {
      try {
        const sps = await services.getAllSps(network)
        return mcpToolRes.success(sps)
      } catch (error) {
        return mcpToolRes.error(error, "fetching storage providers")
      }
    }
  )
}
