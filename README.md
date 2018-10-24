# provide-CLI

Command-line interface for building decentralized and semi-centralized apps using [provide](http://provide.services/) APIs

## Quickstart
***Note**: This CLI builds with [go](https://golang.org/doc/install#install)*

1. Create a [provide account](https://dawn.provide.services/sign-in)<br>

2. Install the CLI by cloning this repository<br>

3. Change into the directory and `go build`<br>

4. Authenticate using `./provide-cli-dev authenticate`<br>

5. Enter the Email/Password used for your [provide account](https://dawn.provide.services/sign-in)<br>

## Basic Usage

Provide makes a `./provide-cli-dev` executable available to you in your terminal:<br>

| Command | Action |
| :--- | :--- |
| `./provide-cli-dev networks list --public` | List available networks |
| `./provide-cli-dev dapps init --name '<myAwesomedApp>' --network <networkId>` | Create your dApp, API token and wallet |
| `./provide-cli-dev deploy <MyFlawlessContract.sol> --application <applicationId> --network <networkId>   --wallet <walletAddress` | Deploy a contract to the testnet |
| `./provide-cli-dev deploy <MyFlattenedContracts.sol> --application <applicationId> --network <networkId>  --wallet <walletAddress>` | Deploy ***multiple*** contracts and dependencies to the testnet |<br>

## Speak Up! <br>

If you see a problem, make us aware and help others by [creating an issue](https://github.com/provideservices/provide-cli/issues)
