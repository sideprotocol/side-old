package segwit_test

import (
	//"fmt"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/sideprotocol/side/bitcoin/keys/segwit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
)

func TestSegwit(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	seed := bip39.NewSeed(mnemonic, "")

	masterKey, chParams := hd.ComputeMastersFromSeed(seed)
	derivedPrivKey, err := hd.DerivePrivateKeyForPath(masterKey, chParams, "m/84'/0'/0'/0/0")
	assert.NoError(t, err, "Private key derivation should not fail")

	privKey := segwit.PrivKey{Key: derivedPrivKey}
	pubKey := privKey.PubKey()
	assert.NotNil(t, pubKey, "Public key should not be nil")

	bech32Address, err := segwit.BitCoinAddr(pubKey.Bytes())
	assert.NoError(t, err)
	t.Logf("Generated SegWit Address: %s", bech32Address)
	// Check if the Bech32 encoded address has the correct prefix and structure.
	assert.True(t, strings.HasPrefix(bech32Address, "bc1q"), "Address should start with 'bc1q'")
	t.Logf("Generated SegWit Address: %s", bech32Address)
	acc, err := sdk.AccAddressFromBech32(bech32Address)
	require.NoError(t, err)
	t.Logf("Generated SegWit Address: %s", acc)

}
