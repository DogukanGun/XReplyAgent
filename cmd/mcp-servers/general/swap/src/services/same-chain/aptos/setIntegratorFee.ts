import {
  Account,
  AccountAddress,
  Aptos,
  AptosConfig,
  Ed25519PrivateKey,
  Network,
} from "@aptos-labs/ts-sdk"
import "dotenv/config"

// Kana Router contract address
export const KANA_ROUTER =
  "0x9538c839fe490ccfaf32ad9f7491b5e84e610ff6edc110ff883f06ebde82463d"

/**
 * Set integrator referral fee for Kana Labs router
 *
 * This function allows integrators to set their referral fee percentage
 * for swaps routed through Kana Labs. The fee is deducted from the output
 * amount and sent to the integrator's wallet.
 *
 * Fee calculation:
 * - 1 bps = 0.01%
 * - 10 bps = 0.1%
 * - 100 bps = 1%
 *
 * @param params - Object containing privateKey, address, and feeBps
 * @returns Transaction result object
 */
export const setIntegratorFee = async (params: {
  privateKey: string
  address: string,
  feeBps: number // Fee in basis points (10 bps = 0.1%)
}) => {

  // Setup signer from private key
  const aptosSigner = Account.fromPrivateKey({
    privateKey: new Ed25519PrivateKey(params.privateKey),
    address: AccountAddress.from(params.address),
    legacy: true,
  })

  // Setup Aptos provider
  const aptosConfig = new AptosConfig({ network: Network.MAINNET })
  const aptos = new Aptos(aptosConfig)

  // Build transaction
  const payload = await aptos.transaction.build.simple({
    sender: aptosSigner.accountAddress.toString(),
    data: {
      function: `${KANA_ROUTER}::KanalabsRouterV2::set_referral_swap_fee`,
      functionArguments: [params.feeBps],
      typeArguments: [],
    },
  })

  // Sign and submit
  const sign = await aptos.signAndSubmitTransaction({
    signer: aptosSigner,
    transaction: payload,
  })

  // Wait for transaction
  const submit = await aptos.waitForTransaction({ transactionHash: sign.hash })
  console.log(submit)
}
