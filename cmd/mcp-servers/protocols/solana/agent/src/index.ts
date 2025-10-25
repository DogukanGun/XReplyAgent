import express, { Request, Response } from "express"
import cors from "cors"
import dotenv from "dotenv"
import { startServer } from "./server"
import { SSEServerTransport } from "@modelcontextprotocol/sdk/server/sse.js"

// Load environment variables
dotenv.config()

async function main() {
    const app = express()
    const server = await startServer()
    app.use(cors())

    // Log the current log level on startup
    console.log(`Starting sse server`)

    // to support multiple simultaneous connections we have a lookup object from
    // sessionId to transport
    const transports: { [sessionId: string]: SSEServerTransport } = {}

    app.get("/sse", async (_: Request, res: Response) => {
      const transport = new SSEServerTransport("/messages", res)
      transports[transport.sessionId] = transport
      console.log("New SSE connection established", {
        sessionId: transport.sessionId
      })

      res.on("close", () => {
        console.log("SSE connection closed", { sessionId: transport.sessionId })
        delete transports[transport.sessionId]
      })

      try {
        await server.connect(transport)
      } catch (error) {
        console.error("Error connecting transport", {
          sessionId: transport.sessionId,
          error
        })
      }
    })

    app.post("/messages", async (req: Request, res: Response) => {
      const sessionId = req.query.sessionId as string
      const transport = transports[sessionId]
      
      if (!transport) {
        res.status(404).json({ error: "Session not found" })
        return
      }

      try {
        await transport.handlePostMessage(req, res)
      } catch (error) {
        console.error("Error handling message", { sessionId, error })
        res.status(500).json({ error: "Internal server error" })
      }
    })

    const PORT = process.env.PORT || 3000
    app.listen(PORT, () => {
      console.log(`MCP Server listening on port ${PORT}`)
    })
}

main()