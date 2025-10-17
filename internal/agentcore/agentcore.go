package agentcore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

// Config carries endpoints and LLM model, supplied by caller (e.g., cmd/agent)
type Config struct {
	XMCP      string
	WalletMCP string
	BNBMCP    string
	SolanaMCP string
	Model     string
}

type mcpHTTP struct {
	base string
	hc   *http.Client
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
		} `json:"result"`
	}
	if _, err := postJSON(m.hc, m.base, map[string]any{
		"jsonrpc": "2.0", "id": 2, "method": "tools/call",
		"params": map[string]any{"name": name, "arguments": args},
	}, &out); err != nil {
		return "", err
	}
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

// xTool mirrors cmd/agent behavior to expose twitter.post_reply as a tool
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
		if schemaVal, ok := t["inputSchema"]; ok && schemaVal != nil {
			if b, err := json.Marshal(schemaVal); err == nil {
				description = fmt.Sprintf("%s\nInput JSON must match schema: %s", description, string(b))
			}
		}
		out = append(out, genericMCPTool{client: bnb, name: name, desc: description})
	}
	return out, nil
}

// Ask returns env-based AskWithConfig
func Ask(ctx context.Context, input string, twitterID string) (string, error) {
	cfg := Config{
		XMCP:      os.Getenv("X_MCP_HTTP"),
		WalletMCP: os.Getenv("WALLET_MCP_HTTP"),
		BNBMCP:    os.Getenv("BNB_MCP_HTTP"),
		SolanaMCP: os.Getenv("SOLANA_MCP_HTTP"),
		Model:     os.Getenv("OPENAI_MODEL"),
	}
	return AskAgent(ctx, input, twitterID, "", "", cfg)
}

// Tweet env-based wrapper
func Tweet(ctx context.Context, replyTo string, text string) error {
	return TweetOnly(ctx, replyTo, text, Config{XMCP: os.Getenv("X_MCP_HTTP")})
}

// CreateWalletForTwitterIDWithConfig ensures a wallet exists
func CreateWalletForTwitterIDWithConfig(ctx context.Context, twitterID string, cfg Config) error {
	if strings.TrimSpace(twitterID) == "" {
		return fmt.Errorf("twitter_id is required")
	}
	if strings.TrimSpace(cfg.WalletMCP) == "" {
		return nil
	}
	wl := newMCP(cfg.WalletMCP)
	_, err := wl.call("create_wallet", map[string]any{"twitter_id": strings.TrimSpace(twitterID)})
	return err
}

// CreateWalletForTwitterID env-based wrapper
func CreateWalletForTwitterID(ctx context.Context, twitterID string) error {
	return CreateWalletForTwitterIDWithConfig(ctx, twitterID, Config{WalletMCP: os.Getenv("WALLET_MCP_HTTP")})
}

func sanitizeFinalAnswer(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	s = strings.TrimPrefix(s, "Final Answer:")
	s = strings.TrimSpace(s)
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

// AskAgent: unified ask that can handle reply/non-reply prompts (no posting)
func AskAgent(ctx context.Context, question string, twitterID string, replyTo string, mentionedUser string, cfg Config) (string, error) {
	q := strings.TrimSpace(question)
	if q == "" || strings.TrimSpace(twitterID) == "" {
		return "", fmt.Errorf("input and twitter_id are required")
	}
	var toolsList []tools.Tool
	if strings.TrimSpace(cfg.XMCP) != "" {
		x := newMCP(cfg.XMCP)
		toolsList = append(toolsList, xTool{client: x})
	}
	if strings.TrimSpace(cfg.BNBMCP) != "" {
		if t, err := bnbDiscoveredTools(newMCP(cfg.BNBMCP)); err == nil {
			toolsList = append(toolsList, t...)
		} else {
			log.Println("failed to discover BNB HTTP tools:", err)
		}
	}
	// Discover Solana tools (optional)
	if strings.TrimSpace(cfg.SolanaMCP) != "" {
		if t, err := bnbDiscoveredTools(newMCP(cfg.SolanaMCP)); err == nil {
			toolsList = append(toolsList, t...)
		} else {
			log.Println("failed to discover Solana HTTP tools:", err)
		}
	}
	if strings.TrimSpace(cfg.WalletMCP) != "" {
		if t, err := wlDiscoveredTools(newMCP(cfg.WalletMCP)); err == nil {
			toolsList = append(toolsList, t...)
		} else {
			log.Println("failed to discover Wallet MCP tools:", err)
		}
	}
	model := cfg.Model
	if model == "" {
		model = "gpt-4.1-mini"
	}
	llm, err := openai.New(openai.WithModel(model))
	if err != nil {
		return "", fmt.Errorf("OPENAI_API_KEY is required or configure a provider supported by LangChainGo")
	}
	exec, err := agents.Initialize(
		llm,
		toolsList,
		agents.ZeroShotReactDescription,
		agents.WithMaxIterations(20),
		agents.WithParserErrorHandler(agents.NewParserErrorHandler(nil)),
	)
	if err != nil {
		return "", err
	}
	prompt := q
	if strings.TrimSpace(replyTo) != "" {
		prompt = fmt.Sprintf("%s Answer this question using the available MCP tools. You are an AI agent that manages user wallets via tweet commands. "+
			"Your reply will be posted on X; write concise, user-facing text. "+
			"Never share private keys or the twitter_id in the reply. Then reply to tweet %s using x_post_reply. Also user\\'s twitter_id is %s. "+
			"If a blockchain transaction is executed (e.g., a transfer), include its transaction hash; for wallet creation or reads, provide the wallet address."+
			"Also if user mentions another user like meaning to transfer something to another users, i am giving you possible user that is in tweet. If the input is empty, act like user has not mention"+
			"anybody. And in this condition always use Kanalabs mcp server or tool. Here is the id: %s",
			prompt, strings.TrimSpace(replyTo), strings.TrimSpace(twitterID), strings.TrimSpace(mentionedUser))
	} else {
		prompt = fmt.Sprintf("%s Answer this question using the available MCP tools. The twitter id of the user is: %s ", q, twitterID)
	}
	out, callErr := exec.Call(ctx, map[string]any{"input": prompt})
	if callErr != nil {
		if v, ok := out["output"].(string); ok && v != "" {
			return v, callErr
		}
		return "", callErr
	}
	ans, _ := out["output"].(string)
	return sanitizeFinalAnswer(ans), nil
}

// TweetOnly: post without modification so caller can truncate/format as desired
func TweetOnly(ctx context.Context, replyTo string, text string, cfg Config) error {
	if strings.TrimSpace(replyTo) == "" || strings.TrimSpace(text) == "" {
		return fmt.Errorf("reply_to and text are required")
	}
	if strings.TrimSpace(cfg.XMCP) == "" {
		return fmt.Errorf("X_MCP_HTTP is required")
	}
	x := newMCP(cfg.XMCP)
	if _, err := x.call("twitter.post_reply", map[string]any{
		"in_reply_to_tweet_id": strings.TrimSpace(replyTo),
		"text":                 text,
	}); err != nil {
		return err
	}
	return nil
}
