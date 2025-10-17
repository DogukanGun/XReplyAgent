package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	agentcore "cg-mentions-bot/internal/agentcore"
)

// InputRequest defines the expected JSON body
type InputRequest struct {
	Input     string `json:"input" binding:"required"`
	TwitterId string `json:"twitter_id" binding:"required"`
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

		answer, _ := agentcore.Ask(c.Request.Context(), req.Input, req.TwitterId)

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

		xURL := os.Getenv("X_MCP_HTTP")
		walletMcpUrl := os.Getenv("WALLET_MCP_HTTP")
		bnbHttpURL := os.Getenv("BNB_MCP_HTTP")
		if xURL == "" {
			fmt.Fprintln(os.Stderr, "Set X_MCP_HTTP (e.g., http://localhost:8081/mcp)")
			os.Exit(1)
		}

		question := flag.String("q", "", "question to ask the agent (fallback: AGENT_INPUT or stdin)")
		replyTo := flag.String("reply-to", "", "tweet id to reply under using x_post_reply (optional)")
		twitterId := flag.String("ti", "", "twitter id of the user that posts it")
		mentionedUser := flag.String("m", "", "mentioned user that is in the tweet")
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

		// Build agentcore config from env
		cfg := agentcore.Config{
			XMCP:      xURL,
			WalletMCP: walletMcpUrl,
			BNBMCP:    bnbHttpURL,
			Model:     os.Getenv("OPENAI_MODEL"),
		}

		// Create wallet for the user if twitter_id is provided (preserve behavior)
		fmt.Println("Users twitter id", *twitterId)
		if strings.TrimSpace(*twitterId) != "" {
			if err := agentcore.CreateWalletForTwitterIDWithConfig(context.Background(), strings.TrimSpace(*twitterId), cfg); err != nil {
				fmt.Fprintln(os.Stderr, "failed to create wallet:", err)
			}
		}

		// Ask using the same prompt logic (reply or non-reply)
		ctx := context.Background()
		answer, err := agentcore.AskAgent(ctx, q, strings.TrimSpace(*twitterId), strings.TrimSpace(*replyTo), strings.TrimSpace(*mentionedUser), cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		// If reply target provided, truncate then post via X
		if strings.TrimSpace(*replyTo) != "" {
			if answer == "" {
				fmt.Fprintln(os.Stderr, "agent produced empty answer; cannot post to X")
				os.Exit(1)
			}
			runes := []rune(answer)
			const maxTweetLen = 270
			if len(runes) > maxTweetLen {
				answer = string(runes[:maxTweetLen]) + "…"
			}
			if err := agentcore.TweetOnly(ctx, strings.TrimSpace(*replyTo), answer, cfg); err != nil {
				fmt.Fprintln(os.Stderr, "failed to post via X:", err)
				os.Exit(1)
			}
		}

		fmt.Println(answer)
	}
}
