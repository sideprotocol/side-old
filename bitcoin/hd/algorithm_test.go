package hd_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	bip39 "github.com/tyler-smith/go-bip39"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	ks "github.com/sideprotocol/side/bitcoin/keyring"
	"github.com/sideprotocol/side/bitcoin/keys/segwit"
)

var TestCodec amino.Codec

func init() {

	interfaceRegistry := types.NewInterfaceRegistry()
	TestCodec = amino.NewProtoCodec(interfaceRegistry)
}

func getCodec() codec.Codec {
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	return codec.NewProtoCodec(registry)
}

const (
	mnemonic        = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	xprv            = "xprv9s21ZrQH143K3GJpoapnV8SFfukcVBSfeCficPSGfubmSFDxo1kuHnLisriDvSnRRuL2Qrg5ggqHKNVpxR86QEC8w35uxmGoggxtQTPvfUu"
	path            = "m/86'/0'/0'/0/0"
	internalPubkey  = "cc8a4bc64d897bddc5fbc2f670f7a8ba0b386779106cf1223c6fc5d7cd6fc115"
	expectedAddress = "bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr"
)

func TestKeyring(t *testing.T) {

	// bip39.NewSeed(mnemonic, nil)
	seed := bip39.NewSeed(mnemonic, "")

	println(seed)
	//masterKey, _ := bip32.NewMasterKey(seed)

	//hd.DerivePrivateKeyForPath([32]byte(masterKey.Key), [32]byte(masterKey.ChainCode), path)
	masterKey, chParams := hd.ComputeMastersFromSeed(seed)
	derivedPrivKey, err := hd.DerivePrivateKeyForPath(masterKey, chParams, "m/84'/0'/0'/0/0")
	require.NoError(t, err)
	keyring.NewInMemory(TestCodec)

	algo, _ := keyring.NewSigningAlgoFromString(segwit.KeyType, ks.SupportedAlgorithmsLedger)
	privKey, err := algo.Derive()(mnemonic, "", path)
	require.NoError(t, err)
	require.Equal(t, derivedPrivKey, privKey)

}
