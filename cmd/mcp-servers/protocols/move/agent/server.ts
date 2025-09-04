import "dotenv/config"
import express, { type Request, type Response } from "express"
import cors from "cors"
import { runAptosAgent } from "./index.js"

const app = express()
app.use(cors())
app.use(express.json())

// Tool schema following MCP format
const APTOS_TOOL = {
	name: "aptos_handle",
	description: "Handle Aptos on-chain operations. Provide message and twitter_id.",
	inputSchema: {
		type: "object",
		properties: {
			message: { type: "string", description: "Natural language instruction for Aptos agent" },
			twitter_id: { type: "string", description: "User's twitter_id for wallet lookup" },
		},
		required: ["message"],
	},
}

app.get("/", (_: Request, res: Response) => res.send("Aptos MCP HTTP Server running"))

// Expose tools metadata in MCP format
app.get("/mcp/tools", (_: Request, res: Response) => {
	res.json({ tools: [APTOS_TOOL] })
})

// Execute the aptos_handle tool
app.post("/mcp", async (req: Request, res: Response) => {
	try {
		const { message, twitter_id } = req.body || {}
		if (!message || typeof message !== "string") {
			return res.status(400).json({ error: "'message' (string) is required" })
		}
		const result = await runAptosAgent(message, typeof twitter_id === "string" ? twitter_id : undefined)
		res.json({
			name: "aptos_handle",
			content: [{ type: "text", text: result.final }],
		})
	} catch (err: any) {
		res.status(500).json({ error: err?.message || String(err) })
	}
})

const PORT = Number(process.env.HTTP_PORT || 3010)
app.listen(PORT, () => {
	console.log(`Aptos MCP HTTP server listening on http://localhost:${PORT}`)
}) 