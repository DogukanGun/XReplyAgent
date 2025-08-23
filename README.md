## XReplyAgent 🤖

Automated X/Twitter mention responder powered by MCP servers and a LangChain agent. It reads mention events (e.g., from an n8n workflow), chooses the right tool (CoinGecko, GoldRush), generates an answer, and optionally posts a reply on X.

### Why
- **Faster support**: Answer crypto data questions at mention time.
- **Tool-orchestration**: Use multiple MCP servers (CoinGecko, GoldRush) via an agent.
- **Hands-free posting**: Reply under the original tweet through an X MCP.

### Hyperliquid extension goal 🧭
XReplyAgent showcases how the Hyperliquid ecosystem can be made more accessible and engaging directly on X:
- **Ask-on-X**: Community members can mention your account and ask about Hyperliquid topics (e.g., HLP price, token movements, transaction details), receiving answers in-reply, in seconds.
- **Right tool for the job**: The agent routes price/market questions to CoinGecko MCP and on‑chain activity questions (balances, transfers, gas, transaction lookups) to the GoldRush MCP.
- **Frictionless onboarding**: No dashboards or query builders—natural‑language questions on X become live, contextual responses, increasing Hyperliquid awareness and engagement.
- **Extensible**: New Hyperliquid‑specific tools can be added as MCP endpoints, and the agent will automatically consider them.

Examples the agent can answer for Hyperliquid users:
- “What’s the price of HLP right now?” → CoinGecko MCP.
- “Show the last 3 ERC20 transfers of the HLP contract on ethereum; brief summary.” → GoldRush MCP.
- “Give me the native balance and recent activity for 0x… on base; 1 line.” → GoldRush MCP.

---

## How it works 🔎


- The bot exposes `POST /mentions` and accepts either a single object or an array of objects (see Payloads).
- The bot runs in **agent mode**: each mention triggers the `agent` binary which selects and calls tools and can post via X.

---

## Components 🧩

- `cmd/mcp-servers/general/coingecko/xmcp` (X MCP)
  - Tool: `twitter.post_reply` (posts under a tweet). Requires X auth (Bearer or OAuth1).

- `cmd/mcp-servers/general/coingecko/cgproxy` (CoinGecko Proxy MCP)
  - Bridges HTTP MCP to an upstream stdio CoinGecko MCP (local via `npx @coingecko/coingecko-mcp` or remote via `mcp-remote`).
  - Discovers and forwards all upstream tools.

- `cmd/mcp-servers/general/goldrush` (GoldRush MCP)
  - Adds tools for on-chain data (Covalent GoldRush): balances, transactions, gas, NFTs, token holders, etc.
  - Requires `GOLDRUSH_AUTH_TOKEN` and often an allow-list (IP) on Covalent.

- `cmd/agent` (LangChainGo ReAct agent)
  - Discovers tools from HTTP MCP servers: CoinGecko + GoldRush + X poster.
  - Runs with OpenAI (or compatible) model. Sanitizes final answer to avoid tool chatter.
  - If `-reply-to` is provided, posts the answer via X MCP and exits non-zero on failures.

- `cmd/bot` (HTTP server)
  - Endpoint: `GET /healthz`, `POST /mentions`.
  - Always runs in **agent mode**: set `AGENT_CMD`; the bot spawns the agent per mention.

---

## n8n integration 🧰

### Payloads accepted by `POST /mentions`
- Single object (`MentionsPayload`) or array of objects.

Example (single):
```json
{
  "count": 1,
  "mentions": [
    {
      "tweet_id": "1958197635065463114",
      "text": "what’s the price of $HLP on Hyperliquid right now?",
      "author_id": "123",
      "author_username": "alice",
      "conversation_id": "1956374656836907309",
      "created_at": "2025-08-15T10:00:00.000Z"
    }
  ],
  "meta": {}
}
```

Example (array of payloads):
```json
[
  {
    "count": 2,
    "mentions": [ { /* ... */ }, { /* ... */ } ],
    "meta": { "result_count": 2 }
  }
]
```

---

## Run locally 🧪

Prereqs: Go 1.21+, Node for `npx` (if using local CoinGecko MCP), OpenAI API key, X credentials.

### 1) Start X MCP (twitter.post_reply) 🐦
```bash
go build -o xmcp ./cmd/mcp-servers/general/coingecko/xmcp
export X_BEARER_TOKEN="<user_bearer_token>"   # or use OAuth1 vars below
# OAuth1 (alternative)
export X_AUTH_MODE=oauth1
export X_CONSUMER_KEY="..."
export X_CONSUMER_SECRET="..."
export X_ACCESS_TOKEN="..."
export X_ACCESS_SECRET="..."
PORT=8081 ./xmcp
```

### 2) Start CoinGecko Proxy MCP 🦎
```bash
go build -o cgproxy ./cmd/mcp-servers/general/coingecko/cgproxy
export CG_MCP_CMD="npx"
export CG_MCP_ARGS="mcp-remote https://mcp.api.coingecko.com/sse"  # or: -y @coingecko/coingecko-mcp
PORT=8082 ./cgproxy
```

