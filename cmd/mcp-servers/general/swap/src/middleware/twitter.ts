import Logger from "../utils/logger.js"
import { MongoClient } from "mongodb"
import { Ed25519PrivateKey } from "@aptos-labs/ts-sdk"

type ToolHandler<P, R> = (params: P) => Promise<R>

export function withTwitterAuth<P extends Record<string, any> & { twitter_id?: string }, R>(
  handler: (params: P & { privateKey: string; address: string }) => Promise<R>,
  toolName?: string
): ToolHandler<P, R> {
  return async (params: P) => {
    const { twitter_id } = params
    if (!twitter_id) {
      throw new Error("twitter_id is required for all MCP tool calls.")
    }

    const user = await getUserByTwitterId(twitter_id)
    if (!user) {
      throw new Error(`No user found for twitter_id=${twitter_id}`)
    }

    Logger.info(`Tool call${toolName ? ` [${toolName}]` : ''} by ${twitter_id}`)

    // Format private key to be AIP-80 compliant (add 0x prefix if not present)
    const formattedPrivateKey = user.private_key.startsWith('0x') ? user.private_key : `0x${user.private_key}`
    
    // Format address to be 64 characters long (pad with 0s if needed)
    const formattedAddress = user.public_key.startsWith('0x') 
      ? user.public_key.slice(2).padStart(64, '0')
      : user.public_key.padStart(64, '0')
    const fullAddress = `0x${formattedAddress}`

    return handler({ ...params, privateKey: formattedPrivateKey, address: fullAddress } as P & {
      privateKey: string;
      address: string
    })
  }
}


async function getUserByTwitterId(twitterId: string) : Promise<{ twitter_id: string, public_key: string, private_key: string } | null> {
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
      Logger.error("Error fetching user by twitter_id:", error)
      return null
    }
}

async function getMongoClient() : Promise<MongoClient> {
    const uri = process.env.MONGO_URI || "mongodb://localhost:27017"
    const client = new MongoClient(uri)
    await client.connect()
    return client
}
  