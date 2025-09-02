package bnb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/tmc/langchaingo/tools"
	"log"
	"os"
	"strings"
	"time"
)

type Agent struct {
	mcpClient *client.Client
}

func ProxyHandler() (Agent, error) {
	ctx := context.Background()
	baseURL := os.Getenv("BNB_AGENT_MCP_SSE")
	if baseURL == "" {
		return Agent{}, fmt.Errorf("BNB_AGENT_MCP_SSE environment variable not set")
	}

	mcpClient, err := client.NewSSEMCPClient(baseURL)
	if err != nil {
		return Agent{}, fmt.Errorf("error creating mcp client: %w", err)
	}

	err = mcpClient.Start(ctx)
	if err != nil {
		return Agent{}, fmt.Errorf("error starting mcp client: %w", err)
	}

	initParams := mcp.InitializeParams{
		ProtocolVersion: "2024-11-05",
		Capabilities:    mcp.ClientCapabilities{},
		ClientInfo: mcp.Implementation{
			Name:    "bnb-proxy",
			Version: "1.0.0",
		},
	}

	if _, err = mcpClient.Initialize(ctx, mcp.InitializeRequest{
		Params: initParams,
	}); err != nil {
		return Agent{}, fmt.Errorf("error initializing mcp client: %w", err)
	}

	log.Println("Bnb proxy started")
	return Agent{mcpClient: mcpClient}, nil
}

func (b *Agent) GetToolsInfo(ctx context.Context) (string, error) {
	log.Println("Bnb proxy tools")
	lgTools, err := b.mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return "", fmt.Errorf("failed to list tools: %w", err)
	}

	var descriptions []string
	descriptions = append(descriptions, "BNB Chain MCP Server Proxy - This tool provides access to the following BNB Chain operations:")

	for _, tool := range lgTools.Tools {
		toolDesc := fmt.Sprintf("\n- %s: %s", tool.Name, tool.Description)

		schemaBytes, err := json.Marshal(tool.InputSchema)
		if err == nil {
			var schemaMap map[string]interface{}
			if err := json.Unmarshal(schemaBytes, &schemaMap); err == nil {
				if properties, ok := schemaMap["properties"].(map[string]interface{}); ok {
					var params []string
					for propName, propDef := range properties {
						if propMap, ok := propDef.(map[string]interface{}); ok {
							propType := "string"
							if t, exists := propMap["type"]; exists {
								propType = fmt.Sprintf("%v", t)
							}
							propDescription := ""
							if desc, exists := propMap["description"]; exists {
								propDescription = fmt.Sprintf(" - %v", desc)
							}
							params = append(params, fmt.Sprintf("%s (%s)%s", propName, propType, propDescription))
						}
					}
					if len(params) > 0 {
						toolDesc += fmt.Sprintf("\n  Parameters: %s", strings.Join(params, ", "))
					}
				}
			}
		}
		descriptions = append(descriptions, toolDesc)
	}

	descriptions = append(descriptions, "\nAnswer all kind of question about opBnb or bsc chain by using thsese tools")

	return strings.Join(descriptions, ""), nil
}

func (b *Agent) CallTool(ctx context.Context, toolName string, parameters map[string]interface{}) (interface{}, error) {
	callRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: parameters,
		},
	}

	result, err := b.mcpClient.CallTool(ctx, callRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to call tool %s: %w", toolName, err)
	}

	return result, nil
}

func (b *Agent) GetAvailableTools(ctx context.Context) ([]string, error) {
	tools, err := b.mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	var toolNames []string
	for _, tool := range tools.Tools {
		toolNames = append(toolNames, tool.Name)
	}

	return toolNames, nil
}

type ProxyTool struct {
	agent *Agent
}

func (b *Agent) AsLangChainTool() tools.Tool {
	return &ProxyTool{agent: b}
}

func (t *ProxyTool) Name() string {
	return "bnb_chain_proxy"
}

func (t *ProxyTool) Description() string {
	ctx := context.Background()
	lgTools, err := t.agent.mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Printf("Error getting tools info: %v", err)
		return "BNB Chain MCP Server Proxy - Error retrieving tool descriptions"
	}

	var descriptions []string
	descriptions = append(descriptions, "BNB Chain MCP Server Proxy - Use this tool to interact with BNB Chain. You must provide a JSON input with 'tool_name' and 'parameters' fields.")
	descriptions = append(descriptions, "\nAvailable tools:")

	for _, tool := range lgTools.Tools {
		toolDesc := fmt.Sprintf("\n- Tool name: '%s' - %s", tool.Name, tool.Description)

		// Parse and display parameters more clearly
		schemaBytes, err := json.Marshal(tool.InputSchema)
		if err == nil {
			var schemaMap map[string]interface{}
			if err := json.Unmarshal(schemaBytes, &schemaMap); err == nil {
				if properties, ok := schemaMap["properties"].(map[string]interface{}); ok {
					var params []string
					required := []string{}
					if req, ok := schemaMap["required"].([]interface{}); ok {
						for _, r := range req {
							if reqStr, ok := r.(string); ok {
								required = append(required, reqStr)
							}
						}
					}

					for propName, propDef := range properties {
						if propMap, ok := propDef.(map[string]interface{}); ok {
							propType := "string"
							if t, exists := propMap["type"]; exists {
								propType = fmt.Sprintf("%v", t)
							}

							isRequired := false
							for _, req := range required {
								if req == propName {
									isRequired = true
									break
								}
							}

							requiredStr := ""
							if isRequired {
								requiredStr = " (required)"
							}

							propDescription := ""
							if desc, exists := propMap["description"]; exists {
								propDescription = fmt.Sprintf(" - %v", desc)
							}
							params = append(params, fmt.Sprintf("%s (%s)%s%s", propName, propType, requiredStr, propDescription))
						}
					}
					if len(params) > 0 {
						toolDesc += fmt.Sprintf("\n  Parameters: %s", strings.Join(params, ", "))
					}
				}
			}
		}
		descriptions = append(descriptions, toolDesc)
	}

	descriptions = append(descriptions, "\n\nUsage: Provide input as JSON: {\"tool_name\": \"exact_tool_name\", \"parameters\": {\"param1\": \"value1\", \"param2\": \"value2\"}}")

	return strings.Join(descriptions, "")
}

type CallRequest struct {
	ToolName   string                 `json:"tool_name"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

func (t *ProxyTool) Call(ctx context.Context, input string) (string, error) {
	// Parse the input to extract tool name and parameters
	var callRequest CallRequest
	if err := json.Unmarshal([]byte(input), &callRequest); err != nil {
		return "", fmt.Errorf("failed to parse input JSON. Expected format: {\"tool_name\": \"...\", \"parameters\": {...}}: %w", err)
	}

	if callRequest.ToolName == "" {
		return "", fmt.Errorf("tool_name is required in input JSON")
	}

	// Create a context with longer timeout for BNB operations
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 360*time.Second)
	defer cancel()

	result, err := t.agent.CallTool(ctxWithTimeout, callRequest.ToolName, callRequest.Parameters)
	if err != nil {
		return "", err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultBytes), nil
}
