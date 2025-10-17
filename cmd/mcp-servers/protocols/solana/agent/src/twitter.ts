import { MongoClient } from "mongodb"

type ToolHandler<P, R> = (params: P) => Promise<R>

export function withTwitterAuth<P extends Record<string, any> & { twitter_id?: string }, R>(
  handler: (params: P & { privateKey: string }) => Promise<R>,
  toolName?: string
): ToolHandler<P, R> {
  return async (params: P) => {
    const { twitter_id } = params
    if (!twitter_id) {
      throw new Error("twitter_id is required for all MCP tool calls.")
    }

    const user = await getUserByTwitterId(twitter_id)
    if (!user) {
      throw new Error(`No user found for twitter_id=${twitter_id}, please create a wallet`)
    }

    console.log(`Tool call${toolName ? ` [${toolName}]` : ''} by ${twitter_id}`)

    return handler({ ...params, privateKey: `${user.solana_private_key}`  } as P & {
      privateKey: string
    })
  }
}


async function getUserByTwitterId(twitterId: string) : Promise<{ twitter_id: string, public_key: string, solana_private_key: string } | null> {
    try {
      const client = await getMongoClient()
  const db = client.db("xreplyagent")
  const collection = db.collection("users")
  
      const user = await collection.findOne<{ 
        twitter_id: string; 
        solana_public_key: string; 
        solana_private_key: string; 
        username?: string 
      }>({ twitter_id: twitterId })
  
      if (!user) {
        return null
      }
  
      return {
        twitter_id: user.twitter_id,
        public_key: user.solana_public_key,
        solana_private_key: user.solana_private_key,
      }
    } catch (error) {
      console.error("Error fetching user by twitter_id:", error)
      return null
    }
}

async function getMongoClient() : Promise<MongoClient> {
    const uri = process.env.MONGO_URI || "mongodb://localhost:27017"
    const client = new MongoClient(uri)
    await client.connect()
    return client
}
  