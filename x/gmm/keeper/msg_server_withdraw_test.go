package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgWithdraw() {
	suite.SetupTest()

	tests := []struct {
		name     string
		poolType types.PoolType
		mutator  func(*types.MsgWithdraw, string)
	}{
		{
			"withdraw from weight pool",
			types.PoolType_WEIGHT,
			func(msg *types.MsgWithdraw, poolID string) {
				msg.Share = sdk.NewCoin(poolID, sdk.NewInt(10))
			},
		},
		{
			"withdraw from stable pool",
			types.PoolType_STABLE,
			func(msg *types.MsgWithdraw, poolID string) {
				msg.Share = sdk.NewCoin(poolID, sdk.NewInt(200))
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			// Create a new pool of the specific type
			poolID := suite.CreateNewPool(tc.poolType)

			// Initialize the MsgWithdraw
			msg := types.MsgWithdraw{
				Sender:   types.Alice,
				Receiver: types.Carol,
				PoolId:   poolID,
			}

			// Use the mutator to set the share for the specific pool type
			tc.mutator(&msg, poolID)

			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.msgServer.Withdraw(ctx, &msg)

			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		})
	}
}

func (suite *KeeperTestSuite) TestMsgWithdrawFail() {
	suite.SetupTest()

	tests := []struct {
		name     string
		poolType types.PoolType
		mutator  func(*types.MsgWithdraw, string)
	}{
		{
			"withdraw from non-existent pool",
			types.PoolType_WEIGHT,
			func(msg *types.MsgWithdraw, poolID string) {
				msg.PoolId = "invalid_pool_id"
			},
		},
		{
			"withdraw with invalid sender",
			types.PoolType_WEIGHT,
			func(msg *types.MsgWithdraw, poolID string) {
				msg.Sender = "invalid_sender_address"
			},
		},
		{
			"withdraw with insufficient balance",
			types.PoolType_WEIGHT,
			func(msg *types.MsgWithdraw, poolID string) {
				msg.Share = sdk.NewCoin(poolID, sdk.NewInt(1000000)) // Some large number
			},
		},
		// Add more failure test cases to cover all scenarios
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			// Create a new pool of the specific type if necessary
			poolID := suite.CreateNewPool(tc.poolType)

			// Initialize the MsgWithdraw
			msg := types.MsgWithdraw{
				Sender:   types.Alice,
				Receiver: types.Carol,
				PoolId:   poolID,
				Share:    sdk.NewCoin(poolID, sdk.NewInt(200)),
			}

			// Use the mutator to set the failure condition
			tc.mutator(&msg, poolID)

			ctx := sdk.WrapSDKContext(suite.ctx)

			// Perform the withdraw
			res, err := suite.msgServer.Withdraw(ctx, &msg)

			// Expect an error and nil result
			suite.Require().Error(err)
			suite.Require().Nil(res)
		})
	}
}
