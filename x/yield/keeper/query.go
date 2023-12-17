package keeper

import (
	"github.com/sideprotocol/side/x/yield/types"
)

var _ types.QueryServer = Keeper{}
