package main

import (
	"cg-mentions-bot/cmd/mcp-servers/wallet/functions"
	"cg-mentions-bot/internal/utils/db"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	client, err := db.ConnectToDB(mongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	wf := &functions.WalletFunctions{
		MongoConnection: client,
		TwitterId:       "",
	}

	// Create MCP server
	walletMcpServer := server.NewMCPServer(
		"Wallet MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithLogging(),
	)

	// Register all tools
	for _, toolInfo := range wf.GenerateEndpointTools() {
		walletMcpServer.AddTool(toolInfo.Tool, toolInfo.Handler)
		log.Printf("Added tool: %s", toolInfo.Tool.Name)
	}

	// Run server
	port := "8085"
	if v, ok := os.LookupEnv("PORT"); ok {
		port = v
	}
	http := server.NewStreamableHTTPServer(walletMcpServer, server.WithEndpointPath("/mcp"), server.WithStateLess(true))
	if err := http.Start(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
