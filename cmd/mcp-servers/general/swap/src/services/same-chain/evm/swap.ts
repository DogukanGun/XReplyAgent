import { SwapAggregator, Environment, NetworkId } from "@kanalabs/aggregator"
import "dotenv/config"
import { ethers, parseEther } from "ethers"

export const kanaswap = async (params: {
  privateKey: string
  address: string,
  amountIn: string, //This number is wo decimals
}) => {

  // Setup Signer
  const polygonRpc = process.env.ETH_POLYGON_RPC as string
  const polygonProvider = new ethers.JsonRpcProvider(polygonRpc)
  const polygonSigner = new ethers.Wallet(params.privateKey, polygonProvider)

  // Setup Kana swap aggregator
  const swap = new SwapAggregator(Environment.production, {
    providers: {
      //@ts-ignore
      polygon: polygonProvider,
    },
    signers: {
      //@ts-ignore
      polygon: polygonSigner,
    },
  })

  // Step 1: Get quotes
  const quotes = await swap.swapQuotes({
    apiKey: process.env.KANA_API_KEY || "",
    inputToken: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
    outputToken: "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
    amountIn: params.amountIn,
    slippage: 0.5,
    network: NetworkId.polygon,
  })
  console.log("Quotes:", quotes)

  // Step 2: Execute swap with best quote
  const executeSwap = await swap.executeSwapInstruction({
    apiKey: process.env.KANA_API_KEY || "",
    quote: quotes.data[0], // Use first (best) quote
    address: await polygonSigner.getAddress(),
  })

  console.log("Transaction hash:", executeSwap)
}
