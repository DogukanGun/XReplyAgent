import { SwapAggregator, Environment, NetworkId } from "@kanalabs/aggregator"
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
  amountIn: string, //This number is wo decimals
}) => {

  // Setup Aptos Signer
  const aptosSigner = Account.fromPrivateKey({
    privateKey: new Ed25519PrivateKey(params.privateKey),
    address: AccountAddress.from(params.address),
    legacy: true,
  })

  // Setup EVM Signer (Arbitrum)
  const evmprivateKey = params.privateKey
  const arbitrumRpc = process.env.ETH_ARBITRUM_RPC as string

  const evmProvider = new ethers.JsonRpcProvider(arbitrumRpc)
  const evmSigner = new ethers.Wallet(evmprivateKey, evmProvider)

  // Setup Aptos provider
  const aptosConfig = new AptosConfig({ network: Network.MAINNET })
  const aptosProvider = new Aptos(aptosConfig)

  // Setup Kana swap aggregator
  const crossChainAggregator = new SwapAggregator(Environment.production, {
    providers: {
      //@ts-ignore
      aptos: aptosProvider,
      arbitrum: evmProvider,
    },
    signers: {
      //@ts-ignore
      aptos: aptosSigner,
      arbitrum: evmSigner,
    },
  })
  // Step 1: Get cross-chain quotes
  const crossChainQuotes = await crossChainAggregator.crossChainQuote({
    apiKey: process.env.KANA_API_KEY || "",
    sourceToken: "0x1::aptos_coin::AptosCoin", //APT
    targetToken: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", //ETH
    sourceChain: NetworkId.aptos,
    targetChain: NetworkId.Arbitrum,
    amountIn: params.amountIn, //0.1 APT
    sourceSlippage: 2, // 2% slippage
    targetSlippage: 2, // 2% slippage
  })
  console.log("Quotes:", crossChainQuotes)

  // Step 2: Execute transfer (source chain transaction)
  const transfer = await crossChainAggregator.executeTransfer({
    apiKey: process.env.KANA_API_KEY || "",
    sourceProvider: aptosProvider,
    sourceAddress: aptosSigner.accountAddress.toString(),
    sourceSigner: aptosSigner,
    quote: crossChainQuotes.data[0],
    targetAddress: (await evmSigner.getAddress()) as string,
  })
  console.log("Transfer executed successfully!")
  console.log("Transaction hash:", transfer)

  // Step 3: Execute claim (target chain transaction)
  const claim = await crossChainAggregator.executeClaim({
    apiKey: process.env.KANA_API_KEY || "",
    txHash: transfer.txHash,
    sourceProvider: aptosProvider,
    targetProvider: evmProvider,
    targetSigner: evmSigner,
    quote: crossChainQuotes.data[0],
    sourceAddress: aptosSigner.accountAddress.toString(),
    targetAddress: (await evmSigner.getAddress()) as string,
  })
  console.log("Tokens claimed successfully!")
  console.log("Transaction hash:", claim)
}