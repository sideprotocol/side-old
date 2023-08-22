package keeper

import (
	"sidechain/x/dex/types"
)

var _ types.QueryServer = Keeper{}
