package emulator_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/flow-go/pkg/types"
	"github.com/dapperlabs/flow-go/sdk/emulator"
	"github.com/dapperlabs/flow-go/sdk/keys"
)

func TestCallScript(t *testing.T) {
	b := emulator.NewEmulatedBlockchain(emulator.DefaultOptions)

	tx := types.Transaction{
		Script:             []byte(addTwoScript),
		ReferenceBlockHash: nil,
		Nonce:              getNonce(),
		ComputeLimit:       10,
		PayerAccount:       b.RootAccountAddress(),
		ScriptAccounts:     []types.Address{b.RootAccountAddress()},
	}

	sig, err := keys.SignTransaction(tx, b.RootKey())
	assert.Nil(t, err)

	tx.AddSignature(b.RootAccountAddress(), sig)

	// Sample call (value is 0)
	value, err := b.CallScript([]byte(sampleCall))
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(0), value)

	// Submit tx1 (script adds 2)
	err = b.SubmitTransaction(tx)
	assert.Nil(t, err)

	// Sample call (value is 2)
	value, err = b.CallScript([]byte(sampleCall))
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(2), value)
}
