import { SwapAggregator, Environment, NetworkId } from "@kanalabs/aggregator"
import "dotenv/config"
import { ethers } from "ethers"

export const kanaswap = async (params: {
  privateKey: string
  address: string,
  amountIn: string, //This number is wo decimals
}) => {

  // Setup EVM Signer (Arbitrum)
  const arbitrumprivateKey = params.privateKey
  const arbitrumRpc = process.env.ETH_ARBITRUM_RPC as string

  const arbitrumProvider = new ethers.JsonRpcProvider(arbitrumRpc)
  const arbitrumSigner = new ethers.Wallet(arbitrumprivateKey, arbitrumProvider)

  // Setup EVM Signer (Base)
  const baseprivateKey = params.privateKey
  const baseRpc = process.env.ETH_BASE_RPC as string

  const baseProvider = new ethers.JsonRpcProvider(baseRpc)
  const baseSigner = new ethers.Wallet(baseprivateKey, baseProvider)

  // Setup Kana swap aggregator
  const crossChainAggregator = new SwapAggregator(Environment.production, {
    providers: {
      //@ts-ignore
      arbitrum: arbitrumProvider,
      base: baseProvider,
    },
    signers: {
      //@ts-ignore
      arbitrum: arbitrumSigner,
      base: baseSigner,
    },
  })

  // Step 1: Get cross-chain quotes
  const crossChainQuotes = await crossChainAggregator.crossChainQuote({
    apiKey: process.env.KANA_API_KEY || "",
    sourceToken: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", //ETH
    targetToken: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", //ETH
    sourceChain: NetworkId.Arbitrum,
    targetChain: NetworkId.base,
    amountIn: params.amountIn, //0.0001 ETH
    sourceSlippage: 2, // 2% slippage
    targetSlippage: 2, // 2% slippage
  })
  console.log("Quotes:", crossChainQuotes)

  // Step 2: Execute transfer (source chain transaction)
  const transfer = await crossChainAggregator.executeTransfer({
    apiKey: process.env.KANA_API_KEY || "",
    sourceProvider: arbitrumProvider,
    sourceAddress: (await arbitrumSigner.getAddress()) as string,
    sourceSigner: arbitrumSigner,
    quote: crossChainQuotes.data[0],
    targetAddress: (await baseSigner.getAddress()) as string,
  })
  console.log("Transfer executed successfully!")
  console.log("Transaction hash:", transfer)

  // Step 3: Execute claim (target chain transaction)
  const claim = await crossChainAggregator.executeClaim({
    apiKey: process.env.KANA_API_KEY || "",
    txHash: transfer.txHash,
    sourceProvider: arbitrumProvider,
    targetProvider: baseProvider,
    targetSigner: baseSigner,
    quote: crossChainQuotes.data[0],
    sourceAddress: (await arbitrumSigner.getAddress()) as string,
    targetAddress: (await baseSigner.getAddress()) as string,
  })
  console.log("Tokens claimed successfully!")
  console.log("Transaction hash:", claim)
}
