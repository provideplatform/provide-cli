# provide-CLI

The command-line interface for building ÐApps that use [Provide](http://provide.services/) APIs

## Environment Setup

1. Create a [Provide account](https://dawn.provide.services/sign-in)<br>

2. Get the bin: [win](https://github.com/provideservices/provide-cli/tree/dev/bin/windows) | [mac](https://github.com/provideservices/provide-cli/tree/dev/bin/osx) | [linux](https://github.com/provideservices/provide-cli/tree/dev/bin/linux)<br>

3. Make sure [solC](https://solidity.readthedocs.io/en/latest/installing-solidity.html) (Solidity compiler) is installed. You can check this with `solc --version` <br>

4. Authenticate using your [Provide](https://dawn.provide.services/sign-in) email/password: `provide-cli authenticate`<br>


## Basic Usage

Your project is fully operable in just 4 steps: <br>


#### 1: View the networks
```provide-cli networks list --public```


#### 2: Create your dApp, API token and wallet
```provide-cli dapps init --name '<myAwesomedApp>' --network <networkId>```



#### 3: Deploy your compiled contract(s)
```provide-cli contracts deploy <MyFlawlessContract.sol> --application <applicationId> --network <networkId> --wallet <walletId>```  <br>


Ctrl + click to [see the full documentation](https://provideservices.github.io/docs/) in a new tab.

## Speak Up! <br>

If you see a problem, help us help you by [creating an issue](https://github.com/provideservices/provide-cli/issues).
