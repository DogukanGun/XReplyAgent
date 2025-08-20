package ask_ai_mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// AIMCPClient represents an HTTP client for the BNB Chain AI MCP server
type AIMCPClient struct {
	baseURL string
	client  *http.Client
}

// Tool represents an MCP tool
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ConnectAiMcpServer creates a new connection to the BNB Chain AI MCP server
func ConnectAiMcpServer() (*AIMCPClient, error) {
	return ConnectAiMcpServerWithURL("https://mcp.inkeep.com/bnbchainorg/mcp")
}

// ConnectAiMcpServerWithURL creates a new connection to an AI MCP server with a custom URL
func ConnectAiMcpServerWithURL(url string) (*AIMCPClient, error) {
	client := &AIMCPClient{
		baseURL: url,
		client:  &http.Client{Timeout: 60 * time.Second},
	}

	// Initialize the connection
	if err := client.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize MCP connection: %w", err)
	}

	return client, nil
}

// initialize sends the initialize request to the MCP server
func (c *AIMCPClient) initialize() error {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2025-06-18",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "bnb-ai-mcp-client",
				"version": "0.1.0",
			},
		},
	}

	_, err := c.postJSON(payload, nil)
	return err
}

// ListTools retrieves the list of available tools from the MCP server
func (c *AIMCPClient) ListTools(ctx context.Context) ([]Tool, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
	}

	var response struct {
		Result struct {
			Tools []Tool `json:"tools"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if _, err := c.postJSON(payload, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", response.Error.Code, response.Error.Message)
	}

	return response.Result.Tools, nil
}

// CallTool calls a specific tool on the MCP server with the given arguments
func (c *AIMCPClient) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (string, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": arguments,
		},
	}

	var response struct {
		Result struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if _, err := c.postJSON(payload, &response); err != nil {
		return "", err
	}

	if response.Error != nil {
		return "", fmt.Errorf("MCP error %d: %s", response.Error.Code, response.Error.Message)
	}

	if len(response.Result.Content) == 0 {
		return "", errors.New("no content returned from tool call")
	}

	// Concatenate all text content
	var result []string
	for _, content := range response.Result.Content {
		if content.Type == "text" && content.Text != "" {
			result = append(result, content.Text)
		}
	}

	if len(result) == 0 {
		return "", errors.New("no text content returned from tool call")
	}

	return fmt.Sprintf("%s", result[0]), nil
}

// Ask is a convenience method that calls a tool with a question parameter
func (c *AIMCPClient) Ask(ctx context.Context, toolName string, question string) (string, error) {
	return c.CallTool(ctx, toolName, map[string]interface{}{
		"question": question,
	})
}

// postJSON sends a JSON POST request to the MCP server
func (c *AIMCPClient) postJSON(payload interface{}, response interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if response != nil {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return resp, fmt.Errorf("failed to decode response: %w", err)
		}
	} else {
		resp.Body.Close()
	}

	return resp, nil
}

// ExampleUsage demonstrates how to use the AI MCP client
func ExampleUsage() error {
	ctx := context.Background()

	// Connect to the BNB Chain AI MCP server
	client, err := ConnectAiMcpServer()
	if err != nil {
		return fmt.Errorf("failed to connect to AI MCP server: %w", err)
	}

	// List available tools
	tools, err := client.ListTools(ctx)
	if err != nil {
		return fmt.Errorf("failed to list tools: %w", err)
	}

	fmt.Printf("Available tools (%d):\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// Example: Ask a question using the first available tool
	if len(tools) > 0 {
		toolName := tools[0].Name
		question := "What is BNB Chain?"

		answer, err := client.Ask(ctx, toolName, question)
		if err != nil {
			return fmt.Errorf("failed to ask question: %w", err)
		}

		fmt.Printf("\nQuestion: %s\n", question)
		fmt.Printf("Answer: %s\n", answer)
	}

	return nil
}
