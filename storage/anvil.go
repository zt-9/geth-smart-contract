// This file contains util functions related to foundry anvil node or on-chain interaction

package storage

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// AnvilContainer defines the testcontainer used for foundry anvil node.
type AnvilContainer struct {
	testcontainers.Container
	URI string
}

// StartAnvilContainer starts an anvil container with the foundry docker image.
//
// The anvil node will start at default http://localhost:8545.
//
// You may need to manully pull the image with `docker pull ghcr.io/foundry-rs/foundry --platform linux/x86_64` first
// if you are on a Apple Silicon Chip and encountered the error
// `no matching manifest for linux/arm64/v8 in the manifest list entries`.
func StartAnvilContainer(ctx context.Context, chainId *big.Int) (*AnvilContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: "ghcr.io/foundry-rs/foundry",
		// this command will start anvil node on a dynamic port
		// the command has to a one string otherwise the anvil node won't start properly
		Cmd:          []string{fmt.Sprintf("anvil --chain-id %s --host %s", chainId, "0.0.0.0")},
		ExposedPorts: []string{"8545/tcp"},
		WaitingFor:   wait.ForListeningPort("8545/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "8545/tcp")
	if err != nil {
		return nil, err
	}

	// eth client should connect to this uri
	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())

	return &AnvilContainer{Container: container, URI: uri}, nil

}

// StopContainer stops the given anvil testcontainer.
func (a *AnvilContainer) StopContainer(ctx context.Context) error {
	if a == nil {
		return nil
	}
	if err := a.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate container: %w", err)
	}
	return nil
}

// deployContract deploys the Storage.sol contract
func deployContract(ethClient *ethclient.Client, transactor *bind.TransactOpts) (common.Address, *Storage, error) {

	contractAddress, deployTX, contractInstance, err := DeployStorage(transactor, ethClient)

	if err != nil {
		log.Printf("contract deployment failed: %v", err)
		return common.Address{}, nil, err
	}

	log.Printf("execution contract deployed at %v", contractAddress)
	log.Printf("deployment transaction hash %v", deployTX)

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(context.Background(), ethClient, deployTX)
	if err != nil {
		log.Printf("Failed to wait for transaction mining: %v", err)
		return common.Address{}, nil, err
	}

	return contractAddress, contractInstance, nil
}
