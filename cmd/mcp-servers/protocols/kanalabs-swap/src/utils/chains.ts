import { NetworkId } from "@kanalabs/aggregator";
import { ethers } from "ethers";

export const getRpcUrl = (networkId: NetworkId): string => {
    switch (networkId) {
        case NetworkId.polygon:
            return process.env.ETH_NODE_URI_POLYGON as string;
        case NetworkId.bsc:
            return process.env.ETH_NODE_URI_BSC as string;
        case NetworkId.ethereum:
            return process.env.ETH_NODE_URI_ETHEREUM as string;
        case NetworkId.Arbitrum:
            return process.env.ETH_NODE_URI_ARBITRUM as string;
        case NetworkId.Avalanche:
            return process.env.ETH_NODE_URI_AVALANCHE as string;
        case NetworkId.base:
            return process.env.ETH_NODE_URI_BASE as string;
        case NetworkId.solana:
            return process.env.SOLANA_NODE_URI as string;
        case NetworkId.aptos:
            return process.env.APTOS_NODE_URI as string;
        case NetworkId.sui:
            return process.env.SUI_NODE_URI as string;
        default:
            throw new Error(`Unsupported network: ${networkId}`);
    }
};

export const getProvider = (networkId: NetworkId) => {
    const rpcUrl = getRpcUrl(networkId);
    return ethers.getDefaultProvider(rpcUrl);
};