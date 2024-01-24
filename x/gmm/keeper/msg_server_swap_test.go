package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/sideprotocol/side/app"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgSwap() {
	suite.SetupTest()

	tests := []struct {
		name     string
		poolType types.PoolType
		mutator  func(*types.MsgSwap, string)
	}{
		{
			"swap in weight pool",
			types.PoolType_WEIGHT,
			func(msg *types.MsgSwap, poolID string) {
				msg.TokenIn = sdk.NewCoin(simapp.DefaultBondDenom, sdk.NewInt(100))
				msg.TokenOut = sdk.NewCoin(simapp.USDC, sdk.NewInt(0))
			},
		},
		{
			"swap in stable pool",
			types.PoolType_STABLE,
			func(msg *types.MsgSwap, poolID string) {
				msg.TokenIn = sdk.NewCoin(simapp.WDAI, sdk.NewInt(100))
				msg.TokenOut = sdk.NewCoin(simapp.WUSDT, sdk.NewInt(0))
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			// Create a new pool of the specific type
			poolID := suite.CreateNewPool(tc.poolType)

			// Initialize the MsgSwap
			msg := types.MsgSwap{
				Sender:   types.Alice,
				PoolId:   poolID,
				Slippage: sdkmath.NewInt(1),
			}

			// Use the mutator to set the token in and token out for the specific pool type
			tc.mutator(&msg, poolID)

			ctx := sdk.WrapSDKContext(suite.ctx)

			// Query the pool before the swap
			queryResBeforeSwap, err := suite.queryClient.Pool(ctx, &types.QueryPoolRequest{
				PoolId: poolID,
			})
			suite.Require().NoError(err)

			pool := queryResBeforeSwap.Pool.ToPool()
			outAssetBeforeSwap := pool.Assets[msg.TokenOut.Denom]
			estimatedOut, err := pool.EstimateSwap(msg.TokenIn, msg.TokenOut.Denom)
			msg.TokenOut = estimatedOut

			suite.Require().NoError(err)

			// Perform the swap
			res, err := suite.msgServer.Swap(ctx, &msg)
			suite.Require().NoError(err)
			suite.Require().NotNil(res)

			// Query the pool after the swap
			queryResAfterSwap, err := suite.queryClient.Pool(ctx, &types.QueryPoolRequest{
				PoolId: poolID,
			})
			suite.Require().NoError(err)
			pool = queryResAfterSwap.Pool.ToPool()
			outAssetAfterSwap := pool.Assets[msg.TokenOut.Denom]
			out := outAssetBeforeSwap.Token.Sub(outAssetAfterSwap.Token)
			suite.Require().Equal(out, estimatedOut)
		})
	}
}

func (suite *KeeperTestSuite) TestMsgSwapFail() {
	suite.SetupTest()

	tests := []struct {
		name     string
		poolType types.PoolType
		mutator  func(*types.MsgSwap, string)
	}{
		{
			"swap with invalid pool ID",
			types.PoolType_WEIGHT,
			func(msg *types.MsgSwap, poolID string) {
				msg.PoolId = "invalid_pool_id"
			},
		},
		{
			"swap with invalid sender",
			types.PoolType_WEIGHT,
			func(msg *types.MsgSwap, poolID string) {
				msg.Sender = "invalid_sender_address"
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			// Create a new pool of the specific type, if necessary
			poolID := suite.CreateNewPool(tc.poolType)

			// Initialize the MsgSwap
			msg := types.MsgSwap{
				Sender:   types.Alice,
				PoolId:   poolID,
				Slippage: sdkmath.NewInt(1),
				TokenIn:  sdk.NewCoin(simapp.WDAI, sdk.NewInt(100)),
				TokenOut: sdk.NewCoin(simapp.WUSDT, sdk.NewInt(50)),
			}

			// Use the mutator to set the failure conditions
			tc.mutator(&msg, poolID)

			ctx := sdk.WrapSDKContext(suite.ctx)

			// Perform the swap
			res, err := suite.msgServer.Swap(ctx, &msg)

			// Expect an error and nil result
			suite.Require().Error(err)
			suite.Require().Nil(res)
		})
	}
}
