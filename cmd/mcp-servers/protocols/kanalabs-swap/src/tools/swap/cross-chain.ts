import { McpServer } from "@modelcontextprotocol/sdk/server/mcp";
import { object, z } from "zod";
import { withTwitterAuth } from "../../middleware/twitter";
import { mcpToolRes } from "../../utils/helper";
import { BridgeId, Environment, EvmSignerType, NetworkId, SwapAggregator } from "@kanalabs/aggregator";
import { ethers, providers, utils } from "ethers";
import { getProvider } from "../../utils/chains";

export const registerSwapCrossChainTool = (server: McpServer) => {
    server.tool(
        "supported_chains",
        "Get supported chains",
        {},
        async () => {
            return mcpToolRes.success(Object.values(NetworkId))
        }
    )
    server.tool(
        "swap_cross_chain_polygon_to_bsc",
        "Swap tokens between two chains. Supported chains: polygon, bsc, ethereum, arbitrum, avalanche, base, solana, aptos, sui",
        {
            sourceChain: z.nativeEnum(NetworkId).describe("The source chain to swap from"),
            destinationChain: z.nativeEnum(NetworkId).describe("The destination chain to swap to"),
            destinationAddress: z.string().describe("The destination address to swap to"),
            sourceToken: z.string().describe("The source token to swap from"),
            targetToken: z.string().describe("The target token to swap to"),
            amount: z.number().describe("The amount of tokens to swap")
        },
        withTwitterAuth(async ({ sender_private_key, receiver_private_key, destinationAddress, amount, sourceToken, targetToken, sourceChain, destinationChain }) => {
            const sourceProvider = getProvider(sourceChain);
            const destinationProvider = getProvider(destinationChain);
            
            const sourceSigner = new ethers.Wallet(sender_private_key, sourceProvider);
            const targetSigner = new ethers.Wallet(receiver_private_key, destinationProvider);
            const APIKEY = process.env.KANALABS_API_KEY as string;
            const crossChainAggregator = new SwapAggregator(Environment.production);

            const quotes = await crossChainAggregator.crossChainQuote({
                apiKey: APIKEY,
                sourceToken: sourceToken,
                targetToken: targetToken,
                sourceChain: sourceChain,
                targetChain: destinationChain,
                amountIn: utils.parseUnits(amount.toString(), 18).toString(),
                sourceSlippage: 0.1,
                targetSlippage: 0.1,
            });
            let shouldClaim = true;
            const optimalQuote = quotes.data[0];
            shouldClaim = !(optimalQuote?.bridge === BridgeId.cctp); // cctp is not supported by claim
            optimalQuote?.bridge
            if (!optimalQuote) {
                return mcpToolRes.error("No quote found", "swap_cross_chain_polygon_to_bsc");
            }
            const sourceAddress = await sourceSigner.getAddress();
            const transfer = await crossChainAggregator.executeTransfer({
                apiKey: APIKEY,
                quote: optimalQuote,
                sourceAddress: sourceAddress,
                targetAddress: destinationAddress,
                sourceProvider: sourceProvider as providers.BaseProvider,
                sourceSigner: sourceSigner as EvmSignerType,
                
            });
            if (!transfer) {
                return mcpToolRes.error("Transfer failed", "swap_cross_chain_polygon_to_bsc");
            }
            if (!transfer.success) {
                const redeem = await crossChainAggregator.redeem({
                    apiKey: APIKEY,
                    sourceChain: sourceChain,
                    targetChain: destinationChain,
                    sourceProvider: sourceProvider as providers.BaseProvider,
                    targetProvider: destinationProvider as providers.BaseProvider,
                    targetSigner: targetSigner as EvmSignerType,
                    SourceHash: transfer.txHash,
                    targetAddress: destinationAddress,
                    BridgeId: optimalQuote?.bridge,
                });
                if (!redeem) {
                    return mcpToolRes.error("Redeem failed", "swap_cross_chain_polygon_to_bsc");
                }
                if (!redeem.success) {
                    return mcpToolRes.error("Redeem failed", "swap_cross_chain_polygon_to_bsc");
                }
                return mcpToolRes.success("Transaction failed but redeem successful");
            }
            if (shouldClaim) {
                await crossChainAggregator.executeClaim({
                    apiKey: APIKEY,
                    txHash: transfer.txHash,
                    sourceProvider: sourceProvider as providers.BaseProvider,
                    targetProvider: destinationProvider as providers.BaseProvider,
                    targetSigner: targetSigner as EvmSignerType,
                    quote: optimalQuote,
                    sourceAddress: sourceAddress,
                    targetAddress: destinationAddress,
                });
            }
            const data = {
                txHash: transfer.txHash,
                success: transfer.success,
                quote: optimalQuote.toString(),
            }
            return mcpToolRes.success(data)
        })
    )

}