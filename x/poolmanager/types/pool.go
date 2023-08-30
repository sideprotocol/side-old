package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/proto"
	sideutils "github.com/sideprotocol/side/sideutils"
)

// PoolI defines an interface for pools that hold tokens.
type PoolI interface {
	proto.Message

	GetAddress() sdk.AccAddress
	String() string
	GetId() uint64
	// Returns whether the pool has swaps enabled at the moment
	IsActive(ctx sdk.Context) bool
	GetType() PoolType
	// AsSerializablePool returns the pool in a serializable form (useful when a model wraps the proto)
	AsSerializablePool() PoolI
}

// NewPoolAddress returns an address for a pool from a given id.
func NewPoolAddress(poolId uint64) sdk.AccAddress {
	return sideutils.NewModuleAddressWithPrefix(ModuleName, "pool", sdk.Uint64ToBigEndian(poolId))
}
