package types

import (
	"math/big"

	"github.com/btcsuite/btcd/chaincfg"
)

const (
	// ModuleName defines the module name
	ModuleName = "btclightclient"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	ParamsStoreKey = []byte{0x1}
	SequenceKey    = []byte{0x2}

	// Host chain keys prefix the HostChain structs
	BtcBlockHeaderHashPrefix   = []byte{0x11} // prefix for each key to a block header, for a hash
	BtcBlockHeaderHeightPrefix = []byte{0x12} // prefix for each key to a block hash, for a height
	BtcBestBlockHeaderKey      = []byte{0x13} // key for the best block height
	BtcSigningRequestPrefix    = []byte{0x14} // prefix for each key to a signing request

	BtcUtxoKeyPrefix      = []byte{0x15} // prefix for each key to a utxo
	BtcOwnerUtxoKeyPrefix = []byte{0x16} // prefix for each key to an owned utxo

	ChainCfg = &chaincfg.MainNetParams
)

func Int64ToBytes(number uint64) []byte {
	big := new(big.Int)
	big.SetUint64(number)
	return big.Bytes()
}

func BtcUtxoKey(hash string, vout uint64) []byte {
	return append(append(BtcUtxoKeyPrefix, []byte(hash)...), Int64ToBytes(vout)...)
}

func BtcOwnerUtxoKey(owner string, hash string, vout uint64) []byte {
	key := append(BtcOwnerUtxoKeyPrefix, []byte(owner)...)
	key = append(key, []byte(hash)...)

	return append(key, Int64ToBytes(vout)...)
}

func BtcBlockHeaderHashKey(hash string) []byte {
	return append(BtcBlockHeaderHashPrefix, []byte(hash)...)
}

func BtcBlockHeaderHeightKey(height uint64) []byte {
	return append(BtcBlockHeaderHeightPrefix, Int64ToBytes(height)...)
}

func BtcSigningRequestKey(sequence uint64) []byte {
	return append(BtcSigningRequestPrefix, Int64ToBytes(sequence)...)
}