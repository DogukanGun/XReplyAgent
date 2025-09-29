import Logger from "../utils/logger.js"
import { MongoClient } from "mongodb"

type ToolHandler<P, R> = (params: P) => Promise<R>

export function withTwitterAuth<P extends Record<string, any> & { sender_twitter_id?: string, receiver_twitter_id?: string }, R>(
  handler: (params: P & { sender_private_key: string, receiver_private_key: string }) => Promise<R>,
  toolName?: string
): ToolHandler<P, R> {
  return async (params: P) => {
    const { sender_twitter_id, receiver_twitter_id } = params
    if (!sender_twitter_id) {
      throw new Error("twitter_id is required for all MCP tool calls.")
    }
    if (!receiver_twitter_id) {
      throw new Error("receiver_twitter_id is required for all MCP tool calls.")
    }

    const user = await getUserByTwitterId(sender_twitter_id)
    const receiverUser = await getUserByTwitterId(receiver_twitter_id)
    if (!user) {
      throw new Error(`No user found for twitter_id=${sender_twitter_id}, please create a wallet`)
    }
    if (!receiverUser) {
      throw new Error(`No user found for twitter_id=${receiver_twitter_id}, please create a wallet`)
    }
    Logger.info(`Tool call${toolName ? ` [${toolName}]` : ''} by ${sender_twitter_id}`)

    return handler({ ...params, sender_private_key: `0x${user.private_key}`, receiver_private_key: `0x${receiverUser.private_key}`  } as P & {
      sender_private_key: string,
      receiver_private_key: string
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
  