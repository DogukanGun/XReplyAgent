package endpoints

import (
	"cg-mentions-bot/internal/types"
)

// GenerateEndpointTools generates MCP tools from all endpoint methods
func (ge *GoldrushEndpoints) GenerateEndpointTools() []types.ToolInfo {
	var tools []types.ToolInfo

	// Generate all endpoint tools
	activityTool, activityHandler := ge.GenerateActivityTool()
	tools = append(tools, types.ToolInfo{Tool: activityTool, Handler: activityHandler})

	bitcoinBalanceTool, bitcoinBalanceHandler := ge.GenerateBitcoinBalanceTool()
	tools = append(tools, types.ToolInfo{Tool: bitcoinBalanceTool, Handler: bitcoinBalanceHandler})

	erc20TransferTool, erc20TransferHandler := ge.GenerateErc20TokenTransferTool()
	tools = append(tools, types.ToolInfo{Tool: erc20TransferTool, Handler: erc20TransferHandler})

	gasPriceTool, gasPriceHandler := ge.GenerateGasPriceTool()
	tools = append(tools, types.ToolInfo{Tool: gasPriceTool, Handler: gasPriceHandler})

	historicalPriceTool, historicalPriceHandler := ge.GenerateHistoricalTokenPriceTool()
	tools = append(tools, types.ToolInfo{Tool: historicalPriceTool, Handler: historicalPriceHandler})

	multichainBalancesTool, multichainBalancesHandler := ge.GenerateMultichainBalancesTool()
	tools = append(tools, types.ToolInfo{Tool: multichainBalancesTool, Handler: multichainBalancesHandler})

	multichainTransactionsTool, multichainTransactionsHandler := ge.GenerateMultichainTransactionsTool()
	tools = append(tools, types.ToolInfo{Tool: multichainTransactionsTool, Handler: multichainTransactionsHandler})

	nativeTokenBalanceTool, nativeTokenBalanceHandler := ge.GenerateNativeTokenBalanceTool()
	tools = append(tools, types.ToolInfo{Tool: nativeTokenBalanceTool, Handler: nativeTokenBalanceHandler})

	nftsForAddressTool, nftsForAddressHandler := ge.GenerateNftsForAddressTool()
	tools = append(tools, types.ToolInfo{Tool: nftsForAddressTool, Handler: nftsForAddressHandler})

	tokenBalanceForAddressTool, tokenBalanceForAddressHandler := ge.GenerateTokenBalanceForAddressTool()
	tools = append(tools, types.ToolInfo{Tool: tokenBalanceForAddressTool, Handler: tokenBalanceForAddressHandler})

	transactionTool, transactionHandler := ge.GenerateTransactionTool()
	tools = append(tools, types.ToolInfo{Tool: transactionTool, Handler: transactionHandler})

	return tools
}
