import * as dotenv from "dotenv";
import { createLangchainTools, KeypairWallet, SolanaAgentKit } from "solana-agent-kit";
import bs58 from "bs58";
import { Keypair } from "@solana/web3.js";
import { z } from "zod"
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { withTwitterAuth } from "./twitter";
import BlinksPlugin from "@solana-agent-kit/plugin-blinks";
import DefiPlugin from "@solana-agent-kit/plugin-defi";
import MiscPlugin from "@solana-agent-kit/plugin-misc";
import NFTPlugin from "@solana-agent-kit/plugin-nft";
import TokenPlugin from "@solana-agent-kit/plugin-token";
import { HumanMessage } from "@langchain/core/messages";
import { MemorySaver } from "@langchain/langgraph";
import { createReactAgent } from "@langchain/langgraph/prebuilt";
import { ChatOpenAI } from "@langchain/openai";
import { mcpToolRes } from "./helper";

export async function startServer() {
    try {
        // Create a new MCP server instance with extended timeout
        const server = new McpServer({
            name: "BNBChain MCP Server",
            version: "1.0.0"
        }, {
            enforceStrictCapabilities: false
        })

        server.registerTool(
            "solana_agent",
            {
                inputSchema: {
                    twitter_id: z.string().describe("User\\'s twitter id"),
                    message: z.string().describe("The thing that user wanna do")
                },
            },
            withTwitterAuth(async ({ privateKey, message }) => {

                const keyPair = Keypair.fromSecretKey(
                    bs58.decode(privateKey as string),
                );
                const wallet = new KeypairWallet(keyPair, process.env.RPC_URL as string);
                const agent = new SolanaAgentKit(wallet, process.env.RPC_URL!, {
                    OPENAI_API_KEY: process.env.OPENAI_API_KEY,
                })
                    .use(TokenPlugin)
                const tools = createLangchainTools(agent, agent.actions);
                const memory = new MemorySaver();
                const llm = new ChatOpenAI({
                    model: "gpt-4o",
                    temperature: 0.3,
                    apiKey: process.env.OPENAI_API_KEY,
                    timeout: 50000,
                    maxRetries: 1,
                });
                console.log("tools are created")
                llm.bindTools(tools as any);
                const langchainAgent = createReactAgent({
                    llm,
                    tools: tools as any,
                    checkpointSaver: memory,
                    messageModifier: `You are a helpful agent that can interact onchain using the Solana Agent Kit. You are
                    empowered to interact onchain using your tools. If you need funds you can request it from the user and provide your wallet details. If there is a 5XX
                    (internal) HTTP error code, ask the user to try again later. If someone asks you to do something you
                    can't do with your currently available tools, you must say so, and encourage them to implement it
                    themselves using the Solana Agent Kit, recommend they go to https://www.solanaagentkit.xyz for more information. Be
                    concise and helpful with your responses. Refrain from restating your tools' descriptions unless it is explicitly requested.`,
                });
                console.log("langchainAgent is created")
                const response = await langchainAgent.invoke(
                    { messages: [new HumanMessage(message)] },
                    { configurable: { thread_id: "default" } }
                );
                console.log("response is created")
                
                // Extract the text content from the last message
                const messages = response.messages || [];
                const lastMessage = messages[messages.length - 1];
                const textContent = lastMessage?.content || "No response generated";
                
                return mcpToolRes.success({ response: textContent })

            })

        )
        // Register all resources, tools, and prompts

        return server
    } catch (error) {
        process.exit(1)
    }
}
