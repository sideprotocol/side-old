package btclightclient

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/btclightclient/keeper"
	"github.com/sideprotocol/side/x/btclightclient/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	k.SetBestBlockHeader(ctx, genState.BestBlockHeader)
	if len(genState.BlockHeaders) > 0 {
		k.SetBlockHeaders(ctx, genState.BlockHeaders)
	}

}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.BestBlockHeader = k.GetBestBlockHeader(ctx)
	genesis.BlockHeaders = k.GetAllBlockHeaders(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
