package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"cg-mentions-bot/internal/agentcore"
)

type AgentAskRequest struct {
	Input        string   `json:"input"`
	ReplyTo      string   `json:"reply_to,omitempty"`
	MentionedIDs []string `json:"mentioned_ids,omitempty"`
}

type AgentAskResponse struct {
	Output string `json:"output"`
}

// AgentAskHandler runs the LangChain agent with provided input and context
//
//	@Summary		Ask the agent
//	@Description	Invoke the agent with input, twitter_id (derived), optional reply_to and mentioned users
//	@Tags			agent
//	@Accept			json
//	@Produce		json
//	@Param			request body AgentAskRequest true "Agent Ask Request"
//	@Success		200 {object} AgentAskResponse
//	@Failure		400 {object} ErrorResponse
//	@Failure		500 {object} ErrorResponse
//	@Security		Bearer
//	@Router			/agent/ask [post]
func AgentAskHandler(w http.ResponseWriter, r *http.Request) {
	var req AgentAskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON in request body"})
		return
	}

	if req.Input == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "input is required"})
		return
	}

	// Derive twitter_id from Firebase UID in context
	uid, _ := r.Context().Value(UidKey).(string)
	var twitterID string
	if uid != "" {
		if user, ok := GetUserByFirebaseID(uid); ok && user != nil {
			twitterID = user.TwitterID
		}
	}
	if twitterID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "twitter_id not found for user"})
		return
	}

	cfg := agentcore.Config{
		XMCP:      os.Getenv("X_MCP_HTTP"),
		WalletMCP: os.Getenv("WALLET_MCP_HTTP"),
		BNBMCP:    os.Getenv("BNB_MCP_HTTP"),
		SolanaMCP: os.Getenv("SOLANA_MCP_HTTP"),
		Model:     os.Getenv("OPENAI_MODEL"),
	}
	agentcore.CreateWalletForTwitterIDWithConfig(r.Context(), twitterID, cfg)
	out, err := agentcore.AskAgent(r.Context(), req.Input, twitterID, "", "", cfg)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(AgentAskResponse{Output: out})
}
