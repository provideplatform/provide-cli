# provide-cli

A command-line interface for managing your decentralized applications and infrastructure using [Provide](https://provide.services).

## Developer Quickstart

1. Create a [Provide account](https://dawn.provide.services/sign-up)

2. Get a precompiled binary for your platform (or `go get github.com/provideservices/provide-cli` and run from source): [win](https://github.com/provideservices/provide-cli/tree/dev/bin/windows) | [mac](https://github.com/provideservices/provide-cli/tree/dev/bin/osx) | [linux](https://github.com/provideservices/provide-cli/tree/dev/bin/linux)

3. Make sure [solc](https://solidity.readthedocs.io/en/latest/installing-solidity.html) (Solidity compiler) is installed. You can check this with `solc --version`

4. Authenticate: `provide-cli authenticate` - this will authorize and cache a Provide API token in your home directory for continued use of the CLI - after you have authenticated you can register an application using the CLI and start building.

## Quickstart

You can deploy a contract and be ready to use one of our [API clients](https://github.com/provideservices) to build dApps using [a variety of underlying protocols](https://github.com/providenetwork/node) in just 4 steps: 

#### 1: View the networks
```provide-cli networks list --public```

#### 2: Create your dApp, API token and wallet
```provide-cli dapps init --name '<myApp>' --network <networkId>```

#### 3: Deploy your compiled contract(s)
```provide-cli contracts deploy <MyContract.sol> --application <applicationId> --network <networkId> --wallet <walletId>```  

Looking for the [API docs](https://docs.provide.services)?
