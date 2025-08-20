package ask_ai_mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestConnectAiMcpServer tests the basic connection function
func TestConnectAiMcpServer(t *testing.T) {
	// Create a mock server that responds to initialize requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Parse the request body
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			return
		}

		// Check if it's an initialize request
		if req["method"] == "initialize" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"result": map[string]interface{}{
					"protocolVersion": "2025-06-18",
					"capabilities":    map[string]interface{}{},
					"serverInfo": map[string]interface{}{
						"name":    "test-server",
						"version": "1.0.0",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	// Test with the mock server URL
	client, err := ConnectAiMcpServerWithURL(server.URL)
	if err != nil {
		t.Errorf("Expected successful connection, got error: %v", err)
	}

	if client == nil {
		t.Error("Expected non-nil client")
		return
	}

	if client.baseURL != server.URL {
		t.Errorf("Expected baseURL to be '%s', got '%s'", server.URL, client.baseURL)
	}

	if client.client == nil {
		t.Error("Expected non-nil HTTP client")
	}
}

// TestConnectAiMcpServerWithURL tests the custom URL connection function
func TestConnectAiMcpServerWithURL(t *testing.T) {
	// Test with invalid URL (should fail)
	_, err := ConnectAiMcpServerWithURL("invalid-url")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}

	// Test with non-existent server (should fail)
	_, err = ConnectAiMcpServerWithURL("http://localhost:99999")
	if err == nil {
		t.Error("Expected error for non-existent server, got nil")
	}
}

// TestAIMCPClient_ListTools tests the ListTools method
func TestAIMCPClient_ListTools(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] == "initialize" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"result":  map[string]interface{}{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if req["method"] == "tools/list" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"result": map[string]interface{}{
					"tools": []map[string]interface{}{
						{
							"name":        "test_tool",
							"description": "A test tool",
							"inputSchema": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"question": map[string]interface{}{
										"type": "string",
									},
								},
							},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := ConnectAiMcpServerWithURL(server.URL)
	if err != nil {
		t.Fatalf("Failed to connect to mock server: %v", err)
	}

	ctx := context.Background()
	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Errorf("Expected successful tool listing, got error: %v", err)
	}

	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}

	if tools[0].Name != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got '%s'", tools[0].Name)
	}

	if tools[0].Description != "A test tool" {
		t.Errorf("Expected tool description 'A test tool', got '%s'", tools[0].Description)
	}
}

