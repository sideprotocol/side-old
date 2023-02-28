package keeper

import (
	"sidechain/x/side/types"
)

var _ types.QueryServer = Keeper{}
