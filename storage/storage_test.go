// in this package we use test to show how to deploy contract and read and write contract
package storage

import (
	"context"
	"log"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stretchr/testify/assert"
)

var (
	chainId          *big.Int
	client           *ethclient.Client
	contractAddress  common.Address
	auth             *bind.TransactOpts
	contractInstance *Storage
	anvilC           *AnvilContainer
)

func setup() {
	var err error
	chainId = big.NewInt(1234)
	// start anvil container
	anvilC, err = StartAnvilContainer(context.Background(), chainId)
	if err != nil {
		log.Fatalf("failed to create anvil testcontainer: %s", err)
	}

	// eth client
	client, err = ethclient.Dial(anvilC.URI)

	if err != nil {
		log.Fatalf("failed to create eth client: %v", err)
	}

	// Load the private key
	// this the default account private key in anvil
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}

	// Create an authorized transactor
	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatalf("failed to create transactor: %v", err)
	}

	// deploy contract
	contractAddress, contractInstance, err = deployContract(client, auth)
	if err != nil {
		log.Fatalf("failed to deploy contract%v", err)
	}
}

func teardown() {
	// stop anvil testcontainer
	if err := anvilC.StopContainer(context.Background()); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	// close eth client
	client.Client().Close()

}

// TestContractInteraction tests contract write and read
func TestContractInteraction(t *testing.T) {
	setup()
	t.Cleanup(teardown)
	var err error

	// Wait for the context to timeout
	time.Sleep(30 * time.Millisecond)

	// write to contract. call Storage.store() function
	user := auth.From
	number := big.NewInt(1234)
	tx, err := contractInstance.Store(auth, user, number)
	assert.NoError(t, err, "failed to write to contract")
	// Wait for the transaction to be mined
	_, err = bind.WaitMined(context.Background(), client, tx)
	assert.NoError(t, err, "failed to mint tx")

	// read from contract. call Storage.retreive() function
	result, err := contractInstance.Retrieve(nil, user)
	assert.NoError(t, err, "failed to read from contract")

	assert.Equal(t, result, number)

}