// TestAIMCPClient_ListTools_Error tests error handling in ListTools
func TestAIMCPClient_ListTools_Error(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] == "initialize" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"result":  map[string]interface{}{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if req["method"] == "tools/list" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"error": map[string]interface{}{
					"code":    -1,
					"message": "Test error",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := ConnectAiMcpServerWithURL(server.URL)
	if err != nil {
		t.Fatalf("Failed to connect to mock server: %v", err)
	}

	ctx := context.Background()
	_, err = client.ListTools(ctx)
	if err == nil {
		t.Error("Expected error from ListTools, got nil")
	}

	if !strings.Contains(err.Error(), "Test error") {
		t.Errorf("Expected error to contain 'Test error', got: %v", err)
	}
}

// TestAIMCPClient_CallTool tests the CallTool method
func TestAIMCPClient_CallTool(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] == "initialize" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"result":  map[string]interface{}{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if req["method"] == "tools/call" {
			params := req["params"].(map[string]interface{})
			toolName := params["name"].(string)

			if toolName == "test_tool" {
				response := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      req["id"],
					"result": map[string]interface{}{
						"content": []map[string]interface{}{
							{
								"type": "text",
								"text": "This is a test response",
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			} else {
				response := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      req["id"],
					"error": map[string]interface{}{
						"code":    -1,
						"message": "Unknown tool",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		}
	}))
	defer server.Close()

	client, err := ConnectAiMcpServerWithURL(server.URL)
	if err != nil {
		t.Fatalf("Failed to connect to mock server: %v", err)
	}

	ctx := context.Background()

	// Test successful tool call
	result, err := client.CallTool(ctx, "test_tool", map[string]interface{}{
		"question": "What is BNB Chain?",
	})
	if err != nil {
		t.Errorf("Expected successful tool call, got error: %v", err)
	}

	if result != "This is a test response" {
		t.Errorf("Expected 'This is a test response', got '%s'", result)
	}

	// Test tool call with unknown tool
	_, err = client.CallTool(ctx, "unknown_tool", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error for unknown tool, got nil")
	}
}

// TestAIMCPClient_Ask tests the Ask convenience method
func TestAIMCPClient_Ask(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] == "initialize" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"result":  map[string]interface{}{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if req["method"] == "tools/call" {
			params := req["params"].(map[string]interface{})
			arguments := params["arguments"].(map[string]interface{})

			// Check if question parameter was passed correctly
			if question, ok := arguments["question"].(string); ok && question == "Test question" {
				response := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      req["id"],
					"result": map[string]interface{}{
						"content": []map[string]interface{}{
							{
								"type": "text",
								"text": "Test answer",
							},
						},
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		}
	}))
	defer server.Close()

	client, err := ConnectAiMcpServerWithURL(server.URL)
	if err != nil {
		t.Fatalf("Failed to connect to mock server: %v", err)
	}

	ctx := context.Background()
	result, err := client.Ask(ctx, "test_tool", "Test question")
	if err != nil {
		t.Errorf("Expected successful ask, got error: %v", err)
	}

	if result != "Test answer" {
		t.Errorf("Expected 'Test answer', got '%s'", result)
	}
}

// TestAIMCPClient_CallTool_EmptyContent tests handling of empty content
func TestAIMCPClient_CallTool_EmptyContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)

		if req["method"] == "initialize" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"result":  map[string]interface{}{},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if req["method"] == "tools/call" {
			response := map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      req["id"],
				"result": map[string]interface{}{
					"content": []map[string]interface{}{},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	client, err := ConnectAiMcpServerWithURL(server.URL)
	if err != nil {
		t.Fatalf("Failed to connect to mock server: %v", err)
	}

	ctx := context.Background()
	_, err = client.CallTool(ctx, "test_tool", map[string]interface{}{})
	if err == nil {
		t.Error("Expected error for empty content, got nil")
	}

	if !strings.Contains(err.Error(), "no content returned") {
		t.Errorf("Expected error about no content, got: %v", err)
	}
}

// TestExampleUsage tests the example usage function (with mock server)
func TestExampleUsage(t *testing.T) {
	// This is more of a smoke test to ensure ExampleUsage doesn't panic
	// In a real scenario, it would try to connect to the actual server

	// We can't easily test this without mocking the global ConnectAiMcpServer function
	// So we'll just ensure the function exists and can be called
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ExampleUsage panicked: %v", r)
		}
	}()

	// This will likely fail due to network connectivity, but shouldn't panic
	_ = ExampleUsage()
}

// Benchmark tests
func BenchmarkConnectAiMcpServer(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  map[string]interface{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := ConnectAiMcpServerWithURL(server.URL)
		if err != nil {
			b.Errorf("Failed to connect: %v", err)
		}
		_ = client
	}
}

// Integration tests (these require the actual server to be available)
func TestIntegrationConnectAiMcpServer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client, err := ConnectAiMcpServer()
	if err != nil {
		t.Skipf("Skipping integration test due to connection failure: %v", err)
	}

	if client.baseURL != "https://mcp.inkeep.com/bnbchainorg/mcp" {
		t.Errorf("Expected correct base URL, got '%s'", client.baseURL)
	}
}

func TestIntegrationListTools(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client, err := ConnectAiMcpServer()
	if err != nil {
		t.Skipf("Skipping integration test due to connection failure: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tools, err := client.ListTools(ctx)
	if err != nil {
		t.Logf("Integration test failed (may be expected if server is unavailable): %v", err)
		return
	}

	t.Logf("Found %d tools", len(tools))
	for _, tool := range tools {
		t.Logf("Tool: %s - %s", tool.Name, tool.Description)
		if tool.Name == "" {
			t.Error("Tool name should not be empty")
		}
	}
}
