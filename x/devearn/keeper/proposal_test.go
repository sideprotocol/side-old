package keeper_test

import (
	"fmt"
	utiltx "github.com/sideprotocol/sidechain/testutil/tx"
	"github.com/sideprotocol/sidechain/x/devearn/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite KeeperTestSuite) TestRegisterIncentive() { //nolint:govet // we can copy locks here because it is a test
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"dev_earn are disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableDevEarn = false
				suite.app.DevearnKeeper.SetParams(suite.ctx, params) //nolint:errcheck
			},
			false,
		},
		{
			"contract doesn't exist",
			func() {
				contract = utiltx.GenerateAddress()
			},
			false,
		},
		{
			"inventive already registered",
			func() {
				regIn := types.NewDevEarn(contract, 0, epochs, suite.priv.PubKey().Address().String())
				suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, regIn)
				suite.Commit()
			},
			false,
		},
		{
			"ok",
			func() {
				// Make sure the non-mint coin has supply
				err := suite.app.BankKeeper.MintCoins(
					suite.ctx,
					types.ModuleName,
					sdk.Coins{sdk.NewInt64Coin(denomMint, 1)},
				)
				suite.Require().NoError(err)
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			suite.deployContracts()

			tc.malleate()

			in, err := suite.app.DevearnKeeper.RegisterDevEarnInfo(
				suite.ctx,
				contract,
				epochs,
				suite.priv.PubKey().Address().String(),
			)
			suite.Commit()

			expIn := &types.DevEarnInfo{
				Contract:     contract.String(),
				Epochs:       epochs,
				GasMeter:     0,
				OwnerAddress: suite.priv.PubKey().Address().String(),
				StartTime:    suite.ctx.BlockTime(),
			}

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(expIn, in)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite KeeperTestSuite) TestCancelIncentive() { //nolint:govet // we can copy locks here because it is a test
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"incentives are disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableDevEarn = false
				suite.app.DevearnKeeper.SetParams(suite.ctx, params) //nolint:errcheck
			},
			false,
		},
		{
			"inventive not registered",
			func() {
			},
			false,
		},
		{
			"ok",
			func() {
				_, err := suite.app.DevearnKeeper.RegisterDevEarnInfo(
					suite.ctx,
					contract,
					epochs,
					suite.priv.PubKey().Address().String(),
				)
				suite.Require().NoError(err)
				suite.Commit()
			},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			suite.deployContracts()

			tc.malleate()

			err := suite.app.DevearnKeeper.CancelDevEarnInfo(suite.ctx, contract)
			suite.Commit()

			_, ok := suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().False(ok, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
				suite.Require().False(ok, tc.name)
			}
		})
	}
}

// TODO: Add asset add and remove tests
