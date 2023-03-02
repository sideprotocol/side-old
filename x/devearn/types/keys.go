package types

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

func KeyPrefix(p string) []byte {
	return []byte(p)
}
