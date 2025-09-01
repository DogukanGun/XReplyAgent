package bnb

import (
	"context"
	"fmt"
	"log"
)

func main() {
	agent := ProxyHandler()

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
}
