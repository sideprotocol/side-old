package keeper

import (
	"github.com/sideprotocol/side/x/router/types"
)

var _ types.QueryServer = Keeper{}
