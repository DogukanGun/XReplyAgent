package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
)

type mcpHTTP struct {
	base string
	hc   *http.Client
}

// InputRequest defines the expected JSON body
type InputRequest struct {
	Input     string `json:"input" binding:"required"`
	TwitterId string `json:"twitter_id" binding:"required"`
}

func newMCP(base string) *mcpHTTP {
	c := &http.Client{Timeout: 60 * time.Second}
	_ = post(c, base, map[string]any{
		"jsonrpc": "2.0", "id": 1, "method": "initialize",
		"params": map[string]any{
			"protocolVersion": "2025-06-18",
			"capabilities":    map[string]any{},
			"clientInfo":      map[string]any{"name": "lc", "version": "0.1"},
		},
	})
	return &mcpHTTP{base: base, hc: c}
}

func (m *mcpHTTP) call(name string, args map[string]any) (string, error) {
	var out struct {
		Result struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
			//IsError bool `json:"isError"`
		} `json:"result"`
	}
	if _, err := postJSON(m.hc, m.base, map[string]any{
		"jsonrpc": "2.0", "id": 2, "method": "tools/call",
		"params": map[string]any{"name": name, "arguments": args},
	}, &out); err != nil {
		return "", err
	}
	// if out.Result.IsError {
	// 	if len(out.Result.Content) > 0 && strings.TrimSpace(out.Result.Content[0].Text) != "" {
	// 		return "", fmt.Errorf(out.Result.Content[0].Text)
	// 	}
	// 	return "", fmt.Errorf("tool call failed: %s", name)
	// }
	if len(out.Result.Content) == 0 {
		return "", nil
	}
	return out.Result.Content[0].Text, nil
}

func (m *mcpHTTP) listTools() ([]map[string]any, error) {
	var out struct {
		Result struct {
			Tools []map[string]any `json:"tools"`
		} `json:"result"`
	}
	if _, err := postJSON(m.hc, m.base, map[string]any{
		"jsonrpc": "2.0", "id": 3, "method": "tools/list",
	}, &out); err != nil {
		return nil, err
	}
	return out.Result.Tools, nil
}

func post(c *http.Client, url string, body any) error {
	_, err := postJSON(c, url, body, nil)
	return err
}

