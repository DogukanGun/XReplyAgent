import { SwapAggregator, Environment, NetworkId } from "@kanalabs/aggregator"
import "dotenv/config"
import { clusterApiUrl, Connection, Keypair } from "@solana/web3.js"
import bs58 from "bs58"
import { ethers } from "ethers"

export const kanaswap = async (params: {
  privateKey: string
  address: string,
  amountIn: string, //This number is wo decimals
}) => {

  // Setup EVM Signer (Arbitrum)
  const evmprivateKey = params.privateKey
  const arbitrumRpc = process.env.ETH_ARBITRUM_RPC as string

  const evmProvider = new ethers.JsonRpcProvider(arbitrumRpc)
  const evmSigner = new ethers.Wallet(evmprivateKey, evmProvider)

  // Setup Solana Signer
  const solanaSigner = Keypair.fromSecretKey(
    bs58.decode(params.privateKey),
  )
  const solanaProvider = new Connection(
    clusterApiUrl("mainnet-beta"),
    "confirmed",
  )
  // Setup Kana swap aggregator
  const crossChainAggregator = new SwapAggregator(Environment.production, {
    providers: {
      //@ts-ignore
      arbitrum: evmProvider,
      solana: solanaProvider,
    },
    signers: {
      //@ts-ignore
      arbitrum: evmSigner,
      solana: solanaSigner,
    },
  })

  // Step 1: Get cross-chain quotes
  const crossChainQuotes = await crossChainAggregator.crossChainQuote({
    apiKey: process.env.KANA_API_KEY || "",
    sourceToken: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", //ETH
    targetToken: "So11111111111111111111111111111111111111112", //SOL
    sourceChain: NetworkId.Arbitrum,
    targetChain: NetworkId.solana,
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
    targetAddress: solanaSigner.publicKey.toBase58(),
  })
  console.log("Transfer executed successfully!")
  console.log("Transaction hash:", transfer)

  // Step 3: Execute claim (target chain transaction)
  const claim = await crossChainAggregator.executeClaim({
    apiKey: process.env.KANA_API_KEY || "",
    txHash: transfer.txHash,
    sourceProvider: evmProvider,
    targetProvider: solanaProvider,
    targetSigner: solanaSigner,
    quote: crossChainQuotes.data[0],
    sourceAddress: (await evmSigner.getAddress()) as string,
    targetAddress: solanaSigner.publicKey.toBase58(),
  })
  console.log("Tokens claimed successfully!")
  console.log("Transaction hash:", claim)
}
