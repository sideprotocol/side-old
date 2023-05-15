package keeper

import (
	"sidechain/x/devearn/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	epochstypes "sidechain/x/epochs/types"
)

// BeforeEpochStart performs a no-op
func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {}

// AfterEpochEnd distributes the contract incentives at the end of each epoch
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) {
	params := k.GetParams(ctx)

	// check if epochIdentifier signal equals the identifier in the params
	if epochIdentifier != params.RewardEpochIdentifier {
		return
	}

	// check if the Incentives are globally enabled
	if !params.EnableDevEarn {
		return
	}
	//get total supply
	totalDenomSupply, _, err := k.bankKeeper.GetPaginatedTotalSupply(ctx, &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      100,
		CountTotal: false,
		Reverse:    false,
	})
	if err != nil {
		ctx.Logger().Error("get total supply err happen, err :", err)
		return
	}
	denom, err := sdk.GetBaseDenom()
	if err != nil {
		ctx.Logger().Error("get base denom err happen, err :", err)
		return
	}

	totalSupply := totalDenomSupply.AmountOf(denom)
	periodProvision := sdk.NewDecFromBigInt(totalSupply.BigInt()).Mul(params.DevEarnInflation_APR)
	var epochProvision sdk.Dec
	switch params.RewardEpochIdentifier {
	case epochstypes.WeekEpochID:
		epochsPerPeriod := sdk.NewDec(365).Quo(sdk.NewDec(7))
		epochProvision = periodProvision.Quo(epochsPerPeriod)
	case epochstypes.DayEpochID:
		epochsPerPeriod := sdk.NewDec(365).Quo(sdk.NewDec(1))
		epochProvision = periodProvision.Quo(epochsPerPeriod)
	default:
		ctx.Logger().Error("RewardEpochIdentifier should be day or week")
		return
	}
	if !epochProvision.IsPositive() {
		k.Logger(ctx).Error(
			"SKIPPING INFLATION: negative epoch mint provision",
			"value", epochProvision.String(),
		)
		return
	}
	//mint token
	mintedCoin := sdk.Coin{
		Denom:  denom,
		Amount: epochProvision.TruncateInt(),
	}
	coins := sdk.Coins{mintedCoin}
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		k.Logger(ctx).Error(
			"SKIPPING INFLATION: mint coin err",
			"err", err,
		)
		return
	}

	//send token to contract owner
	if err := k.DistributeRewards(ctx); err != nil {
		panic(err)
	}
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for dev earn keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// BeforeEpochStart implements EpochHooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

// AfterEpochEnd implements EpochHooks
func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
