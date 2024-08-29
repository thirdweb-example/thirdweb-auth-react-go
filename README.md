![Group 1](https://github.com/thirdweb-example/thirdweb-auth-express/assets/17715009/06383e68-9c65-4265-8505-e88e573443f9)

# React + Go Auth Starter Template

[<img alt="thirdweb SDK" src="https://img.shields.io/npm/v/thirdweb?label=Thirdweb SDK&style=for-the-badge&logo=npm" height="30">](https://www.npmjs.com/package/thirdweb)
[<img alt="Discord" src="https://img.shields.io/discord/834227967404146718.svg?color=7289da&label=discord&logo=discord&style=for-the-badge" height="30">](https://discord.gg/thirdweb)

The template is a standard monorepo separated into a client and a server. The client is a React app using Vite and the server is an Go API using Gin.

The client uses thirdweb's Auth SDK to communicate with the server via your own routes. When a user is authenticated, the server sets a JWT to keep the user authenticated on future requests. Thirdweb's Auth SDK takes care of the difficult pieces like payload generation and signing while allowing you to control how authentication happens between your backend and frontend.

## Getting Started

> Note, all client-side commands can be substituted with npm, yarn, or pnpm if preferred.

First, install dependencies on both the client and server:

```bash
cd client && bun install
cd server && go get .
```

### Environment variables

Copy the `.env.example` files to `.env`:

```bash
cp client/.env.example client/.env
cp server/.env.example server/.env
```

In each `.env` file, complete any necessary values. There are comments in each `.env.example` file explaining what each value does.

### Running the app

Then, run the client and server in development mode:

```bash
cd client && bun run dev
cd server && go run .
```

## How it works

All authentication logic happens in just two files: one in the client, and one on the server.

On the client, `ConnectButtonAuth.tsx` uses the thirdweb `ConnectButton` to specify the `auth` object, which contains the callbacks to manage all auth requests your application makes. This same process can be done with your own custom components using the underlying functions in `thirdweb/auth`. The connect button is a template for the four functions every application's auth needs to implement:

-   `isLoggedIn` checks whether or not the current user is logged in, the criteria for this is up to you but the default uses a JWT in local cookie storage
-   `getLoginPayload` retrieves the message to sign for login from the server. The function to generate this payload is provided for you in the SDK, you simply need to send it to the server
-   `doLogin` attempts to log the user with a signed payload. Normally, this involves sending the payload to your Express server.
-   `doLogout` logs the user out. This function is entirely up to you and is based on how you handle authentication. In this app for example, it clears the local storage.

> Note: All complex client-side logic such as signing payloads, caching and invalidating state, and more is handled for you by the SDK. Even if you're using your own components, thirdweb provides the underlying functions to execute these steps. Check out our documentation to learn more.

On the server, `main.go` handles all routes for the app. The SDK provides a number of useful functions to keep your backend routes as simple as possible:

-   `generatePayload` generates a login payload for a given user address. This is the message the user's accont signs on the client and returns to the server for verification.
-   `verifyPayload` verifies the user's signed payload to either log them in, or reject the attempt.
-   `generateJWT` generates a JWT for the user to keep them logged in.
-   `verifyJWT` checks if a given JWT is valid when authenticating requests.

## Documentation

-   [TypeScript SDK](https://portal.thirdweb.com/typescript/v5)
-   [Sign in with Ethereum Spec](https://eips.ethereum.org/EIPS/eip-4361)
-   [Gin](https://github.com/gin-gonic/gin)

## Support

For help or feedback, please [visit our support site](https://thirdweb.com/support)
