package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/sideprotocol/side/app"
	"github.com/sideprotocol/side/testutil/sample"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgAddLiquidity() {
	suite.SetupTest()

	tests := []struct {
		name     string
		poolType types.PoolType
		mutator  func(*types.MsgAddLiquidity, string)
	}{
		{
			"add liquidity to weight pool",
			types.PoolType_WEIGHT,
			func(msg *types.MsgAddLiquidity, poolID string) {
				msg.Liquidity = []sdk.Coin{
					sdk.NewCoin(simapp.DefaultBondDenom, sdk.NewInt(100)),
				}
			},
		},
		{
			"add liquidity to stable pool",
			types.PoolType_STABLE,
			func(msg *types.MsgAddLiquidity, poolID string) {
				msg.Liquidity = []sdk.Coin{
					sdk.NewCoin(simapp.WDAI, sdk.NewInt(100)),
					sdk.NewCoin(simapp.WUSDT, sdk.NewInt(100)),
				}
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			// Create a new pool of the specific type
			poolID := suite.CreateNewPool(tc.poolType)

			// Initialize the MsgAddLiquidity
			msg := types.MsgAddLiquidity{
				Sender: types.Carol,
				PoolId: poolID,
			}

			// Use the mutator to set the liquidity for the specific pool type
			tc.mutator(&msg, poolID)

			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.msgServer.AddLiquidity(ctx, &msg)

			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		})
	}
}

func (suite *KeeperTestSuite) TestMsgAddLiquidityFail() {
	var msg *types.MsgAddLiquidity
	// Create a new pool
	poolID := suite.CreateNewPool(types.PoolType_WEIGHT)

	testCases := []struct {
		name   string
		mallet func(msg *types.MsgAddLiquidity, poolID string)
	}{
		{
			"invalid sender",
			func(msg *types.MsgAddLiquidity, poolID string) {},
		},
		{
			"invalid poolID",
			func(msg *types.MsgAddLiquidity, poolID string) {
				msg = &types.MsgAddLiquidity{
					Sender: sample.AccAddress(),
					PoolId: "",
					Liquidity: []sdk.Coin{
						sdk.NewCoin(
							simapp.DefaultBondDenom,
							sdk.NewInt(100),
						),
					},
				}
			},
		},
		{
			"not enough funds",
			func(msg *types.MsgAddLiquidity, poolID string) {
				msg = &types.MsgAddLiquidity{
					Sender: sample.AccAddress(),
					PoolId: poolID,
					Liquidity: []sdk.Coin{
						sdk.NewCoin(
							simapp.DefaultBondDenom,
							sdk.NewInt(100),
						),
					},
				}
			},
		},
		{
			"invalid asset type",
			func(msg *types.MsgAddLiquidity, poolID string) {
				msg.Liquidity = []sdk.Coin{
					sdk.NewCoin("INVALID_ASSET_TYPE", sdk.NewInt(100)),
				}
			},
		},
		{
			"zero liquidity",
			func(msg *types.MsgAddLiquidity, poolID string) {
				msg.Liquidity = []sdk.Coin{
					sdk.NewCoin(simapp.DefaultBondDenom, sdk.NewInt(0)),
				}
			},
		},
	}

	for _, tc := range testCases {

		msg = types.NewMsgAddLiquidity(
			"",
			poolID,
			[]sdk.Coin{},
		)
		tc.mallet(msg, poolID)

		res, err := suite.msgServer.AddLiquidity(sdk.WrapSDKContext(suite.ctx), msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}
}
