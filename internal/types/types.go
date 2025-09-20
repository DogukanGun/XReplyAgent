package types

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

// Mention represents one mention payload item.
type Mention struct {
	TweetID        string `json:"tweet_id"`
	Text           string `json:"text"`
	AuthorID       string `json:"twitter_id"`
	AuthorUsername string `json:"author_username"`
	ConversationID string `json:"conversation_id"`
	CreatedAt      string `json:"created_at"`
}

// MentionsPayload is the full body we receive from n8n.
type MentionsPayload struct {
	Count    int            `json:"count"`
	Mentions []Mention      `json:"mentions"`
	Meta     map[string]any `json:"meta,omitempty"`
}

// ToolInfo contains a tool and its handler function
type ToolInfo struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
