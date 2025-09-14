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
import { clusterApiUrl, Connection, Keypair } from "@solana/web3.js"
import bs58 from "bs58"

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

  // Setup Solana Signer
  const solanaSigner = Keypair.fromSecretKey(
    bs58.decode(params.privateKey),
  )
  const solanaProvider = new Connection(
    clusterApiUrl("mainnet-beta"),
    "confirmed",
  )
  // Setup Aptos provider
  const aptosConfig = new AptosConfig({ network: Network.MAINNET })
  const aptosProvider = new Aptos(aptosConfig)

  // Setup Kana swap aggregator
  const crossChainAggregator = new SwapAggregator(Environment.production, {
    providers: {
      //@ts-ignore
      aptos: aptosProvider,
      solana: solanaProvider,
    },
    signers: {
      //@ts-ignore
      aptos: aptosSigner,
      solana: solanaSigner,
    },
  })

  // Step 1: Get cross-chain quotes
  const crossChainQuotes = await crossChainAggregator.crossChainQuote({
    apiKey: process.env.KANA_API_KEY || "",
    sourceToken: "So11111111111111111111111111111111111111112", //SOL
    targetToken: "0x1::aptos_coin::AptosCoin", //APT
    sourceChain: NetworkId.solana,
    targetChain: NetworkId.aptos,
    amountIn: params.amountIn, //0.01 SOL
    sourceSlippage: 2, // 2% slippage
    targetSlippage: 2, // 2% slippage
  })
  console.log("Quotes:", crossChainQuotes)

  // Step 2: Execute transfer (source chain transaction)
  const transfer = await crossChainAggregator.executeTransfer({
    apiKey: process.env.KANA_API_KEY || "",
    sourceProvider: solanaProvider,
    sourceAddress: solanaSigner.publicKey.toBase58(),
    sourceSigner: solanaSigner,
    quote: crossChainQuotes.data[0],
    targetAddress: aptosSigner.accountAddress.toString(),
  })
  console.log("Transfer executed successfully!")
  console.log("Transaction hash:", transfer)

  // Step 3: Execute claim (target chain transaction)
  const claim = await crossChainAggregator.executeClaim({
    apiKey: process.env.KANA_API_KEY || "",
    txHash: transfer.txHash,
    sourceProvider: solanaProvider,
    targetProvider: aptosProvider,
    targetSigner: aptosSigner,
    quote: crossChainQuotes.data[0],
    sourceAddress: solanaSigner.publicKey.toBase58(),
    targetAddress: aptosSigner.accountAddress.toString(),
  })
  console.log("Tokens claimed successfully!")
  console.log("Transaction hash:", claim)
}
