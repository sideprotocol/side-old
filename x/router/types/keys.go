package types

const (
	// ModuleName defines the module name
	ModuleName = "router"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_router"

	// Version defines the current version the IBC module supports
	Version = "router-1"

	// PortID is the default port id that module binds to
	PortID = "router"
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("router-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

type ProcessedKey struct{}
