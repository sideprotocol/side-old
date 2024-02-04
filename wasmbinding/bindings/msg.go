package bindings

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SideMsg struct {
	/// Contracts can create swap using gmm pools
	Swap *Swap `json:"swap,omitempty"`
}

type Swap struct {
	PoolId   string      `json:"pool_id"`
	TokenIn  sdk.Coin    `json:"token_in"`
	TokenOut sdk.Coin    `json:"token_out"`
	Slippage sdkmath.Int `json:"slippage"`
}
