package functions

import (
	"cg-mentions-bot/internal/types"
)

// GenerateEndpointTools generates MCP tools from all WalletFunctions methods
func (wf *WalletFunctions) GenerateEndpointTools() []types.ToolInfo {
	var tools []types.ToolInfo

	// Create wallet
	createWalletTool, createWalletHandler := wf.GenerateCreateWalletTool()
	tools = append(tools, types.ToolInfo{Tool: createWalletTool, Handler: createWalletHandler})

	// Read wallet
	readWalletTool, readWalletHandler := wf.GenerateReadWalletTool()
	tools = append(tools, types.ToolInfo{Tool: readWalletTool, Handler: readWalletHandler})

	// Sign transaction
	signTxTool, signTxHandler := wf.GenerateSignTransactionTool()
	tools = append(tools, types.ToolInfo{Tool: signTxTool, Handler: signTxHandler})

	// Transfer asset
	transferAssetTool, transferAssetHandler := wf.GenerateTransferAssetTool()
	tools = append(tools, types.ToolInfo{Tool: transferAssetTool, Handler: transferAssetHandler})

	return tools
}
