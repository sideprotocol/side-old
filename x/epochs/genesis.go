package epochs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sidechain/x/epochs/keeper"
	"sidechain/x/epochs/types"
	"time"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// set epoch info from genesis
	for _, epoch := range genState.Epochs {
		// Initialize empty epoch values via Cosmos SDK
		if epoch.StartTime.Equal(time.Time{}) || epoch.StartTime.IsZero() {
			epoch.StartTime = ctx.BlockTime()
		}

		epoch.CurrentEpochStartHeight = ctx.BlockHeight()

		k.SetEpochInfo(ctx, epoch)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Epochs: k.AllEpochInfos(ctx),
	}
}