func postJSON(c *http.Client, url string, body any, out any) (*http.Response, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if out == nil {
			resp.Body.Close()
		}
	}()
	if out != nil {
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

type genericMCPTool struct {
	client *mcpHTTP
	name   string
	desc   string
}

func (t genericMCPTool) Name() string        { return t.name }
func (t genericMCPTool) Description() string { return t.desc }
func (t genericMCPTool) Call(ctx context.Context, input string) (string, error) {
	var a map[string]any
	_ = json.Unmarshal([]byte(input), &a)
	return t.client.call(t.name, a)
}

type xTool struct{ client *mcpHTTP }

func (t xTool) Name() string { return "x_post_reply" }
func (t xTool) Description() string {
	return "Reply under a tweet via X MCP. Input JSON: {\"in_reply_to_tweet_id\":\"...\",\"text\":\"...\"}"
}
func (t xTool) Call(ctx context.Context, input string) (string, error) {
	var a map[string]any
	_ = json.Unmarshal([]byte(input), &a)
	return t.client.call("twitter.post_reply", a)
}

// func cgDiscoveredTools(cg *mcpHTTP) ([]tools.Tool, error) {
// 	raw, err := cg.listTools()
// 	if err != nil {
// 		return nil, err
// 	}
// 	out := make([]tools.Tool, 0, len(raw))
// 	for _, t := range raw {
// 		name, _ := t["name"].(string)
// 		if name == "" {
// 			continue
// 		}
// 		description, _ := t["description"].(string)
// 		// include inputSchema (if any) as a compact JSON to guide the LLM
// 		// if schemaVal, ok := t["inputSchema"]; ok && schemaVal != nil {
// 		// 	if b, err := json.Marshal(schemaVal); err == nil {
// 		// 		description = fmt.Sprintf("%s\nInput JSON must match schema: %s", description, string(b))
// 		// 	}
// 		// }
// 		out = append(out, genericMCPTool{client: cg, name: name, desc: description})
// 	}
// 	return out, nil
// }

// grDiscoveredTools discovers tools exposed by the GoldRush HTTP MCP server and
// converts them to LangChainGo tools for the agent.
func grDiscoveredTools(gr *mcpHTTP) ([]tools.Tool, error) {
	raw, err := gr.listTools()
	if err != nil {
		return nil, err
	}
	out := make([]tools.Tool, 0, len(raw))
	for _, t := range raw {
		name, _ := t["name"].(string)
		if name == "" {
			continue
		}
		description, _ := t["description"].(string)
		// include inputSchema (if any) as a compact JSON to guide the LLM
		if schemaVal, ok := t["inputSchema"]; ok && schemaVal != nil {
			if b, err := json.Marshal(schemaVal); err == nil {
				description = fmt.Sprintf("%s\nInput JSON must match schema: %s", description, string(b))
			}
		}
		out = append(out, genericMCPTool{client: gr, name: name, desc: description})
	}
	return out, nil
}

// wlDiscoveredTools discovers tools exposed by the Wallet HTTP MCP server and
// includes inputSchema details to guide the LLM on argument structure.
func wlDiscoveredTools(wl *mcpHTTP) ([]tools.Tool, error) {
	raw, err := wl.listTools()
	if err != nil {
		return nil, err
	}
	out := make([]tools.Tool, 0, len(raw))
	for _, t := range raw {
		name, _ := t["name"].(string)
		if name == "" {
			continue
		}
		description, _ := t["description"].(string)
		if schemaVal, ok := t["inputSchema"]; ok && schemaVal != nil {
			if b, err := json.Marshal(schemaVal); err == nil {
				description = fmt.Sprintf("%s\nInput JSON must match schema: %s", description, string(b))
			}
		}
		out = append(out, genericMCPTool{client: wl, name: name, desc: description})
	}
	return out, nil
}

// bnbDiscoveredTools discovers tools exposed by the BNB HTTP MCP proxy server
// and converts them to LangChainGo tools for the agent.
func bnbDiscoveredTools(bnb *mcpHTTP) ([]tools.Tool, error) {
	raw, err := bnb.listTools()
	if err != nil {
		return nil, err
	}
	out := make([]tools.Tool, 0, len(raw))
	for _, t := range raw {
		name, _ := t["name"].(string)
		if name == "" {
			continue
		}
		description, _ := t["description"].(string)
		// include inputSchema (if any) as a compact JSON to guide the LLM
		if schemaVal, ok := t["inputSchema"]; ok && schemaVal != nil {
			if b, err := json.Marshal(schemaVal); err == nil {
				description = fmt.Sprintf("%s\nInput JSON must match schema: %s", description, string(b))
			}
		}
		out = append(out, genericMCPTool{client: bnb, name: name, desc: description})
	}
	return out, nil
}

func askAgentAndGetXMcp(question string, twitterId string) (string, *mcpHTTP) {
	// cgURL := os.Getenv("CG_MCP_HTTP") // CoinGecko disabled to reduce token usage
	xURL := os.Getenv("X_MCP_HTTP")
	// goldrushURL := os.Getenv("GOLDRUSH_MCP_HTTP") // GoldRush disabled for now
	walletMcpUrl := os.Getenv("WALLET_MCP_HTTP")
	bnbHttpURL := os.Getenv("BNB_MCP_HTTP")
	if xURL == "" {
		fmt.Fprintln(os.Stderr, "Set X_MCP_HTTP (e.g., http://localhost:8081/mcp)")
		return "", nil
	}

	q := strings.TrimSpace(question)
	if q == "" {
		if v := strings.TrimSpace(os.Getenv("AGENT_INPUT")); v != "" {
			q = v
		} else if fi, _ := os.Stdin.Stat(); fi != nil && (fi.Mode()&os.ModeCharDevice) == 0 {
			b, _ := io.ReadAll(os.Stdin)
			q = strings.TrimSpace(string(b))
		}
	}
	if q == "" {
		fmt.Fprintln(os.Stderr, "Provide a question with -q, AGENT_INPUT, or piped stdin.")
		return "", nil
	}

	// cg disabled to reduce token usage
	// cg := newMCP(cgURL)
	x := newMCP(xURL)
	// gr := newMCP(goldrushURL) // GoldRush disabled
	wl := newMCP(walletMcpUrl)

	// Initialize BNB tools: prefer HTTP MCP proxy; fallback to SSE proxy tool
	var bnbTools []tools.Tool
	if strings.TrimSpace(bnbHttpURL) != "" {
		if t, err := bnbDiscoveredTools(newMCP(bnbHttpURL)); err == nil {
			bnbTools = t
		} else {
			fmt.Fprintln(os.Stderr, "failed to discover BNB HTTP tools:", err)
		}
	}

	// cgTools, err := cgDiscoveredTools(cg)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "failed to discover CG tools:", err)
	// }
	// grTools, err := grDiscoveredTools(gr)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, "failed to discover GoldRush tools:", err)
	// }
	wlTools, err := wlDiscoveredTools(wl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to discover Wallet Mcp tools:", err)
	}

	toolsList := make([]tools.Tool, 0, len(wlTools)+len(bnbTools)+2)
	toolsList = append(toolsList, xTool{client: x})
	toolsList = append(toolsList, bnbTools...)
	toolsList = append(toolsList, wlTools...)
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4.1-mini"
	}
	llm, err := openai.New(openai.WithModel(model))
	if err != nil {
		fmt.Fprintln(os.Stderr, "OPENAI_API_KEY is required or configure a provider supported by LangChainGo")
		return "", nil
	}

	peh := agents.NewParserErrorHandler(nil)
	exec, err := agents.Initialize(
		llm,
		toolsList,
		agents.ZeroShotReactDescription,
		agents.WithMaxIterations(10),
		agents.WithParserErrorHandler(peh),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", nil
	}

	prompt := fmt.Sprintf("%s Answer this question using the available MCP tools. The twitter id of the user is: %s ", q, twitterId)

	// Debug: print the exact LLM prompt to stderr (not returned to clients)
	fmt.Fprintln(os.Stderr, "LLM prompt:", prompt)
	ctx := context.Background()
	log.Printf("Asking the question: %s", prompt)
	out, err := exec.Call(ctx, map[string]any{"input": prompt})
	if err != nil {
		// Print best effort output or last observation, but exit non-zero so callers can detect failure
		if steps, ok := out["intermediateSteps"].([]schema.AgentStep); ok && len(steps) > 0 {
			fmt.Println(steps[len(steps)-1].Observation)
		}
		if v, ok := out["output"].(string); ok && v != "" {
			fmt.Println(v)
		}
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Extract and sanitize final answer
	answer, _ := out["output"].(string)
	answer = sanitizeFinalAnswer(answer)
	return answer, x
}

func provideTester() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Status(200)
	})

	// POST /receive - accepts JSON { "input": "some string" }
	r.POST("/receive", func(c *gin.Context) {
		var req InputRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			// Binding failed (missing or wrong type)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body, expected JSON { \"input\": " +
				"\"string\" }"})
			return
		}

		answer, _ := askAgentAndGetXMcp(req.Input, req.TwitterId)

		c.JSON(http.StatusOK, gin.H{
			"response": answer,
			"length":   len(answer),
		})
	})

	// Start server on port 8080
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func main() {
	testModeOn := os.Getenv("TEST_MODE")
	if testModeOn == "YES" {
		provideTester()
	} else {
		// cgURL := os.Getenv("CG_MCP_HTTP") // CoinGecko disabled
		xURL := os.Getenv("X_MCP_HTTP")
		// goldrushURL := os.Getenv("GOLDRUSH_MCP_HTTP") // GoldRush disabled
		walletMcpUrl := os.Getenv("WALLET_MCP_HTTP")
		bnbHttpURL := os.Getenv("BNB_MCP_HTTP")
		if xURL == "" {
			fmt.Fprintln(os.Stderr, "Set X_MCP_HTTP (e.g., http://localhost:8081/mcp)")
			os.Exit(1)
		}

		question := flag.String("q", "", "question to ask the agent (fallback: AGENT_INPUT or stdin)")
		replyTo := flag.String("reply-to", "", "tweet id to reply under using x_post_reply (optional)")
		twitterId := flag.String("ti", "", "twitter id of the user that posts it")
		flag.Parse()

		q := strings.TrimSpace(*question)
		if q == "" {
			if v := strings.TrimSpace(os.Getenv("AGENT_INPUT")); v != "" {
				q = v
			} else if fi, _ := os.Stdin.Stat(); fi != nil && (fi.Mode()&os.ModeCharDevice) == 0 {
				b, _ := io.ReadAll(os.Stdin)
				q = strings.TrimSpace(string(b))
			}
		}
		if q == "" {
			fmt.Fprintln(os.Stderr, "Provide a question with -q, AGENT_INPUT, or piped stdin.")
			os.Exit(1)
		}

		// cg disabled to reduce token usage
		// cg := newMCP(cgURL)
		x := newMCP(xURL)
		// gr := newMCP(goldrushURL)
		wl := newMCP(walletMcpUrl)

		// Initialize BNB tools: prefer HTTP MCP proxy; fallback to SSE proxy tool
		var bnbTools []tools.Tool
		if strings.TrimSpace(bnbHttpURL) != "" {
			if t, err := bnbDiscoveredTools(newMCP(bnbHttpURL)); err == nil {
				bnbTools = t
			} else {
				fmt.Fprintln(os.Stderr, "failed to discover BNB HTTP tools:", err)
			}
		}

		// cgTools, err := cgDiscoveredTools(cg)
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, "failed to discover CG tools:", err)
		// 	os.Exit(1)
		// }
		// grTools, err := grDiscoveredTools(gr)
		// if err != nil {
		// 	fmt.Fprintln(os.Stderr, "failed to discover GoldRush tools:", err)
		// 	os.Exit(1)
		// }
		wlTools, err := wlDiscoveredTools(wl)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to discover Wallet Mcp tools:", err)
		}

		toolsList := make([]tools.Tool, 0, len(bnbTools)+len(wlTools)+2)
		toolsList = append(toolsList, xTool{client: x})
		toolsList = append(toolsList, bnbTools...)
		toolsList = append(toolsList, wlTools...)
		model := os.Getenv("OPENAI_MODEL")
		if model == "" {
			model = "gpt-4.1-mini"
		}
		llm, err := openai.New(openai.WithModel(model))
		if err != nil {
			fmt.Fprintln(os.Stderr, "OPENAI_API_KEY is required or configure a provider supported by LangChainGo")
			os.Exit(1)
		}

		peh := agents.NewParserErrorHandler(nil)
		exec, err := agents.Initialize(
			llm,
			toolsList,
			agents.ZeroShotReactDescription,
			agents.WithMaxIterations(10),
			agents.WithParserErrorHandler(peh),
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Create wallet for the user if twitter_id is provided and wallet MCP is available
		fmt.Println("Users twitter id", *twitterId)
		if strings.TrimSpace(*twitterId) != "" {
			resFromCreateWallet, err := wl.call("create_wallet", map[string]any{
				"twitter_id": strings.TrimSpace(*twitterId),
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to create wallet:", err)
			}
			fmt.Println("here is the result of create new wallet: ", resFromCreateWallet)
		}

		prompt := q
		prompt = fmt.Sprintf("%s . User\\'s twitter_id is %s", prompt, strings.TrimSpace(*twitterId))
		if strings.TrimSpace(*replyTo) != "" {
			prompt = fmt.Sprintf("%s Answer this question using the available MCP tools. You are an AI agent that manages user wallets via tweet commands. "+
				"Your reply will be posted on X; write concise, user-facing text. "+
				"Never share private keys or the twitter_id in the reply. Then reply to tweet %s using x_post_reply. Also user\\'s twitter_id is %s. "+
				"If a blockchain transaction is executed (e.g., a transfer), include its transaction hash; for wallet creation or reads, provide the wallet address.",
				prompt, strings.TrimSpace(*replyTo), strings.TrimSpace(*twitterId))
		}

		// Debug: print the exact LLM prompt to stderr (not returned to clients)
		fmt.Fprintln(os.Stderr, "LLM prompt:", prompt)
		ctx := context.Background()
		out, err := exec.Call(ctx, map[string]any{"input": prompt})
		if err != nil {
			// Print best effort output or last observation, but exit non-zero so callers can detect failure
			if steps, ok := out["intermediateSteps"].([]schema.AgentStep); ok && len(steps) > 0 {
				fmt.Println(steps[len(steps)-1].Observation)
				os.Exit(1)
			}
			if v, ok := out["output"].(string); ok && v != "" {
				fmt.Println(v)
				os.Exit(1)
			}
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// Extract and sanitize final answer
		answer, _ := out["output"].(string)
		answer = sanitizeFinalAnswer(answer)

		// If a reply target was given, post the agent's answer using the X MCP directly.
		if strings.TrimSpace(*replyTo) != "" {
			if answer == "" {
				fmt.Fprintln(os.Stderr, "agent produced empty answer; cannot post to X")
				os.Exit(1)
			}
			// best-effort truncate to fit X limits
			runes := []rune(answer)
			const maxTweetLen = 270
			if len(runes) > maxTweetLen {
				answer = string(runes[:maxTweetLen]) + "…"
			}
			if _, postErr := x.call("twitter.post_reply", map[string]any{
				"in_reply_to_tweet_id": strings.TrimSpace(*replyTo),
				"text":                 answer,
			}); postErr != nil {
				fmt.Fprintln(os.Stderr, "failed to post via X:", postErr)
				os.Exit(1)
			}
		}

		fmt.Println(answer)
	}
}

// sanitizeFinalAnswer removes ReAct/tool-calling artifacts so only the answer remains
func sanitizeFinalAnswer(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	// Remove common prefixes
	s = strings.TrimPrefix(s, "Final Answer:")
	s = strings.TrimSpace(s)
	// Remove obvious tool/action lines
	re := regexp.MustCompile(`(?i)^(action|action input|observation|thought|tool|intermediate steps)\s*:`)
	lines := strings.Split(s, "\n")
	filtered := make([]string, 0, len(lines))
	for _, ln := range lines {
		if re.MatchString(strings.TrimSpace(ln)) {
			continue
		}
		filtered = append(filtered, ln)
	}
	s = strings.TrimSpace(strings.Join(filtered, "\n"))
	s = strings.ReplaceAll(s, "```", "")
	return strings.TrimSpace(s)
}
