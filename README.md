# provide-CLI

The command-line interface for building decentralized and semi-centralized apps using [provide](http://provide.services/) APIs

## Quickstart
***Note**: This CLI builds with [go](https://golang.org/doc/install#install)*

1. Create a [provide account](https://dawn.provide.services/sign-in)<br>

2. Install the CLI by cloning this repository<br>

3. Change into the directory and `go build`<br>

4. Authenticate using `./provide-cli authenticate`<br>

5. Enter the Email/Password used for your [provide account](https://dawn.provide.services/sign-in)<br>

## Basic Usage

Provide makes a `./provide-cli` executable available in your terminal:<br>

| Command | Action |
| :--- | :--- |
| `./provide-cli networks list --public` | List available networks |
| `./provide-cli dapps init --name '<myAwesomedApp>' --network <networkId>` | Create your dApp, API token and wallet |
| `./provide-cli deploy <MyFlawlessContract.sol> --application <applicationId> --network <networkId>   --wallet <walletAddress>` | Deploy a contract to the testnet | <br>

[Show me more documentation](https://provideservices.github.io/docs/)

## Speak Up! <br>

If you see a problem, help us help you by [creating an issue](https://github.com/provideservices/provide-cli/issues)
