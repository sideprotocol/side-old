package keeper

import (
	"github.com/sideprotocol/side/x/gmm/types"
)

var _ types.QueryServer = Keeper{}
