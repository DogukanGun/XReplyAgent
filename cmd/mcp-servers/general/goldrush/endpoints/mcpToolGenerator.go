package endpoints

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

// ToolInfo contains a tool and its handler function
type ToolInfo struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// GenerateEndpointTools generates MCP tools from all endpoint methods
func (ge *GoldrushEndpoints) GenerateEndpointTools() []ToolInfo {
	var tools []ToolInfo

	// Generate all endpoint tools
	activityTool, activityHandler := ge.GenerateActivityTool()
	tools = append(tools, ToolInfo{Tool: activityTool, Handler: activityHandler})

	bitcoinBalanceTool, bitcoinBalanceHandler := ge.GenerateBitcoinBalanceTool()
	tools = append(tools, ToolInfo{Tool: bitcoinBalanceTool, Handler: bitcoinBalanceHandler})

	erc20TransferTool, erc20TransferHandler := ge.GenerateErc20TokenTransferTool()
	tools = append(tools, ToolInfo{Tool: erc20TransferTool, Handler: erc20TransferHandler})

	gasPriceTool, gasPriceHandler := ge.GenerateGasPriceTool()
	tools = append(tools, ToolInfo{Tool: gasPriceTool, Handler: gasPriceHandler})

	historicalPriceTool, historicalPriceHandler := ge.GenerateHistoricalTokenPriceTool()
	tools = append(tools, ToolInfo{Tool: historicalPriceTool, Handler: historicalPriceHandler})

	multichainBalancesTool, multichainBalancesHandler := ge.GenerateMultichainBalancesTool()
	tools = append(tools, ToolInfo{Tool: multichainBalancesTool, Handler: multichainBalancesHandler})

	multichainTransactionsTool, multichainTransactionsHandler := ge.GenerateMultichainTransactionsTool()
	tools = append(tools, ToolInfo{Tool: multichainTransactionsTool, Handler: multichainTransactionsHandler})

	nativeTokenBalanceTool, nativeTokenBalanceHandler := ge.GenerateNativeTokenBalanceTool()
	tools = append(tools, ToolInfo{Tool: nativeTokenBalanceTool, Handler: nativeTokenBalanceHandler})

	nftsForAddressTool, nftsForAddressHandler := ge.GenerateNftsForAddressTool()
	tools = append(tools, ToolInfo{Tool: nftsForAddressTool, Handler: nftsForAddressHandler})

	tokenBalanceForAddressTool, tokenBalanceForAddressHandler := ge.GenerateTokenBalanceForAddressTool()
	tools = append(tools, ToolInfo{Tool: tokenBalanceForAddressTool, Handler: tokenBalanceForAddressHandler})

	transactionTool, transactionHandler := ge.GenerateTransactionTool()
	tools = append(tools, ToolInfo{Tool: transactionTool, Handler: transactionHandler})

	return tools
}
