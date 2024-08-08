// in this package we use test to show how to deploy contract and read and write contract
package storage

import (
	"context"
	"log"
	"math/big"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stretchr/testify/assert"
)

var (
	client           *ethclient.Client
	contractAddress  common.Address
	auth             *bind.TransactOpts
	contractInstance *Storage
	chainID          *big.Int
	stopAnvil        func()
)

func setup() {
	chainID = big.NewInt(1234)
	log.Printf("set chain id %v", chainID)
	stopAnvil = startAnvil()

	err := deployContract()
	if err != nil {
		log.Fatalf("failed to deploy contract%v", err)
	}
}

func teardown() {
	stopAnvil()
}

// TestContractInteraction tests contract write and read
func TestContractInteraction(t *testing.T) {
	setup()
	t.Cleanup(teardown)

	// test eth client
	ethChainid, err := client.ChainID(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, chainID, ethChainid)

	// write to contract. call Storage.store() function
	user := auth.From
	number := big.NewInt(1234)
	_, err = contractInstance.Store(auth, user, number)
	assert.NoError(t, err, "failed to write to contract")

	// read from contract. call Storage.retreive() function
	result, err := contractInstance.Retrieve(nil, user)
	assert.NoError(t, err, "failed to read from contract")

	assert.Equal(t, result, number)

}

// for start and stop anvil node
func startAnvil() func() {
	cmd := exec.Command("anvil", "--chain-id", chainID.String())

	// Redirect stdout and stderr to the current process
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the Anvil process
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start Anvil: %v", err)
	}

	// Give the Anvil node some time to start
	time.Sleep(2 * time.Second)

	// Return a function to stop the Anvil process
	return func() {
		if err := cmd.Process.Kill(); err != nil {
			log.Fatalf("Failed to stop Anvil: %v", err)
		}
	}
}

func deployContract() error {
	// Connect to the Anvil node
	var err error
	client, err = ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Printf("Failed to connect to the Ethereum client: %v", err)
		return err
	}

	// Load the private key
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Printf("Failed to load private key: %v", err)
		return err
	}

	// Create an authorized transactor
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Printf("Create an authorized transactor failed")
		return err
	}

	// Deploy contract
	var deployTX *types.Transaction
	contractAddress, deployTX, contractInstance, err = DeployStorage(auth, client)

	if err != nil {
		log.Printf("contract deployment failed:%s", err)
		return err
	}

	log.Printf("execution contract deployed at %v", contractAddress)
	log.Printf("deployment transaction hash %v", deployTX)

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(context.Background(), client, deployTX)
	if err != nil {
		log.Printf("Failed to wait for transaction mining: %v", err)
		return err
	}

	return nil
}
