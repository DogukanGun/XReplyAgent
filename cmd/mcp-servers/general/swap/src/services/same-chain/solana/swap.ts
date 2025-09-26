import { SwapAggregator, Environment, NetworkId } from "@kanalabs/aggregator"
import "dotenv/config"
import bs58 from "bs58"
import { Connection, Keypair, clusterApiUrl } from "@solana/web3.js"

export const kanaswap = async (params: {
  privateKey: string
  address: string,
  amountIn: string, //This number is wo decimals
}) => {

  // Setup Signer
  const solanaSigner = Keypair.fromSecretKey(
    bs58.decode(params.privateKey),
  )
  const solanaProvider = new Connection(
    clusterApiUrl("mainnet-beta"),
    "confirmed",
  )

  // Setup Kana swap aggregator
  const swap = new SwapAggregator(Environment.production, {
    providers: {
      //@ts-ignore
      solana: solanaProvider,
    },
    signers: {
      //@ts-ignore
      solana: solanaSigner,
    },
  })

  // Step 1: Get quotes
  const quotes = await swap.swapQuotes({
    apiKey: process.env.KANA_API_KEY || "",
    inputToken: "So11111111111111111111111111111111111111112",
    outputToken: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
    amountIn: params.amountIn,
    slippage: 0.5,
    network: NetworkId.solana,
  })
  console.log("Quotes:", quotes)

  // Step 2: Execute swap with best quote
  const executeSwap = await swap.executeSwapInstruction({
    apiKey: process.env.KANA_API_KEY || "",
    quote: quotes.data[0], // Use first (best) quote
    address: solanaSigner.publicKey.toBase58(),
  })

  console.log("Transaction hash:", executeSwap)
}
