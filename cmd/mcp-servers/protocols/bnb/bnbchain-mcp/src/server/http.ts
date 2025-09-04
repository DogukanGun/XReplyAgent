import "dotenv/config"

import cors from "cors"
import express from "express"
import type { Request, Response } from "express"

import Logger from "../utils/logger.ts"
import { startServer } from "./base.ts"

// Simple HTTP transport implementation
class HttpTransport {
  public sessionId: string
  private req: Request
  private res: Response

  constructor(req: Request, res: Response) {
    this.req = req
    this.res = res
    this.sessionId = `http-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }

  async handlePostMessage(req: Request, res: Response, body: any) {
    try {
      // For now, just echo back a success response
      // In a full implementation, this would route to the actual MCP server handlers
      res.json({
        jsonrpc: "2.0",
        id: body.id,
        result: { success: true, method: body.method }
      })
    } catch (error) {
      res.status(500).json({
        jsonrpc: "2.0",
        id: body?.id,
        error: {
          code: -32603,
          message: "Internal error",
          data: String(error)
        }
      })
    }
  }
}

export const startHttpServer = async () => {
  try {
    const app = express()
    const server = startServer()
    app.use(cors())
    app.use(express.json())

    // Log the current log level on startup
    Logger.info(`Starting HTTP server with log level: ${Logger.getLevel()}`)

    // to support multiple simultaneous connections we have a lookup object from
    // sessionId to transport
    const transports: { [sessionId: string]: HttpTransport } = {}

    app.get("/", (_, res) => {
      res.send("BNBChain MCP HTTP Server is running")
    })

    app.get("/mcp", (_, res) => {
      res.send("BNBChain MCP HTTP Server - MCP endpoint")
    })

    app.post("/mcp", async (req: Request, res: Response) => {
      const sessionId = req.query.sessionId as string || `http-${Date.now()}`
      let transport = transports[sessionId]

      if (!transport) {
        // Create a new transport for this session
        transport = new HttpTransport(req, res)
        transports[transport.sessionId] = transport
        Logger.info("New HTTP connection established", {
          sessionId: transport.sessionId
        })

        try {
          await server.connect(transport as any)
        } catch (error) {
          Logger.error("Error connecting transport", {
            sessionId: transport.sessionId,
            error
          })
        }
      }

      try {
        Logger.debug("Handling message", { sessionId, body: req.body })
        await transport.handlePostMessage(req, res, req.body)
      } catch (error) {
        Logger.error("Error handling message", { sessionId, error })
        res.status(500).send("Internal server error")
      }
    })

    const PORT = process.env.HTTP_PORT || 3001
    app.listen(PORT, () => {
      Logger.info(
        `BNBChain MCP HTTP Server is running on http://localhost:${PORT}`
      )
    })

    return server
  } catch (error) {
    Logger.error("Error starting BNBChain MCP HTTP Server:", error)
  }
} 