package main

import (
	"log"
	"os"

	"cg-mentions-bot/cmd/mcp-servers/general/goldrush/endpoints"

	"github.com/mark3labs/mcp-go/server"
)

/*
This mcp server is created for GoldRush and requires these env variables
-	PORT // default 8083
-	GOLDRUSH_AUTH_TOKEN  // required
*/
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

	port := "8083"
	if v, ok := os.LookupEnv("PORT"); ok {
		port = v
	}
	http := server.NewStreamableHTTPServer(mcpServer, server.WithEndpointPath("/mcp"), server.WithStateLess(true))
	if err := http.Start(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
