package keeper

import (
	"github.com/sideprotocol/sidechain/x/devearn/types"
)

var _ types.QueryServer = Keeper{}
