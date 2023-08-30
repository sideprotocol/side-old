package keeper

import (
	"github.com/sideprotocol/side/x/incentive/types"
)

var _ types.QueryServer = Keeper{}
