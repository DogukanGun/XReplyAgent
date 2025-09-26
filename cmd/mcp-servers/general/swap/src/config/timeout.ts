// Timeout configuration for blockchain operations
export const TIMEOUTS = {
  // Default MCP request timeout (5 minutes)
  MCP_REQUEST: parseInt(process.env.MCP_REQUEST_TIMEOUT || "300000"),
  
  // Blockchain transaction timeout (3 minutes)
  BLOCKCHAIN_TX: parseInt(process.env.BLOCKCHAIN_TX_TIMEOUT || "180000"),
  
  // Cross-chain swap timeout (5 minutes)
  CROSS_CHAIN_SWAP: parseInt(process.env.CROSS_CHAIN_SWAP_TIMEOUT || "300000"),
  
  // Same-chain swap timeout (2 minutes)
  SAME_CHAIN_SWAP: parseInt(process.env.SAME_CHAIN_SWAP_TIMEOUT || "120000")
}

// Helper function to create timeout promise
export const withTimeout = <T>(promise: Promise<T>, timeoutMs: number, operation: string): Promise<T> => {
  return Promise.race([
    promise,
    new Promise<never>((_, reject) =>
      setTimeout(() => reject(new Error(`${operation} timed out after ${timeoutMs}ms`)), timeoutMs)
    )
  ])
} 