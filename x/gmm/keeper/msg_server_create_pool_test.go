package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/sideprotocol/side/app"
	"github.com/sideprotocol/side/testutil/sample"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgCreatePool() {
	var msg *types.MsgCreatePool
	suite.SetupTest()

	msg = types.NewMsgCreatePool(
		types.Alice,
		types.PoolParams{
			Type:    types.PoolType_WEIGHT,
			SwapFee: sdkmath.LegacyDec(sdk.NewInt(100)),
		},
		[]types.PoolAsset{
			types.PoolAsset{
				Token:   sdk.NewCoin(simapp.DefaultBondDenom, sdkmath.NewInt(100)),
				Weight:  sdk.NewInt(50),
				Decimal: sdk.NewInt(6),
			},
			types.PoolAsset{
				Token:   sdk.NewCoin(simapp.AltDenom, sdkmath.NewInt(100)),
				Weight:  sdk.NewInt(50),
				Decimal: sdk.NewInt(6),
			},
		},
	)

	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.msgServer.CreatePool(ctx, msg)

	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	// Check pool is created or not
	pool, err := suite.queryClient.Pool(ctx, &types.QueryPoolRequest{
		PoolId: res.PoolId,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(pool.Pool.PoolId, res.PoolId)
}

func (suite *KeeperTestSuite) TestMsgCreatePoolFail() {
	var msg *types.MsgCreatePool
	suite.SetupTest()

	testCases := []struct {
		name   string
		mallet func()
	}{
		{
			"invalid sender",
			func() {},
		},
		{
			"not enough funds",
			func() {
				msg = types.NewMsgCreatePool(
					sample.AccAddress(),
					types.PoolParams{
						Type:    types.PoolType_WEIGHT,
						SwapFee: sdkmath.LegacyDec(sdk.NewInt(100)),
					},
					[]types.PoolAsset{
						types.PoolAsset{
							Token:   sdk.NewCoin("aside", sdkmath.NewInt(100)),
							Weight:  sdk.NewInt(50),
							Decimal: sdk.NewInt(6),
						},
						types.PoolAsset{
							Token:   sdk.NewCoin("usdc", sdkmath.NewInt(100)),
							Weight:  sdk.NewInt(50),
							Decimal: sdk.NewInt(6),
						},
					},
				)

			},
		},
	}

	for _, tc := range testCases {

		msg = types.NewMsgCreatePool(
			"",
			types.PoolParams{},
			[]types.PoolAsset{},
		)
		tc.mallet()

		res, err := suite.msgServer.CreatePool(sdk.WrapSDKContext(suite.ctx), msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}
}

// Helper function to create a new pool
func (suite *KeeperTestSuite) CreateNewPool() string {
	var msg *types.MsgCreatePool
	suite.SetupTest()

	msg = types.NewMsgCreatePool(
		types.Alice,
		types.PoolParams{
			Type:    types.PoolType_WEIGHT,
			SwapFee: sdkmath.LegacyDec(sdk.NewInt(100)),
		},
		[]types.PoolAsset{
			types.PoolAsset{
				Token:   sdk.NewCoin("aside", sdkmath.NewInt(100)),
				Weight:  sdk.NewInt(50),
				Decimal: sdk.NewInt(6),
			},
			types.PoolAsset{
				Token:   sdk.NewCoin("usdc", sdkmath.NewInt(100)),
				Weight:  sdk.NewInt(50),
				Decimal: sdk.NewInt(6),
			},
		},
	)

	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.msgServer.CreatePool(ctx, msg)

	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	return res.PoolId
}
