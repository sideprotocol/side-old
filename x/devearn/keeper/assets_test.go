package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "sidechain/testutil/keeper"
	"sidechain/testutil/nullify"
	"sidechain/x/devearn/keeper"
	"sidechain/x/devearn/types"
)

func createNAssets(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Assets {
	items := make([]types.Assets, n)
	for i := range items {
		items[i].Id = keeper.AppendAssets(ctx, items[i])
	}
	return items
}

func TestAssetsGet(t *testing.T) {
	keeper, ctx := keepertest.DevearnKeeper(t)
	items := createNAssets(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetAssets(ctx, item.Id)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestAssetsRemove(t *testing.T) {
	keeper, ctx := keepertest.DevearnKeeper(t)
	items := createNAssets(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveAssets(ctx, item.Id)
		_, found := keeper.GetAssets(ctx, item.Id)
		require.False(t, found)
	}
}

func TestAssetsGetAll(t *testing.T) {
	keeper, ctx := keepertest.DevearnKeeper(t)
	items := createNAssets(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllAssets(ctx)),
	)
}

func TestAssetsCount(t *testing.T) {
	keeper, ctx := keepertest.DevearnKeeper(t)
	items := createNAssets(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetAssetsCount(ctx))
}