### 3) Start GoldRush MCP ⛓️
```bash
go build -o goldrush ./cmd/mcp-servers/general/goldrush
export GOLDRUSH_AUTH_TOKEN="cqt_..."
PORT=8083 ./goldrush
```
Tip: If Covalent enforces allow-lists, add the egress IP (Elastic IP or NAT EIP). Validate the token with:
```bash
curl -i -H "Authorization: Bearer $GOLDRUSH_AUTH_TOKEN" https://api.covalenthq.com/v1/chains/
```

### 4) Run bot in Agent mode (recommended) 🤖
```bash
go build -o agent ./cmd/agent
go build -o bot ./cmd/bot
export AGENT_CMD="$(pwd)/agent"
export AGENT_CG_MCP_HTTP="http://localhost:8082/mcp"
export AGENT_X_MCP_HTTP="http://localhost:8081/mcp"
export AGENT_GOLDRUSH_MCP_HTTP="http://localhost:8083/mcp"
export OPENAI_API_KEY="<your_openai_key>"
PORT=8080 ./bot
```

### 5) Test a mention 📬
```bash
curl -s -X POST http://localhost:8080/mentions \
  -H 'Content-Type: application/json' \
  -d '{
    "count":1,
    "mentions":[{"tweet_id":"1958197635065463114","text":"what’s the price of $HLP on Hyperliquid right now?","author_id":"123","author_username":"alice","conversation_id":"1956374656836907309","created_at":"2025-08-15T10:00:00.000Z"}],
    "meta":{}
  }'
```

Replies under a tweet are handled automatically by the agent when invoked by the bot; no extra flags are needed in normal operation.

---

## Environment variables 📦

### X MCP
- `X_BEARER_TOKEN` (or OAuth1: `X_AUTH_MODE=oauth1`, `X_CONSUMER_KEY`, `X_CONSUMER_SECRET`, `X_ACCESS_TOKEN`, `X_ACCESS_SECRET`)
- `PORT` (default 8081)

### CoinGecko Proxy MCP
- `CG_MCP_CMD` (e.g., `npx`)
- `CG_MCP_ARGS` (e.g., `-y @coingecko/coingecko-mcp` or `mcp-remote https://mcp.api.coingecko.com/sse`)
- `PORT` (default 8082)

### GoldRush MCP
- `GOLDRUSH_AUTH_TOKEN` (Covalent token; may require allow-list)
- `PORT` (default 8083)

### Agent
- `CG_MCP_HTTP` (e.g., `http://localhost:8082/mcp`)
- `X_MCP_HTTP`  (e.g., `http://localhost:8081/mcp`)
- `GOLDRUSH_MCP_HTTP` (e.g., `http://localhost:8083/mcp`)
- `OPENAI_API_KEY`, `OPENAI_MODEL` (default `gpt-4.1-mini`)
- Flags: `-q`, `-reply-to`

### Bot
- `AGENT_CMD`, `AGENT_CG_MCP_HTTP`, `AGENT_X_MCP_HTTP`, `AGENT_GOLDRUSH_MCP_HTTP`, `OPENAI_API_KEY`
- `WEBHOOK_SECRET` (optional), `PORT` (default 8080)

---

## GoldRush usage guidance 🧭
- Keep queries **lightweight**: include `chain_name`, a specific token contract when sensible, and small limits (e.g., last 3).
- Prefer single-chain endpoints to multichain when possible.
- Example prompts that steer the agent:
  - "GoldRush: last 3 ERC20 transfers for 0x742d… on ethereum for USDC 0xa0b8…e48; brief."
  - "GoldRush: native balance for 0x742d… on ethereum; 1 line."
  - "GoldRush: tx details for 0x<txhash> on base; concise."

---

## Troubleshooting 🛠️

### Port already in use
```bash
lsof -nP -iTCP:8081 -sTCP:LISTEN
kill <PID>            # or: pkill -f xmcp
```

### n8n 400 Bad Request
- Ensure Body is valid JSON and not a literal object value; always `JSON.stringify(...)` the full payload.
- Content-Type must be `application/json`.

### X posting succeeds but tweet not visible
- Verify token has `tweet.write` and you reply to a valid, visible tweet id.
- Agent now exits non-zero on posting failures; check the bot JSON response `error` field.

### GoldRush 401 / 402
- `401 Unauthorized`: bad/expired token or allow-list not set. Validate with curl (see above).
- `402 Payment Required`: credits/quota issue.

### Remove tool chatter from replies
- The agent sanitizes its final answer to remove lines like `Action:` / `Observation:`; update to latest build if you still see them.

---

## Development notes 🧑‍💻
- High-level tool discovery is done over HTTP MCP (`initialize`, `tools/list`, `tools/call`).
- The agent maps discovered tools to LangChainGo `tools.Tool` instances, and calls them by name.
- For CoinGecko, `cgproxy` forwards to an upstream stdio MCP (local or remote).

---

## Security 🔐
- Keep all API keys in env vars; avoid committing secrets.
- For Covalent GoldRush, use an Elastic IP or NAT Gateway EIP in the allow-list.
- The bot supports a `WEBHOOK_SECRET` header to protect `/mentions`.

---

## License 📄
MIT (see `LICENSE`).


