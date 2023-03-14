package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sidechain/x/devearn/types"
)

func (suite *KeeperTestSuite) TestDistributeIncentives() {
	const (
		mintAmount   int64  = 100
		gasUsed      uint64 = 500
		totalGasUsed        = gasUsed * 2
	)

	// check module balance
	moduleAddr := suite.app.AccountKeeper.GetModuleAddress(types.ModuleName)

	testCases := []struct {
		name       string
		epochs     uint32
		denom      string
		mintAmount int64
		expPass    bool
	}{
		{
			"pass - with capped reward",
			epochs,
			denomMint,
			1000000,
			true,
		},
		{
			"pass - with non-mint denom and no remaining epochs",
			1,
			denomMint,
			mintAmount,
			true,
		},
		{
			"pass - with non-mint denom and remaining epochs",
			1,
			denomMint,
			1000000,
			true,
		},
		{
			"pass - with mint denom and remaining epochs",
			epochs,
			denomMint,
			mintAmount,
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			suite.deployContracts()

			// Mint tokens in module account
			err := suite.app.BankKeeper.MintCoins(
				suite.ctx,
				types.ModuleName,
				sdk.Coins{sdk.NewInt64Coin(tc.denom, tc.mintAmount)},
			)
			suite.Require().NoError(err)

			// create incentive
			_, err = suite.app.DevearnKeeper.RegisterDevEarn(
				suite.ctx,
				contract,
				tc.epochs,
				suite.priv.PubKey().Address().String(),
			)
			suite.Require().NoError(err)

			regIn, found := suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract)
			suite.Require().True(found)
			suite.Require().Zero(regIn.GasMeter)

			balance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, tc.denom)
			suite.Require().True(balance.IsPositive())

			// Set total gas meter
			suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, types.NewDevEarn(contract, gasUsed, tc.epochs, ownerPriv1.PubKey().Address().String()))
			suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, types.NewDevEarn(contract2, gasUsed, tc.epochs, ownerPriv2.PubKey().Address().String()))

			suite.Commit()
			err = suite.app.DevearnKeeper.DistributeRewards(suite.ctx)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)

				// distributes the rewards to all participants
				sdkParticipant := sdk.AccAddress(ownerPriv1.PubKey().Address().Bytes())
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdkParticipant, tc.denom)
				gasRatio := sdk.NewDec(int64(gasUsed)).QuoInt64(int64(totalGasUsed))
				coinAllocated := sdk.NewDec(tc.mintAmount)
				expBalance := coinAllocated.Mul(gasRatio)
				suite.Require().Equal(expBalance.TruncateInt(), balance.Amount, tc.name)

				// updates the remaining epochs of each incentive and sets the cumulative
				// totalGas to zero OR deletes incentive
				regIn, found = suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract)
				if regIn.IsActive() {
					suite.Require().True(found)
					suite.Require().Equal(tc.epochs-1, regIn.Epochs)
					suite.Require().Zero(regIn.GasMeter)
				} else {
					suite.Require().False(found)
				}

			} else {
				suite.Require().Error(err)
			}
		})
	}
}
