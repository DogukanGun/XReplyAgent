import {
  Environment,
  getSameChainInstruction,
  getSameChainQuote,
  NetworkId,
} from "@kanalabs/aggregator"
import "dotenv/config"

export const kanaswap = async (params: {
  privateKey: string
  address: string,
  amountIn: string, //This number is wo decimals
}) => {

  // Step 1: Get quotes
  const quotes = await getSameChainQuote(
    {
      apiKey: process.env.KANA_API_KEY || "",
      inputToken: "0x1::aptos_coin::AptosCoin",
      outputToken:
        "0x6f986d146e4a90b828d8c12c14b6f4e003fdff11a8eecceceb63744363eaac01::mod_coin::MOD",
      amountIn: params.amountIn,
      slippage: 0.5,
      network: NetworkId.aptos,
    },
    Environment.production,
  )

  // Step 2: Get swap instruction
  const payload = await getSameChainInstruction(
    {
      apiKey: process.env.KANA_API_KEY || "",
      quote: quotes.data[0],
      address: params.address,
    },
    Environment.production,
  )

  console.log("swap Instruction", payload)
}
