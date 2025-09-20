import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js"

import { registerBlocks } from "./modules/blocks"
import { registerContracts } from "./modules/contracts"
import { registerNetwork } from "./modules/network"
import { registerNFT } from "./modules/nft"
import { registerTokens } from "./modules/tokens"
import { registerTransactions } from "./modules/transactions"
import { registerWallet } from "./modules/wallet"

export function registerEVM(server: McpServer) {
  registerBlocks(server)
  registerContracts(server)
  registerNetwork(server)
  registerTokens(server)
  registerTransactions(server)
  registerWallet(server)
  registerNFT(server)
}
