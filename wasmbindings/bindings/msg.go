package bindings

import sdk "github.com/cosmos/cosmos-sdk/types"

type OsmosisMsg struct {
	/// Contracts can create swap using gmm pools
	Swap *Swap `json:"swap,omitempty"`
}

type Swap struct {
	PoolId   string   `json:"pool_id"`
	TokenIn  sdk.Coin `json:"token_in"`
	TokenOut sdk.Coin `json:"token_out"`
}
