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
    sourceToken: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", //ETH
    targetToken: "0x1::aptos_coin::AptosCoin", //APT
    sourceChain: NetworkId.Arbitrum,
    targetChain: NetworkId.aptos,
    amountIn: params.amountIn, //0.0001 ETH
    sourceSlippage: 2, // 2% slippage
    targetSlippage: 2, // 2% slippage
  })
  console.log("Quotes:", crossChainQuotes)

  // Step 2: Execute transfer (source chain transaction)
  const transfer = await crossChainAggregator.executeTransfer({
    apiKey: process.env.KANA_API_KEY || "",
    sourceProvider: evmProvider,
    sourceAddress: (await evmSigner.getAddress()) as string,
    sourceSigner: evmSigner,
    quote: crossChainQuotes.data[0],
    targetAddress: aptosSigner.accountAddress.toString(),
  })
  console.log("Transfer executed successfully!")
  console.log("Transaction hash:", transfer)

  // Step 3: Execute claim (target chain transaction)
  const claim = await crossChainAggregator.executeClaim({
    apiKey: process.env.KANA_API_KEY || "",
    txHash: transfer.txHash,
    sourceProvider: evmProvider,
    targetProvider: aptosProvider,
    targetSigner: aptosSigner,
    quote: crossChainQuotes.data[0],
    sourceAddress: (await evmSigner.getAddress()) as string,
    targetAddress: aptosSigner.accountAddress.toString(),
  })
  console.log("Tokens claimed successfully!")
  console.log("Transaction hash:", claim)
}
