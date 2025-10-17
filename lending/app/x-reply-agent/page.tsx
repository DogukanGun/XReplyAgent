import Image from "next/image";

export default function XReplyAgentPage() {
  return (
    <div className="font-sans grid grid-rows-[min-content_1fr_min-content] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
      <main className="flex flex-col gap-[32px] row-start-2 items-center sm:items-start max-w-4xl mx-auto">
        <div className="flex items-center gap-4 mb-4 w-full justify-center sm:justify-start">
          <div className="flex items-center justify-center w-16 h-16 bg-blue-500 rounded-lg">
            <Image src="/globe.svg" alt="X Reply Agent Logo" width={32} height={32} className="invert" />
          </div>
          <div className="flex flex-col">
            <h1 className="text-4xl font-bold">X Reply Agent</h1>
            <div className="flex gap-4 mt-2">
              <a href="https://github.com/DogukanGun/XReplyAgent" target="_blank" rel="noopener noreferrer" className="text-gray-600 hover:text-black dark:text-gray-400 dark:hover:text-white flex items-center gap-2">
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                GitHub
              </a>
            </div>
          </div>
        </div>


        <p className="text-lg text-gray-700 dark:text-gray-300 mb-6">
          Our project is an AI agent that uses n8n to catch mentions on X and
          automatically replies with real-time Hyperliquid insights. By
          integrating CoinGecko MCP and Hyperliquid’s GoldRush MCP, it delivers
          price, funding rate, and open interest data directly into
          conversations. This solves the problem of fragmented, delayed
          information by bringing HL-native analytics to the platforms where the
          community already engages such as X. It strengthens Hyperliquid’s
          visibility and culture while giving traders instant, reliable answers.
          <br />
          <a
            href="https://taikai.network/hl-hackathon-organizers/hackathons/hl-hackathon/projects/cmejaxqeh023aw40ispuniili/idea"
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-500 hover:underline"
          >
            Learn more on TAIKAI
          </a>
        </p>

        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Team</h2>
          <ul className="list-disc list-inside text-gray-700 dark:text-gray-300">
            <li>Cem Denizsel (Owner)</li>
            <li>Dogukan Ali Gundogan (Member)</li>
          </ul>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Categories</h2>
          <ul className="list-disc list-inside text-gray-700 dark:text-gray-300">
            <li>02. 🚀 Hyperliquid Frontier Track</li>
            <li>17. Best use of GoldRush</li>
          </ul>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Why</h2>
          <ul className="list-disc list-inside text-gray-700 dark:text-gray-300">
            <li>
              <strong>Faster support</strong>: Answer crypto data questions at
              mention time.
            </li>
            <li>
              <strong>Tool-orchestration</strong>: Use multiple MCP servers
              (CoinGecko, GoldRush) via an agent.
            </li>
            <li>
              <strong>Hands-free posting</strong>: Reply under the original
              tweet through an X MCP.
            </li>
          </ul>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">
            Hyperliquid extension goal 🧭
          </h2>
          <ul className="list-disc list-inside text-gray-700 dark:text-gray-300">
            <li>
              <strong>Ask-on-X</strong>: Community members can mention your
              account and ask about Hyperliquid topics (e.g., HLP price, token
              movements, transaction details), receiving answers in-reply, in
              seconds.
            </li>
            <li>
              <strong>Right tool for the job</strong>: The agent routes
              price/market questions to CoinGecko MCP and on‑chain activity
              questions (balances, transfers, gas, transaction lookups) to the
              GoldRush MCP.
            </li>
            <li>
              <strong>Frictionless onboarding</strong>: No dashboards or query
              builders—natural‑language questions on X become live, contextual
              responses, increasing Hyperliquid awareness and engagement.
            </li>
            <li>
              <strong>Extensible</strong>: New Hyperliquid‑specific tools can be
              added as MCP endpoints, and the agent will automatically consider
              them.
            </li>
          </ul>
          <p className="mt-4 text-gray-700 dark:text-gray-300">
            Examples the agent can answer for Hyperliquid users:
          </p>
          <ul className="list-disc list-inside ml-5 text-gray-700 dark:text-gray-300">
            <li>“What’s the price of HLP right now?” → CoinGecko MCP.</li>
            <li>
              “Show the last 3 ERC20 transfers of the HLP contract on ethereum;
              brief summary.” → GoldRush MCP.
            </li>
            <li>
              “Give me the native balance and recent activity for 0x… on base; 1
              line.” → GoldRush MCP.
            </li>
          </ul>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Components 🧩</h2>
          <ul className="list-disc list-inside text-gray-700 dark:text-gray-300">
            <li>
              <code>cmd/mcp-servers/general/coingecko/xmcp</code> (X MCP)
              <ul>
                <li>Tool: &quot;twitter.post_reply&quot; (posts under a tweet). Requires X auth (Bearer or OAuth1).</li>
              </ul>
            </li>
            <li>
              <code>cmd/mcp-servers/general/coingecko/cgproxy</code> (CoinGecko Proxy MCP)
              <ul>
                <li>Bridges HTTP MCP to an upstream stdio CoinGecko MCP (local via &quot;npx @coingecko/coingecko-mcp&quot; or remote via &quot;mcp-remote&quot;).</li>
                <li>Discovers and forwards all upstream tools.</li>
              </ul>
            </li>
            <li>
              <code>cmd/mcp-servers/general/goldrush</code> (GoldRush MCP)
              <ul>
                <li>Adds tools for on-chain data (Covalent GoldRush): balances, transactions, gas, NFTs, token holders, etc.</li>
                <li>Requires &quot;GOLDRUSH_AUTH_TOKEN&quot; and often an allow-list (IP) on Covalent.</li>
              </ul>
            </li>
            <li>
              <code>cmd/agent</code> (LangChainGo ReAct agent)
              <ul>
                <li>Discovers tools from HTTP MCP servers: CoinGecko + GoldRush + X poster.</li>
                <li>Runs with OpenAI (or compatible) model. Sanitizes final answer to avoid tool chatter.</li>
                <li>If `-reply-to` is provided, posts the answer via X MCP and exits non-zero on failures.</li>
              </ul>
            </li>
            <li>
              <code>cmd/bot</code> (HTTP server)
              <ul>
                <li>Endpoint: `GET /healthz`, `POST /mentions`.</li>
                <li>Always runs in <strong>agent mode</strong>: set `AGENT_CMD`; the bot spawns the agent per mention.</li>
              </ul>
            </li>
          </ul>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">n8n integration 🧰</h2>
          <h3 className="text-xl font-medium mb-2">
            Payloads accepted by &quot;POST /mentions&quot;
          </h3>
          <ul className="list-disc list-inside text-gray-700 dark:text-gray-300">
            <li>Single object (&quot;MentionsPayload&quot;) or array of objects.</li>
          </ul>
          <h4 className="text-lg font-medium mt-4 mb-2">Example (single):</h4>
          <pre className="bg-gray-100 dark:bg-gray-800 p-4 rounded-md text-sm overflow-x-auto">
            <code className="language-json">
              {`{
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
}`}
            </code>
          </pre>
        </section>

        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Quick Setup 🚀</h2>
          <div className="bg-gray-100 dark:bg-gray-800 p-6 rounded-lg">
            <p className="text-gray-700 dark:text-gray-300 mb-4">To get started with X Reply Agent:</p>
            <ol className="list-decimal list-inside space-y-2 text-gray-700 dark:text-gray-300">
              <li>Clone the repository:
                <pre className="whitespace-pre-wrap break-words bg-gray-200 dark:bg-gray-700 p-2 mt-1 rounded"><code>git clone https://github.com/DogukanGun/XReplyAgent.git</code></pre>
              </li>
              <li>Set up environment variables:
                <pre className="whitespace-pre-wrap break-words bg-gray-200 dark:bg-gray-700 p-2 mt-1 rounded"><code>export AGENT_CMD=&quot;$(pwd)/agent&quot;
export AGENT_CG_MCP_HTTP=&quot;http://localhost:8082/mcp&quot;
export X_MCP_HTTP=&quot;http://localhost:8081/mcp&quot;
export AGENT_GOLDRUSH_MCP_HTTP=&quot;http://localhost:8083/mcp&quot;
export OPENAI_API_KEY=&quot;your_openai_key&quot;</code></pre>
              </li>
              <li>Install dependencies and build the project</li>
              <li>Start the bot server and required MCP servers</li>
            </ol>
            <p className="mt-4 text-gray-700 dark:text-gray-300">For detailed setup instructions and configuration options, visit our <a href="https://github.com/DogukanGun/XReplyAgent" target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">GitHub repository</a>.</p>
          </div>
        </section>
      </main>
    </div>
  );
} 