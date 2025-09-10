#!/usr/bin/env node

import { startSSEServer } from "./server/sse.ts"

async function main() {
  let server = await startSSEServer()

  if (!server) {
    console.error("Failed to start server")
    process.exit(1)
  }

  const handleShutdown = async () => {
    await server.close()
    process.exit(0)
  }
  // Handle process termination
  process.on("SIGINT", handleShutdown)
  process.on("SIGTERM", handleShutdown)
}

main()
