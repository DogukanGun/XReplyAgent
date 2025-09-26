package main

import (
	"context"
	"log"
	"os"
	"time"

	mcpclient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Upstream MOVE MCP over SSE
	upstream := os.Getenv("MOVE_MCP_SSE")
	if upstream == "" {
		log.Fatal("MOVE_MCP_SSE is required (e.g., http://localhost:3003/sse)")
	}

	ctx := context.Background()
	mcpCl, err := mcpclient.NewSSEMCPClient(upstream)
	if err != nil {
		log.Fatalf("failed to create SSE client: %v", err)
	}
	if err := mcpCl.Start(ctx); err != nil {
		log.Fatalf("failed to start SSE client: %v", err)
	}
	if _, err := mcpCl.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo:      mcp.Implementation{Name: "move-proxy", Version: "0.1.0"},
		},
	}); err != nil {
		log.Fatalf("failed to initialize upstream: %v", err)
	}

	// List upstream tools once at startup
	lt, err := mcpCl.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatalf("failed to list tools from upstream: %v", err)
	}

	// Local HTTP MCP server exposing the same tools
	s := server.NewMCPServer("move-proxy", "0.1.0", server.WithToolCapabilities(true), server.WithLogging())
	for _, tool := range lt.Tools {
		t := tool // capture
		s.AddTool(t, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Forward to upstream via the persistent SSE client
			log.Printf("Forwarding tool call: %s with args: %v", req.Params.Name, req.GetArguments())
			ctx2, cancel := context.WithTimeout(ctx, 120*time.Second)
			defer cancel()
			res, err := mcpCl.CallTool(ctx2, mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      req.Params.Name,
					Arguments: req.GetArguments(),
				},
			})
			if err != nil {
				log.Printf("Error calling upstream tool: %v", err)
				return mcp.NewToolResultError(err.Error()), nil
			}
			log.Printf("Received response from upstream: %+v", res)
			return res, nil
		})
	}

	port := getEnv("PORT", "8086")
	http := server.NewStreamableHTTPServer(s, server.WithEndpointPath("/mcp"), server.WithStateLess(true))
	log.Printf("move-proxy MCP server listening on :%s/mcp (forwarding to SSE %s)", port, upstream)
	if err := http.Start(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}
