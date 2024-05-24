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
	// Host chain keys prefix the HostChain structs
	BtcBlockHeaderHashPrefix   = []byte{0x11} // prefix for each key to a block header, for a hash
	BtcBlockHeaderHeightPrefix = []byte{0x12} // prefix for each key to a block hash, for a height
	BtcBestBlockHeaderKey      = []byte{0x13} // key for the best block height

	ChainCfg = &chaincfg.MainNetParams
)

func Int64ToBytes(number uint64) []byte {
	big := new(big.Int)
	big.SetUint64(number)
	return big.Bytes()
}

func BtcBlockHeaderHashKey(hash string) []byte {
	return append(BtcBlockHeaderHashPrefix, []byte(hash)...)
}

func BtcBlockHeaderHeightKey(height uint64) []byte {
	return append(BtcBlockHeaderHashPrefix, Int64ToBytes(height)...)
}
