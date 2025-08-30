package main

import (
	"cg-mentions-bot/cmd/mcp-servers/protocols/bnb"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"log"
)

func main() {
	agent := bnb.BnbProxy()

	ctx := context.Background()

	toolsInfo, err := agent.GetToolsInfo(ctx)
	if err != nil {
		log.Printf("Error getting tools info: %v", err)
		return
	}

	fmt.Println("Available tools:")
	fmt.Println(toolsInfo)

	langchainTool := agent.AsLangChainTool()
	fmt.Printf("\nLangChain Tool Name: %s\n", langchainTool.Name())
	fmt.Printf("LangChain Tool Description:\n%s\n", langchainTool.Description())

	toolRequest := bnb.CallRequest{
		ToolName: "get_latest_block",
		Parameters: map[string]interface{}{
			"network": "bsc",
		},
	}
	bytes, err := json.Marshal(toolRequest)
	if err != nil {
		return
	}
	res, err := langchainTool.Call(ctx, string(bytes))
	if err != nil {
		log.Printf("Error calling What is bnb doing: %v", err)
	}
	fmt.Printf("Res: %s", res)
	toolRequest = bnb.CallRequest{
		ToolName: "get_chain_info",
		Parameters: map[string]interface{}{
			"network": "bsc",
		},
	}
	bytes, err = json.Marshal(toolRequest)
	if err != nil {
		log.Printf("Error calling What is bnb doing: %v", err)
	}
	res, err = langchainTool.Call(ctx, string(bytes))
	fmt.Printf("\nRes2: %s", res)
}
