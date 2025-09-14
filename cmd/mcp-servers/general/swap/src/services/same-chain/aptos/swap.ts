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

export const kanaswap = async (params: {
  privateKey: string
  address: string,
  amountIn: string, //This number is wo decimals
}) => {

  // Setup Signer
  const aptosSigner = Account.fromPrivateKey({
    privateKey: new Ed25519PrivateKey(params.privateKey),
    address: AccountAddress.from(params.address),
    legacy: true,
  })

  // Setup Aptos provider
  const aptosConfig = new AptosConfig({ network: Network.MAINNET })
  const aptosProvider = new Aptos(aptosConfig)

  // Setup Kana swap aggregator
  const swap = new SwapAggregator(Environment.production, {
    providers: {
      //@ts-ignore
      aptos: aptosProvider,
    },
    signers: {
      //@ts-ignore
      aptos: aptosSigner,
    },
  })

  // Step 1: Get quotes
  const quotes = await swap.swapQuotes({
    apiKey: process.env.KANA_API_KEY || "",
    inputToken: "0x1::aptos_coin::AptosCoin",
    outputToken:
      "0x6f986d146e4a90b828d8c12c14b6f4e003fdff11a8eecceceb63744363eaac01::mod_coin::MOD",
    amountIn: params.amountIn,
    slippage: 0.5,
    network: NetworkId.aptos,
  })
  console.log("Quotes:", quotes)

  // Step 2: Execute swap with best quote
  const executeSwap = await swap.executeSwapInstruction({
    apiKey: process.env.KANA_API_KEY || "",
    quote: quotes.data[0], // Use first (best) quote
    address: aptosSigner.accountAddress.toString(),
  })

  console.log("Transaction hash:", executeSwap)
}
