import { SwapAggregator, Environment, NetworkId, BridgeId } from "@kanalabs/aggregator"
import "dotenv/config"
import {
  Account,
  AccountAddress,
  Aptos,
  AptosConfig,
  Ed25519PrivateKey,
  Network,
} from "@aptos-labs/ts-sdk"
import { ethers } from "ethers"

export const kanaswap = async (params: {
  privateKey: string
  address: string,
  txHash: string,
  sourceChain: string,
  targetChain: string
}) => {

  // Setup Aptos Signer
  const aptosSigner = Account.fromPrivateKey({
    privateKey: new Ed25519PrivateKey(params.privateKey),
    address: AccountAddress.from(params.address),
    legacy: true,
  })

  // Setup EVM Signer (Polygon)
  const evmprivateKey = params.privateKey
  const polygonRpc = process.env.ETH_POLYGON_RPC as string

  const evmProvider = new ethers.JsonRpcProvider(polygonRpc)
  const evmSigner = new ethers.Wallet(evmprivateKey, evmProvider)

  // Setup Aptos provider
  const aptosConfig = new AptosConfig({ network: Network.MAINNET })
  const aptosProvider = new Aptos(aptosConfig)

  // Setup Kana swap aggregator
  const crossChainAggregator = new SwapAggregator(Environment.production, {
    providers: {
      //@ts-ignore
      aptos: aptosProvider,
      polygon: evmProvider,
    },
    signers: {
      //@ts-ignore
      aptos: aptosSigner,
      polygon: evmSigner,
    },
  })

  // Redeem cross-chain transaction
  const claim = await crossChainAggregator.redeem({
    apiKey: process.env.KANA_API_KEY || "",
    sourceChain: NetworkId.polygon,
    targetChain: NetworkId.aptos,
    sourceProvider: evmProvider,
    targetProvider: aptosProvider,
    SourceHash: params.txHash,
    targetAddress: aptosSigner.accountAddress.toString(),
    targetSigner: aptosSigner,
    BridgeId: BridgeId.cctp
  })
  console.log("Tokens redeemed successfully!")
  console.log("Transaction hash:", claim)
}
