package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tmc/langchaingo/tools"
)

// MCPServer represents a connection to an MCP server
type MCPServer struct {
	Name    string
	BaseURL string
	client  *http.Client
	timeout time.Duration
}

// MCPToolCreator manages multiple MCP servers and creates tools from them
type MCPToolCreator struct {
	servers map[string]*MCPServer
}

// NewMCPToolCreator creates a new MCP tool creator
func NewMCPToolCreator() *MCPToolCreator {
	return &MCPToolCreator{
		servers: make(map[string]*MCPServer),
	}
}

// AddServer adds an MCP server to the creator
func (mtc *MCPToolCreator) AddServer(name, baseURL string, timeout time.Duration) error {
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	server := &MCPServer{
		Name:    name,
		BaseURL: baseURL,
		client:  &http.Client{Timeout: timeout},
		timeout: timeout,
	}

	// Initialize the MCP server
	if err := server.initialize(); err != nil {
		return fmt.Errorf("failed to initialize MCP server %s: %w", name, err)
	}

	mtc.servers[name] = server
	return nil
}

// GetAllTools discovers and returns all tools from all connected MCP servers
func (mtc *MCPToolCreator) GetAllTools() ([]tools.Tool, error) {
	var allTools []tools.Tool

	for serverName, server := range mtc.servers {
		tools, err := server.discoverTools()
		if err != nil {
			return nil, fmt.Errorf("failed to discover tools from server %s: %w", serverName, err)
		}
		allTools = append(allTools, tools...)
	}

	return allTools, nil
}

// GetToolsFromServer returns tools from a specific MCP server
func (mtc *MCPToolCreator) GetToolsFromServer(serverName string) ([]tools.Tool, error) {
	server, exists := mtc.servers[serverName]
	if !exists {
		return nil, fmt.Errorf("server %s not found", serverName)
	}

	return server.discoverTools()
}

// ListServers returns the names of all connected servers
func (mtc *MCPToolCreator) ListServers() []string {
	var names []string
	for name := range mtc.servers {
		names = append(names, name)
	}
	return names
}

// initialize sends the initialize request to the MCP server
func (server *MCPServer) initialize() error {
	initRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2025-06-18",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    fmt.Sprintf("langchain-go-client-%s", server.Name),
				"version": "0.1.0",
			},
		},
	}

	return server.post(initRequest, nil)
}

// discoverTools fetches all available tools from the MCP server
func (server *MCPServer) discoverTools() ([]tools.Tool, error) {
	listRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
	}

	var response struct {
		Result struct {
			Tools []map[string]interface{} `json:"tools"`
		} `json:"result"`
		Error *struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		} `json:"error,omitempty"`
	}

	if err := server.post(listRequest, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("MCP error %d: %s", response.Error.Code, response.Error.Message)
	}

	var tools []tools.Tool
	for _, toolData := range response.Result.Tools {
		tool := server.createGenericTool(toolData)
		if tool != nil {
			tools = append(tools, tool)
		}
	}

	return tools, nil
}

// createGenericTool creates a generic MCP tool from tool metadata
func (server *MCPServer) createGenericTool(toolData map[string]interface{}) tools.Tool {
	name, ok := toolData["name"].(string)
	if !ok || name == "" {
		return nil
	}

	description, _ := toolData["description"].(string)

	// Include input schema in description to guide the LLM
	if schemaVal, exists := toolData["inputSchema"]; exists && schemaVal != nil {
		if schemaBytes, err := json.Marshal(schemaVal); err == nil {
			description = fmt.Sprintf("%s\n\nInput JSON schema: %s", description, string(schemaBytes))
		}
	}

	return &GenericMCPTool{
		server:      server,
		toolName:    name,
		description: description,
	}
}

// callTool invokes a specific tool on the MCP server
func (server *MCPServer) callTool(toolName string, arguments map[string]interface{}) (string, error) {
	callRequest := map[string]interface{}{
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
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Data    interface{} `json:"data,omitempty"`
		} `json:"error,omitempty"`
	}

	if err := server.post(callRequest, &response); err != nil {
		return "", err
	}

	if response.Error != nil {
		return "", fmt.Errorf("MCP tool call error %d: %s", response.Error.Code, response.Error.Message)
	}

	if len(response.Result.Content) == 0 {
		return "", nil
	}

	// Concatenate all content parts
	var result string
	for _, content := range response.Result.Content {
		if content.Type == "text" {
			result += content.Text
		}
	}

	return result, nil
}

// post sends a JSON-RPC request to the MCP server
func (server *MCPServer) post(request interface{}, response interface{}) error {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, server.BaseURL, bytes.NewReader(requestBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := server.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// GenericMCPTool implements the tools.Tool interface for MCP tools
type GenericMCPTool struct {
	server      *MCPServer
	toolName    string
	description string
}

// Name returns the tool name
func (t *GenericMCPTool) Name() string {
	// Prefix with server name to avoid conflicts
	return fmt.Sprintf("%s_%s", t.server.Name, t.toolName)
}

// Description returns the tool description
func (t *GenericMCPTool) Description() string {
	return fmt.Sprintf("[%s] %s", t.server.Name, t.description)
}

// Call invokes the tool with the given input
func (t *GenericMCPTool) Call(ctx context.Context, input string) (string, error) {
	var arguments map[string]interface{}

	// Try to parse input as JSON
	if err := json.Unmarshal([]byte(input), &arguments); err != nil {
		// If parsing fails, treat as a simple string argument
		arguments = map[string]interface{}{
			"input": input,
		}
	}

	return t.server.callTool(t.toolName, arguments)
}

// Example usage and helper functions

// CreateMCPToolsFromConfig creates tools from a configuration
type MCPServerConfig struct {
	Name    string        `json:"name"`
	URL     string        `json:"url"`
	Timeout time.Duration `json:"timeout,omitempty"`
}

// SetupMCPTools creates and configures MCP tools from server configurations
func SetupMCPTools(configs []MCPServerConfig) ([]tools.Tool, error) {
	creator := NewMCPToolCreator()

	// Add all servers
	for _, config := range configs {
		if err := creator.AddServer(config.Name, config.URL, config.Timeout); err != nil {
			return nil, fmt.Errorf("failed to add server %s: %w", config.Name, err)
		}
	}

	// Get all tools
	return creator.GetAllTools()
}
