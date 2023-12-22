package types

const (
	// ModuleName defines the module name
	ModuleName = "yield"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_yield"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	// Host chain keys prefix the HostChain structs
	HostChainKey          = "HostChain-value-"
	DepositRecordKey      = "DepositRecord-value-"
	DepositRecordCountKey = "DepositRecord-count-"
)
