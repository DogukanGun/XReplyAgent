package cg

import (
	"context"
	"fmt"

	mcpclient "cg-mentions-bot/internal/mcp"
)

// NewAsker returns a function that builds a concise prompt and asks the MCP tool.
func NewAsker(mcpCmd string, mcpTool string) func(ctx context.Context, text string, twitterId string) (string, error) {
	return func(ctx context.Context, text string, twitterId string) (string, error) {
		prompt := fmt.Sprintf("Answer the user's question about CoinGecko data:\n\n%s\n\nRespond concisely. User\\'s twitter_id is %s", text, twitterId)
		return mcpclient.Ask(ctx, mcpCmd, mcpTool, prompt, twitterId)
	}
}
