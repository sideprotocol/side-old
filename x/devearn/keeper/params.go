package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sidechain/x/devearn/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.EnableDevEarn(ctx),
		k.DevEarnEpoch(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// EnableDevEarn returns the EnableDevEarn param
func (k Keeper) EnableDevEarn(ctx sdk.Context) (res bool) {
	k.paramstore.Get(ctx, types.KeyEnableDevEarn, &res)
	return
}

// DevEarnEpoch returns the DevEarnEpoch param
func (k Keeper) DevEarnEpoch(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyDevEarnEpoch, &res)
	return
}
