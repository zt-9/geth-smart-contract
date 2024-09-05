// Tests for anvil.go

package storage

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func TestAnvil_StartAnvilContainer(t *testing.T) {

	ctx := context.Background()
	chainId := big.NewInt(1111)
	anvilC, err := StartAnvilContainer(ctx, chainId)
	assert.NoError(t, err, "failed to start anvil container")

	// check connection with the anvil node
	client, err := ethclient.Dial(anvilC.URI)
	assert.NoError(t, err, "failed to connect to the Ethereum client")

	chainIdRes, err := client.ChainID(context.Background())
	assert.NoError(t, err, "failed to get chainId")
	assert.Equal(t, chainId, chainIdRes)

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		// avoid nil pointer when anvilC is not correctly set
		if anvilC == nil {
			return
		}
		if err := anvilC.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
		client.Client().Close()
	})

}
