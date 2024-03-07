package hd

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	amino "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
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
	mnemonic        = "picnic rent average infant boat squirrel federal assault mercy purity very motor fossil wheel verify upset box fresh horse vivid copy predict square regret"
	xprv            = "xprv9s21ZrQH143K3GJpoapnV8SFfukcVBSfeCficPSGfubmSFDxo1kuHnLisriDvSnRRuL2Qrg5ggqHKNVpxR86QEC8w35uxmGoggxtQTPvfUu"
	path            = "m/86'/0'/0'/0/0"
	internalPubkey  = "cc8a4bc64d897bddc5fbc2f670f7a8ba0b386779106cf1223c6fc5d7cd6fc115"
	expectedAddress = "bc1p5cyxnuxmeuwuvkwfem96lqzszd02n6xdcjrs20cac6yqjjwudpxqkedrcr"
)

func TestKeyring(t *testing.T) {

	// bip39.NewSeed(mnemonic, nil)
	seed, _ := bip39.MnemonicToByteArray(mnemonic)

	println(seed)
	masterKey, _ := bip32.NewMasterKey(seed)
	hd.DerivePrivateKeyForPath([32]byte(masterKey.Key), 0, path)

	keyring.NewInMemory(TestCodec)

	pk := masterKey.PublicKey()

	fmt.Println("PK", pk.String())
	// hd.NewMasterKeyFromMnemonic(mnemonic, "", "bitcoin")

	// algo, _ := keyring.NewSigningAlgoFromString("secp256k1", ks.options.SupportedAlgosLedger)
	// algo.Derive()(&masterKey, "", path)

}
