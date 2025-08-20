import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"
import type { Address } from "viem"
import { z } from "zod"

import * as services from "../../services"
import { mcpToolRes } from "../../../utils/helper.ts"
import { defaultNetworkParam } from "../common/types.ts"

export function registerNftTools(server: McpServer) {
  // Get NFT (ERC721) information
  server.tool(
    "get_nft_info",
    "Get detailed information about a specific NFT (ERC721 token), including collection name, symbol, token URI, and current owner if available.",
    {
      tokenAddress: z
        .string()
        .describe(
          "The contract address of the NFT collection (e.g., '0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D' for Bored Ape Yacht Club)"
        ),
      tokenId: z
        .string()
        .describe("The ID of the specific NFT token to query (e.g., '1234')"),
      network: defaultNetworkParam
    },
    async ({ tokenAddress, tokenId, network }) => {
      try {
        const metadata = await services.getERC721TokenMetadata(
          tokenAddress as Address,
          BigInt(tokenId),
          network
        )

        return mcpToolRes.success(metadata)
      } catch (error) {
        return mcpToolRes.error(error, "fetching NFT metadata")
      }
    }
  )

  // Add tool for getting ERC1155 token URI
  server.tool(
    "get_erc1155_token_metadata",
    "Get the metadata for an ERC1155 token (multi-token standard used for both fungible and non-fungible tokens). The metadata typically points to JSON metadata about the token.",
    {
      tokenAddress: z
        .string()
        .describe(
          "The contract address of the ERC1155 token collection (e.g., '0x76BE3b62873462d2142405439777e971754E8E77')"
        ),
      tokenId: z
        .string()
        .describe(
          "The ID of the specific token to query metadata for (e.g., '1234')"
        ),
      network: defaultNetworkParam
    },
    async ({ tokenAddress, tokenId, network }) => {
      try {
        const metadata = await services.getERC1155TokenMetadata(
          tokenAddress as Address,
          BigInt(tokenId),
          network
        )

        return mcpToolRes.success(metadata)
      } catch (error) {
        return mcpToolRes.error(error, "fetching ERC1155 token URI")
      }
    }
  )

}
