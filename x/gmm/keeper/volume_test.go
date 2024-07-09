package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/sideprotocol/side/testutil/keeper"
	"github.com/sideprotocol/side/x/gmm/types"
	"github.com/stretchr/testify/require"
)

func TestVolumeQuery(t *testing.T) {
	keeper, ctx := testkeeper.GmmKeeper(t)
	amp := sdkmath.NewInt(100)
	// params := types.DefaultParams()
	mockAssets := []types.PoolAsset{}
	weight := sdkmath.NewInt(6)
	tokenIn := sdk.NewCoin("usdt", sdk.NewInt(100))
	tokenOut := sdk.NewCoin("usdc", sdk.NewInt(80))
	mockAssets = append(mockAssets, types.PoolAsset{
		Decimal: sdkmath.NewInt(6),
		Weight:  &weight,
		Token:   sdk.NewCoin("usdt", sdk.NewInt(1000000)),
	})
	mockAssets = append(mockAssets, types.PoolAsset{
		Decimal: sdkmath.NewInt(6),
		Weight:  &weight,
		Token:   sdk.NewCoin("usdc", sdk.NewInt(1000000)),
	})

	pool := types.Pool{
		PoolId: "test",
		Sender: types.Alice,
		PoolParams: types.PoolParams{
			ExitFee: sdkmath.LegacyDec(sdkmath.NewInt(100)),
			SwapFee: sdkmath.LegacyDec(sdkmath.NewInt(100)),
			Amp:     &amp,
		},
		Assets: mockAssets,
	}
	keeper.AppendPool(ctx, pool)
	// err := keeper.ObserveVolumeFromPool(ctx, "test", tokenIn, tokenOut)
	// require.NoError(t, err)
	// Loop to simulate 1000 observations at different timestamps
	for i := 0; i < 1000; i++ {
		// Increment the block time for each observation
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(60 * time.Minute))
		// Perform the observation
		err := keeper.ObserveVolumeFromPool(ctx, pool.PoolId, tokenIn, tokenOut)
		require.NoError(t, err)
	}

	totalVolume := keeper.GetTotalVolume(ctx, pool.PoolId)

	// Calculate expected total volume
	expectedTotalVolumeUsdt := sdk.NewInt(100 * 1000) // 100 USDT per observation * 1000 observations
	expectedTotalVolumeUsdc := sdk.NewInt(80 * 1000)  // 80 USDC per observation * 1000 observations

	// Assert total volume
	require.Len(t, totalVolume, 2, "Total volume should have two coins")
	require.Equal(t, expectedTotalVolumeUsdc, totalVolume[0].Amount, "Total USDT volume does not match")
	require.Equal(t, expectedTotalVolumeUsdt, totalVolume[1].Amount, "Total USDC volume does not match")

	volumeInDay := keeper.GetVolume24(ctx, pool.PoolId)

	// Calculate expected 24-hour volume
	expected24HourVolumeUsdt := sdk.NewInt(100 * 24) // 100 USDT per observation * 24 observations
	expected24HourVolumeUsdc := sdk.NewInt(80 * 24)  // 80 USDC per observation * 24 observations

	// Assert 24-hour volume
	require.Len(t, volumeInDay, 2, "24-hour volume should have two coins")
	require.Equal(t, expected24HourVolumeUsdc, volumeInDay[0].Amount, "24-hour USDT volume does not match")
	require.Equal(t, expected24HourVolumeUsdt, volumeInDay[1].Amount, "24-hour USDC volume does not match")
}
