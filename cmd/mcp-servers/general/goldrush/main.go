package main

import (
	"log"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new Goldrush endpoints instance
	goldrushEndpoints := endpoints.NewGoldrushEndpoints("")

	// Create the MCP server
	mcpServer := server.NewMCPServer(
		"GoldRush MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Generate all endpoint tools
	tools := goldrushEndpoints.GenerateEndpointTools()

	// Add all tools to the server
	for _, toolInfo := range tools {
		mcpServer.AddTool(toolInfo.Tool, toolInfo.Handler)
		log.Printf("Added tool: %s", toolInfo.Tool.Name)
	}

	log.Printf("GoldRush MCP Server started with %d tools", len(tools))

	// Start the stdio server
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Printf("Server error: %v\n", err)
	}
}
