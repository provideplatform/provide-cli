# provide-CLI

The command-line interface for building ÐApps that use [Provide](http://provide.services/) APIs

## Quickstart

1. Create a [Provide account](https://dawn.provide.services/sign-in)<br>

2. Download the binary: [provide-cli.exe](https://github.com/provideservices/provide-cli/tree/dev/binary)<br>

3. Make the file executable in your terminal with `chmod +x provide-cli`<br>

4. Authenticate using `./provide-cli authenticate`<br>

5. Enter the Email/Password used for your [Provide account](https://dawn.provide.services/sign-in)<br>

## Basic Usage

Here's a taste of what your shiny new `./provide-cli` executable does:<br>

| Command | Action |
| :--- | :--- |
| `./provide-cli networks list --public` | List available networks |
| `./provide-cli dapps init --name '<myAwesomedApp>' --network <networkId>` | Create your dApp, API token and wallet |
| `./provide-cli deploy <MyFlawlessContract.sol> --application <applicationId> --network <networkId> --wallet <walletId>` | Deploy a contract | <br>

[Show me the full documentation](https://provideservices.github.io/docs/)

## Speak Up! <br>

If you see a problem, help us help you by [creating an issue](https://github.com/provideservices/provide-cli/issues)
