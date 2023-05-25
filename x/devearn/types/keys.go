package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "devearn"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_devearn"
)

// ModuleAddress is the native module address for incentives module
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

const (
	prefixDevEarn = iota + 1
)

// KVStore key prefixes
var (
	KeyPrefixDevEarn = []byte{prefixDevEarn}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	AssetsKey = "Assets/value/"
)
