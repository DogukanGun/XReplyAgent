import { McpServer } from "@modelcontextprotocol/sdk/server/mcp";
import { z } from "zod";
import { withTwitterAuth } from "../../middleware/twitter";
import { mcpToolRes } from "../../utils/helper";
import { Environment, NetworkId, SwapAggregator } from "@kanalabs/aggregator";
import { ethers, utils } from "ethers";
import { getProvider } from "../../utils/chains";
import { Keypair } from "@solana/web3.js";
import { Ed25519Keypair } from '@mysten/sui.js/keypairs/ed25519';
import bs58 from "bs58";
import { Account, AccountAddress, Ed25519PrivateKey } from "@aptos-labs/ts-sdk";

export const registerSwapSameChainTool = (server: McpServer) => {
    server.tool(
        "swap_same_chain",
        "Swap tokens on the same chain",
        {
            network: z.nativeEnum(NetworkId).describe("The network to perform the swap on"),
            inputToken: z.string().describe("The input token address to swap from"),
            outputToken: z.string().describe("The output token address to swap to"),
            amountIn: z.number().describe("The amount of input tokens to swap"),
            slippage: z.number().optional().default(0.5).describe("Slippage tolerance percentage (default: 0.5)")
        },
        withTwitterAuth(async ({ sender_private_key, network, inputToken, outputToken, amountIn, slippage }) => {
            try {
                const provider = getProvider(network);
                let signer;
                const APIKEY = process.env.KANALABS_API_KEY as string;
                let address;
                if(network === NetworkId.solana) {
                    signer = Keypair.fromSecretKey(bs58.decode(sender_private_key));
                    address = signer.publicKey.toBase58();
                } else if(network === NetworkId.aptos) {
                    signer = Account.fromPrivateKey({
                        privateKey: new Ed25519PrivateKey(process.env.APTOS_PRIVATEKEY || ''),
                        address:  AccountAddress.from(process.env.APTOS_ADDRESS || ''),
                        legacy: true,
                    });
                    address = signer.publicKey;
                } else if(network === NetworkId.sui) {
                    signer = Ed25519Keypair.deriveKeypair(process.env.SUIMNEMONICS || '');
                    address = signer.getPublicKey()
                }  else {
                    signer = new ethers.Wallet(sender_private_key, provider);
                    address = signer.publicKey
                } 
                const swapAggregator = new SwapAggregator(
                    Environment.production,
                    {
                        
                    }
                );

                const quotes = await swapAggregator.swapQuotes({
                    apiKey: APIKEY,
                    inputToken: inputToken,
                    outputToken: outputToken,
                    amountIn: utils.parseUnits(amountIn.toString(), 18).toString(),
                    slippage: slippage,
                    network: network
                });

                if (!quotes.data || quotes.data.length === 0) {
                    return mcpToolRes.error("No quotes found for the swap", "swap_same_chain");
                }

                const optimalQuote = quotes.data[0];
                
                if (optimalQuote === undefined) {
                    return mcpToolRes.error("No quote found for the swap", "swap_same_chain");
                }
                const executeSwap = await swapAggregator.executeSwapInstruction({
                    apiKey: APIKEY,
                    quote: optimalQuote!,
                    address: address.toString(),
                });

                if (!executeSwap) {
                    return mcpToolRes.error("Swap execution failed", "swap_same_chain");
                }

                const data = {
                    result: executeSwap,
                    quote: optimalQuote,
                    address: address.toString()
                };

                return mcpToolRes.success(data);
            } catch (error) {
                return mcpToolRes.error(`Swap error: ${error instanceof Error ? error.message : 'Unknown error'}`, "swap_same_chain");
            }
        })
    );
}