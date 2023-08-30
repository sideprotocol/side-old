package keeper

import (
	"github.com/sideprotocol/side/x/poolmanager/types"
)

var _ types.QueryServer = Keeper{}
