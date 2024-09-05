# Geth Smart Contract
This guide provides tutorials for using Geth to deploy and interact with smart contracts.

## Prerequisite
- [Golang](https://go.dev/doc/install)
- [Docker](https://www.docker.com/)

## 1. Generate the ABI and BIN
We can use `solc` or tools like Hardhat or Foundry to generate the ABI and BIN files for our contracts.

### Using `solc`
1. Download `solc` using [solc-select](https://github.com/crytic/solc-select.git):
   ```bash
   git clone https://github.com/crytic/solc-select.git
   cd solc-select
   ./solc-select install
   ```
2. Install the desired Solidity version:
   ```bash
   solc-select use 0.8.26 --always-install
   ```
3. Navigate to the `contract` folder and run the command to generate the ABI and BIN files:
   ```bash
   solc --abi --bin -o build Storage.sol
   ```

## 2. Generate Go Files for Contract Deployment and Interaction
1. Install the `go-ethereum` package, which is required for the `abigen` tool:
   ```bash
   go get github.com/ethereum/go-ethereum
   ```
2. Install the `abigen` tool:
   ```bash
    go install github.com/ethereum/go-ethereum/cmd/abigen@latest
   ```
3. Create a folder named `storage` for our Go package:
   ```bash
   mkdir storage
   cd storage
   ```
4. Run the `abigen` command in the `storage` folder to generate the Go files for the `storage` package:
   ```bash
   abigen --abi=../contract/build/Storage.abi --bin=../contract/build/Storage.bin --pkg=storage --out=storage.go
   ```

## 3. Contract Deployment and Interaction
The generated `storage.go` file provides the `DeployStorage()` function to deploy the contract, as well as functions to interact with it.

To call the `store()` function in the `Storage` contract, use the generated `Store()` function in the `storage.go` file.

Refer to the example in `storage/storage_test.go` for contract deployment and interaction. 
In this test, we use Anvil to start a local ethereum testnet node. Ensure you have [Foundry](https://book.getfoundry.sh/getting-started/installation) installed before running the test:
```bash
make test
```

## Notes
- [Testcontainers](https://golang.testcontainers.org/) is used for the tests. so we need to start docker before running the tests.
-  The Testcontainers fetches the the foundry docker image from `ghcr.io/foundry-rs/foundry`. The fetch may fail with error `no matching manifest for linux/arm64/v8 in the manifest list entries` on Apple Silicon chips. You may need to manully pull the image with `docker pull ghcr.io/foundry-rs/foundry --platform linux/x86_64` first if you are on a Apple Silicon Chip and encountered the error.
- Run `make clean` if you encounter any Go module or package load errors.
