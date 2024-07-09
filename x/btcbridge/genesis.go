package btcbridge

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/btcbridge/keeper"
	"github.com/sideprotocol/side/x/btcbridge/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	k.SetBestBlockHeader(ctx, genState.BestBlockHeader)
	if len(genState.BlockHeaders) > 0 {
		err := k.SetBlockHeaders(ctx, genState.BlockHeaders)
		if err != nil {
			panic(err)
		}
	}
	// import utxos
	for _, utxo := range genState.Utxos {
		k.SetUTXO(ctx, utxo)
	}
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.BestBlockHeader = k.GetBestBlockHeader(ctx)
	genesis.BlockHeaders = k.GetAllBlockHeaders(ctx)
	genesis.Utxos = k.GetAllUTXOs(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
