package hd

import (
	bip39 "github.com/tyler-smith/go-bip39"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/sideprotocol/side/bitcoin/keys/segwit"
)

const (
	// EthSecp256k1Type defines the ECDSA secp256k1 used on Ethereum
	BTCSecp256k1Type = hd.PubKeyType(segwit.KeyType)
)

var (
	// SupportedAlgorithms defines the list of signing algorithms used on Ethermint:
	//  - eth_secp256k1 (Ethereum)
	//  - secp256k1 (Tendermint)
	SupportedAlgorithms = keyring.SigningAlgoList{BtcSecp256k1, hd.Secp256k1}
	// SupportedAlgorithmsLedger defines the list of signing algorithms used on Ethermint for the Ledger device:
	//  - eth_secp256k1 (Ethereum)
	//  - secp256k1 (Tendermint)
	SupportedAlgorithmsLedger = keyring.SigningAlgoList{BtcSecp256k1, hd.Secp256k1}
)

// BthSecp256k1Option defines a function keys options for the ethereum Secp256k1 curve.
// It supports eth_secp256k1 and secp256k1 keys for accounts.
func BthSecp256k1Option() keyring.Option {
	return func(options *keyring.Options) {
		options.SupportedAlgos = SupportedAlgorithms
		options.SupportedAlgosLedger = SupportedAlgorithmsLedger
	}
}

var (
	_ keyring.SignatureAlgo = BtcSecp256k1

	// BtcSecp256k1 uses the Bitcoin secp256k1 ECDSA parameters.
	BtcSecp256k1 = btcSecp256k1Algo{}
)

type btcSecp256k1Algo struct{}

// Name returns eth_secp256k1
func (s btcSecp256k1Algo) Name() hd.PubKeyType {
	return BTCSecp256k1Type
}

// Derive derives and returns the eth_secp256k1 private key for the given mnemonic and HD path.
func (s btcSecp256k1Algo) Derive() hd.DeriveFn {
	return func(mnemonic, bip39Passphrase, path string) ([]byte, error) {

		seed := bip39.NewSeed(mnemonic, bip39Passphrase)
		masterKey, chParams := hd.ComputeMastersFromSeed(seed)
		derivedPrivKey, err := hd.DerivePrivateKeyForPath(masterKey, chParams, "m/84'/0'/0'/0/0")
		//masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams) // Adjust chaincfg as needed for testnet or others
		if err != nil {
			return nil, err
		}

		return derivedPrivKey, nil // This returns the serialized private key
	}
}

// Generate generates a btc_secp256k1 private key from the given bytes.
func (s btcSecp256k1Algo) Generate() hd.GenerateFn {
	return func(bz []byte) cryptotypes.PrivKey {
		bzArr := make([]byte, segwit.PrivKeySize)
		copy(bzArr, bz)

		// TODO: modulo P
		return &segwit.PrivKey{
			Key: bzArr,
		}
	}
}
