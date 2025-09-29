package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"cg-mentions-bot/internal/types"
)

// MentionsHandler handles POST /mentions events.
type MentionsHandler struct {
	Secret string
	Ask    func(ctx context.Context, text string, twitterId string) (string, error)
	Reply  func(ctx context.Context, in ReplyIn) error
	// If set, uses the agent binary to both answer and post per mention.
	AgentRun func(ctx context.Context, question string, replyTo string, twitterId string, mentionedPeople []string) (string, error)
}

// XIdResponse used for x id of the user
type XIdResponse struct {
	Data *struct {
		CreatedAt string `json:"created_at"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		Protected bool   `json:"protected"`
		Username  string `json:"username"`
	} `json:"data,omitempty"`
	Errors []struct {
		Detail string `json:"detail"`
		Status int    `json:"status"`
		Title  string `json:"title"`
		Type   string `json:"type"`
	} `json:"errors,omitempty"`
	Includes *struct {
		Users []struct {
			CreatedAt string `json:"created_at"`
			ID        string `json:"id"`
			Name      string `json:"name"`
			Protected bool   `json:"protected"`
			Username  string `json:"username"`
		} `json:"users,omitempty"`

		Tweets []struct {
			AuthorID string `json:"author_id"`
			ID       string `json:"id"`
			Text     string `json:"text"`
		} `json:"tweets,omitempty"`
	} `json:"includes,omitempty"`
}

// ReplyIn contains minimal info to reply to a tweet.
type ReplyIn struct {
	InReplyTo string
	Text      string
}

// Handle verifies secret (if configured), processes mentions, and returns a summary.
func (h MentionsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if h.Secret != "" && r.Header.Get("X-Webhook-Secret") != h.Secret {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Accept either a single payload object or an array of payloads
	var payload types.MentionsPayload
	var payloads []types.MentionsPayload
	if err := json.Unmarshal(body, &payloads); err != nil {
		// Not an array, try single
		if err2 := json.Unmarshal(body, &payload); err2 != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		payloads = []types.MentionsPayload{payload}
	}

	// Flatten mentions
	mentions := make([]types.Mention, 0)
	received := 0
	for _, p := range payloads {
		if p.Count > 0 {
			received += p.Count
		}
		if len(p.Mentions) > 0 {
			mentions = append(mentions, p.Mentions...)
		}
	}
	if received == 0 {
		received = len(mentions)
	}

	type res struct {
		TweetID string `json:"tweet_id"`
		Posted  bool   `json:"posted"`
		Error   string `json:"error,omitempty"`
	}

	results := make([]res, 0, len(mentions))

	for _, m := range mentions {
		q := normalizeTweetText(m.Text)
		mentionedUsers := handleMentions(m.Text)
		if h.AgentRun != nil {
			if _, err := h.AgentRun(r.Context(), q, m.TweetID, m.AuthorID, mentionedUsers); err != nil {
				results = append(results, res{TweetID: m.TweetID, Posted: false, Error: err.Error()})
			} else {
				results = append(results, res{TweetID: m.TweetID, Posted: true})
			}
			continue
		}

		ans, err := h.Ask(r.Context(), q, m.AuthorID)
		if err != nil {
			results = append(results, res{TweetID: m.TweetID, Posted: false, Error: err.Error()})
			continue
		}

		if postErr := h.Reply(r.Context(), ReplyIn{InReplyTo: m.TweetID, Text: ans}); postErr != nil {
			results = append(results, res{TweetID: m.TweetID, Posted: false, Error: postErr.Error()})
			continue
		}

		results = append(results, res{TweetID: m.TweetID, Posted: true})
	}

	summary := struct {
		Received  int   `json:"received"`
		Processed int   `json:"processed"`
		Results   []res `json:"results"`
	}{
		Received:  received,
		Processed: len(mentions),
		Results:   results,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(summary)
}

// handleMentions returns mentioned users from a tweet
func handleMentions(s string) []string {
	// Regex: @ followed by word characters (letters, numbers, underscore)
	re := regexp.MustCompile(`@(\w+)`)

	// Find all matches (only capture group 1, without "@")
	matches := re.FindAllStringSubmatch(s, -1)

	var users []string
	for _, match := range matches {
		if len(match) > 1 {
			xId := getXId(match[1])
			users = append(users, xId)
		}
	}

	return users
}

// normalizeTweetText removes handles and URLs and trims whitespace to form a concise question input.
func normalizeTweetText(s string) string {
	// Remove URLs
	urlRe := regexp.MustCompile(`https?://\S+`)
	s = urlRe.ReplaceAllString(s, " ")
	// Remove @handles
	handleRe := regexp.MustCompile(`@[A-Za-z0-9_]+`)
	s = handleRe.ReplaceAllString(s, " ")
	// Collapse whitespace
	s = strings.TrimSpace(strings.Join(strings.Fields(s), " "))
	return s
}

func getXId(username string) string {
	url := "https://api.x.com/2/users/by/username/" + username
	log.Printf("url is : %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for %s: %s", url, err)
		return ""
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("XAUTH_TOKEN"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return ""
	}
	if res == nil {
		log.Println("response is nil")
		return ""
	}
	defer func() {
		if res.Body != nil {
			if err := res.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}
	}()

	// Check HTTP status code
	if res.StatusCode != http.StatusOK {
		log.Printf("Error fetching user id: %s", res.Status)
		return ""
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return ""
	}

	var resp XIdResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Printf("Error unmarshaling X API response for user %s: %v", username, err)
		return ""
	}

	// Check if resp.Data is nil before accessing its fields
	if resp.Data == nil {
		log.Printf("No user data found for username: %s", username)
		return ""
	}

	// Check if ID is empty
	if resp.Data.ID == "" {
		log.Printf("Empty user ID for username: %s", username)
		return ""
	}

	println(resp.Data.ID)

	return resp.Data.ID
}
