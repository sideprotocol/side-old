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

func TestAPRCalculation(t *testing.T) {
	keeper, ctx := testkeeper.GmmKeeper(t)
	amp := sdkmath.NewInt(100)
	// params := types.DefaultParams()
	mockAssets := []types.PoolAsset{}
	weight := sdkmath.NewInt(6)
	tokenIn := sdk.NewCoin("usdt", sdk.NewInt(100))
	// tokenOut := sdk.NewCoin("usdc", sdk.NewInt(80))
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
	// Loop to simulate 1000 observations at different timestamps
	for i := 0; i < 365; i++ {
		// Perform the observation
		err := keeper.ObserveFeeFromPool(ctx, pool.PoolId, tokenIn)
		// Increment the block time for each observation
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(24 * time.Hour))
		require.NoError(t, err)
	}

	// Calculate the APR for the pool
	apr := keeper.GetAPR(ctx, pool.PoolId)
	expectedAPR := sdk.NewCoin("usdt", sdkmath.NewInt(6500000))
	// Assert APR calculation
	require.Equal(t, expectedAPR.Amount.LTE(apr[0].Amount), true, "Calculated APR does not match expected APR")
}
