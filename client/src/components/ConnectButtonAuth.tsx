import { ConnectButton } from "thirdweb/react";
import { client } from "../lib/client";
import { LoginPayload, VerifyLoginPayloadParams } from "thirdweb/auth";
import { get, post } from "../lib/api";
import { sepolia } from "thirdweb/chains";

export default function ConnectButtonAuth() {
  return (
    <ConnectButton
      theme="dark"
      client={client}
      /**
       * Looking to authenticate with account abstraction enabled? Uncomment the following lines:
       *
       * accountAbstraction={{
       *  chain: sepolia,
       *  factoryAddress: "0x5cA3b8E5B82D826aF6E8e9BA9E4E8f95cbC177F4",
       *  gasless: true,
       * }}
       */
      auth={{
        /**
         * 	`getLoginPayload` should @return {VerifyLoginPayloadParams} object.
         * 	This can be generated on the server with the generatePayload method.
         */
        getLoginPayload: async (params: {
          address: string;
        }): Promise<LoginPayload> => {
          const payload = (await get({
            url: process.env.AUTH_API + "/login",
            params: {
              address: params.address,
              chainId: sepolia.id.toString(),
            },
          })).payload;

          return {
            uri: payload.uri,
            domain: payload.domain,
            address: payload.address,
            statement: payload.statement,
            expiration_time: payload.expiration_time,
            issued_at: payload.issued_at,
            nonce: payload.nonce,
            version: payload.version,
            invalid_before: payload.invalid_before,
            chain_id: payload.chain_id,
          }
        },
        /**
         * 	`doLogin` performs any logic necessary to log the user in using the signed payload.
         * 	In this case, this means sending the payload to the server for it to set a JWT cookie for the user.
         */
        doLogin: async (params: VerifyLoginPayloadParams) => {
          await post({
            url: process.env.AUTH_API + "/login",
            params,
          });
        },
        /**
         * 	`isLoggedIn` returns true or false to signal if the user is logged in.
         * 	Here, this is done by calling the server to check if the user has a valid JWT cookie set.
         */
        isLoggedIn: async () => {
          const response = await get({
            url: process.env.AUTH_API + "/isLoggedIn",
          });
          return response.result;
        },
        /**
         * 	`doLogout` performs any logic necessary to log the user out.
         * 	In this case, this means sending a request to the server to clear the JWT cookie.
         */
        doLogout: async () => {
          await post({
            url: process.env.AUTH_API + "/logout",
          });
        },
      }}
    />
  );
}
