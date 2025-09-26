import * as dotenv from "dotenv";
import { Aptos, AptosConfig, Ed25519PrivateKey, HexInput, Network, PrivateKey, PrivateKeyVariants } from "@aptos-labs/ts-sdk";
import { MongoClient } from "mongodb"
import { AgentRuntime, createAptosTools, LocalSigner } from "move-agent-kit";
import { ChatOpenAI } from "@langchain/openai";
import { MemorySaver } from "@langchain/langgraph";
import { createReactAgent } from "@langchain/langgraph/prebuilt"
import { HumanMessage } from "@langchain/core/messages";

dotenv.config();

export type RunAptosAgentResult = {
	messages: string[]
	final: string
}

export async function runAptosAgent(message: string, twitterId: string): Promise<RunAptosAgentResult> {
	const outputs: string[] = []

	// aptos setup
	const aptosConfig = new AptosConfig({
		network: Network.TESTNET,
	})
	const aptos = new Aptos(aptosConfig)

	let privateKeyHex: string | undefined
	if (twitterId) {
		const user = await getUserByTwitterId(twitterId)
		if (user?.private_key) {
			privateKeyHex = user.private_key
		}
	}

	if (!privateKeyHex) {
		throw new Error("No private key available. Provide via twitter_id lookup.")
	}

	const account = await aptos.deriveAccountFromPrivateKey({
		privateKey: new Ed25519PrivateKey(
			PrivateKey.formatPrivateKey(privateKeyHex as HexInput, PrivateKeyVariants.Ed25519)
		),
	})

	const signer = new LocalSigner(account, Network.TESTNET)
	const agentRuntime = new AgentRuntime(signer, aptos, {
		PANORA_API_KEY: process.env.PANORA_API_KEY,
	})
	const tools = createAptosTools(agentRuntime)

	const llm = new ChatOpenAI({
		temperature: 0.7,
		model: "gpt-4.1-mini",
		openAIApiKey: process.env.OPENAI_API_KEY,
	});

	const memory = new MemorySaver();

	const agent = createReactAgent({
		llm,
		tools,
		checkpointSaver: memory,
		messageModifier: `
		You are a helpful agent that can interact onchain using the Aptos Agent Kit. You are
		empowered to interact onchain using your tools. Whatever is asked, you must respond with a natural language response. So answer questions for Aptos
	`,
	});
	const config = { configurable: { thread_id: "Aptos Agent Kit!" } }
	const stream = await agent.stream(
		{
			messages: [new HumanMessage(message)],
		},
		config
	);

	let final = ""
	for await (const chunk of stream) {
		if ("agent" in chunk) {
			const text = String(chunk.agent.messages[0].content)
			outputs.push(text)
			final = text
		} else if ("tools" in chunk) {
			const text = String(chunk.tools.messages[0].content)
			outputs.push(text)
		}
	}

	return { messages: outputs, final }
}

async function getUserByTwitterId(twitterId: string): Promise<{ twitter_id: string, public_key: string, private_key: string } | null> {
	try {
		const client = await getMongoClient()
		const db = client.db("User")
		const collection = db.collection("Wallet")

		const user = await collection.findOne<{
			twitter_id: string;
			public_key: string;
			private_key: string;
			username?: string
		}>({ twitter_id: twitterId })

		if (!user) {
			return null
		}

		return {
			twitter_id: user.twitter_id,
			public_key: user.public_key,
			private_key: user.private_key,
		}
	} catch (error) {
		console.log("Error fetching user by twitter_id:", error)
		return null
	}
}

async function getMongoClient(): Promise<MongoClient> {
	const uri = process.env.MONGO_URI || "mongodb://localhost:27017"
	const client = new MongoClient(uri)
	await client.connect()
	return client
}