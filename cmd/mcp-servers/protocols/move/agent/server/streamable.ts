import "dotenv/config"

import { StreamableHTTPServerTransport } from "@modelcontextprotocol/sdk/server/streamableHttp.js"
import cors from "cors"
import express from "express"
import type { Request, Response } from "express"

import { startServer } from "./base.js"

export const startStreamableServer = async () => {
  try {
    const app = express()
    const server = startServer()
    app.use(cors())
    app.use(express.json())

    console.log("Starting Streamable HTTP server")

    // Store transports for session management
    const transports: { [sessionId: string]: StreamableHTTPServerTransport } = {}

    // Root endpoint for debugging
    app.get("/", (_, res) => {
      res.json({ 
        message: "Aptos MCP Server",
        transport: "Streamable HTTP",
        endpoints: {
          mcp: "/mcp"
        }
      })
    })

    // Modern Streamable HTTP endpoint - following official SDK pattern
    app.all('/mcp', async (req: Request, res: Response) => {
      const sessionId = req.query.sessionId as string || `session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
      
      let transport = transports[sessionId]
      if (!transport) {
        transport = new StreamableHTTPServerTransport({
          sessionIdGenerator: () => sessionId
        })
        transports[sessionId] = transport
        
        console.log("New Streamable HTTP connection established", { sessionId })
        
        res.on("close", () => {
          console.log("Streamable HTTP connection closed", { sessionId })
          delete transports[sessionId]
        })
        
        try {
          await server.connect(transport)
        } catch (error) {
          console.error("Error connecting transport", { sessionId, error })
        }
      }
      
      try {
        await transport.handleRequest(req, res)
      } catch (error) {
        console.error("Error handling request", { sessionId, error })
        res.status(500).send("Internal server error")
      }
    })

    const PORT = Number(process.env.HTTP_PORT || 3010)
    app.listen(PORT, () => {
      console.log(`Aptos MCP Streamable HTTP Server is running on http://localhost:${PORT}`)
      console.log(`MCP Inspector URL: http://localhost:${PORT}`)
    })

    return server
  } catch (error) {
    console.error("Error starting Aptos MCP Streamable HTTP Server:", error)
  }
} 