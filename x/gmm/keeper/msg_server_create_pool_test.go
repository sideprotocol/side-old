package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/sideprotocol/side/app"
	"github.com/sideprotocol/side/testutil/sample"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgCreatePool() {
	suite.SetupTest()

	tests := []struct {
		name    string
		mutator func(msg *types.MsgCreatePool)
	}{
		{
			"weight pool",
			func(msg *types.MsgCreatePool) {
				msg.Params.Type = types.PoolType_WEIGHT
			},
		},
		{
			"stable pool",
			func(msg *types.MsgCreatePool) {
				amp := sdkmath.NewInt(100)
				msg.Params.Type = types.PoolType_STABLE
				msg.Params.Amp = &amp
				msg.Liquidity = suite.createStablePoolLiquidity()
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			msg := suite.defaultMsgCreatePool()
			tc.mutator(msg)

			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.msgServer.CreatePool(ctx, msg)

			suite.Require().NoError(err)
			suite.Require().NotNil(res)

			// Check pool is created or not
			pool, err := suite.queryClient.Pool(ctx, &types.QueryPoolRequest{
				PoolId: res.PoolId,
			})
			suite.Require().NoError(err)
			suite.Require().Equal(pool.Pool.Id, res.PoolId)
		})
	}
}

// Helper method to create a default MsgCreatePool for testing
func (suite *KeeperTestSuite) defaultMsgCreatePool() *types.MsgCreatePool {
	weight := sdkmath.NewInt(50)
	amp := sdkmath.NewInt(100)
	return types.NewMsgCreatePool(
		types.Alice,
		types.PoolParams{
			Type:    types.PoolType_WEIGHT,
			SwapFee: sdkmath.LegacyDec(sdkmath.NewInt(100)),
			Amp:     &amp,
		},
		[]types.PoolAsset{
			{
				Token:   sdk.NewCoin(simapp.DefaultBondDenom, sdkmath.NewInt(100)),
				Weight:  &weight,
				Decimal: sdkmath.NewInt(6),
			},
			{
				Token:   sdk.NewCoin(simapp.USDC, sdkmath.NewInt(100)),
				Weight:  &weight,
				Decimal: sdkmath.NewInt(6),
			},
		},
	)
}

// Helper method to create stable pool liquidity for testing
func (suite *KeeperTestSuite) createStablePoolLiquidity() []types.PoolAsset {
	weight := sdkmath.NewInt(50)
	return []types.PoolAsset{
		{
			Token:   sdk.NewCoin(simapp.WDAI, sdkmath.NewInt(100)),
			Weight:  &weight,
			Decimal: sdkmath.NewInt(6),
		},
		{
			Token:   sdk.NewCoin(simapp.WUSDT, sdkmath.NewInt(100)),
			Weight:  &weight,
			Decimal: sdkmath.NewInt(6),
		},
	}
}

func (suite *KeeperTestSuite) TestMsgCreatePoolFail() {
	var msg *types.MsgCreatePool
	suite.SetupTest()
	weight := sdkmath.NewInt(50)
	amp := sdkmath.NewInt(100)

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
						SwapFee: sdkmath.LegacyDec(sdkmath.NewInt(100)),
						Amp:     &amp,
					},
					[]types.PoolAsset{
						{
							Token:   sdk.NewCoin(simapp.DefaultBondDenom, sdkmath.NewInt(100)),
							Weight:  &weight,
							Decimal: sdkmath.NewInt(6),
						},
						{
							Token:   sdk.NewCoin(simapp.USDC, sdkmath.NewInt(100)),
							Weight:  &weight,
							Decimal: sdkmath.NewInt(6),
						},
					},
				)
			},
		},
		{
			"create pool with zero liquidity",
			func() {
				msg = types.NewMsgCreatePool(
					sample.AccAddress(),
					types.PoolParams{
						Type:    types.PoolType_WEIGHT,
						SwapFee: sdkmath.LegacyDec(sdkmath.NewInt(100)),
						Amp:     &amp,
					},
					[]types.PoolAsset{},
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

func (suite *KeeperTestSuite) CreateNewPool(poolType types.PoolType) string {
	switch poolType {
	case types.PoolType_STABLE:
		return suite.createNewStablePool()
	default:
		return suite.createNewWeightPool()
	}
}

// Helper function to create a new pool
func (suite *KeeperTestSuite) createNewWeightPool() string {
	var msg *types.MsgCreatePool
	suite.SetupTest()

	weight := sdkmath.NewInt(50)
	amp := sdkmath.NewInt(100)

	msg = types.NewMsgCreatePool(
		types.Alice,
		types.PoolParams{
			Type:    types.PoolType_WEIGHT,
			SwapFee: sdkmath.LegacyDec(sdkmath.NewInt(100)),
			Amp:     &amp,
		},
		[]types.PoolAsset{
			{
				Token:   sdk.NewCoin(simapp.DefaultBondDenom, sdkmath.NewInt(100)),
				Weight:  &weight,
				Decimal: sdkmath.NewInt(6),
			},
			{
				Token:   sdk.NewCoin(simapp.USDC, sdkmath.NewInt(100)),
				Weight:  &weight,
				Decimal: sdkmath.NewInt(6),
			},
		},
	)

	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.msgServer.CreatePool(ctx, msg)

	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	return res.PoolId
}

func (suite *KeeperTestSuite) createNewStablePool() string {
	var msg *types.MsgCreatePool
	suite.SetupTest()

	weight := sdkmath.NewInt(50)
	amp := sdkmath.NewInt(1)

	msg = types.NewMsgCreatePool(
		types.Alice,
		types.PoolParams{
			Type:    types.PoolType_STABLE,
			SwapFee: sdkmath.LegacyDec(sdkmath.NewInt(100)),
			Amp:     &amp,
		},
		[]types.PoolAsset{
			{
				Token:   sdk.NewCoin(simapp.WDAI, sdkmath.NewInt(100)),
				Weight:  &weight,
				Decimal: sdkmath.NewInt(6),
			},
			{
				Token:   sdk.NewCoin(simapp.WUSDT, sdkmath.NewInt(100)),
				Weight:  &weight,
				Decimal: sdkmath.NewInt(6),
			},
		},
	)

	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.msgServer.CreatePool(ctx, msg)

	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	return res.PoolId
}
