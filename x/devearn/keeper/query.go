package keeper

import (
	"sidechain/x/devearn/types"
)

var _ types.QueryServer = Keeper{}
