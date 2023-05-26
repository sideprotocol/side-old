package keeper_test

import (
	"fmt"
	"sidechain/x/devearn/types"
	epochstypes "sidechain/x/epochs/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestEpochIdentifierAfterEpochEnd() {
	testCases := []struct {
		name            string
		epochIdentifier string
		epochs          uint32
		denom           string
	}{
		{
			"correct epoch identifier",
			epochstypes.WeekEpochID,
			epochs,
			denomMint,
		},
		{
			"incorrect epoch identifier",
			epochstypes.WeekEpochID,
			epochs,
			denomMint,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()
			suite.deployContracts()
			// Deploy Contract
			_, err := suite.DeployContract("COIN TOKEN", "COIN", erc20Decimals)
			suite.Require().NoError(err)
			suite.Commit()
			_, err = suite.app.DevearnKeeper.RegisterDevEarnInfo(
				suite.ctx,
				contract,
				tc.epochs,
				ownerPriv1.PubKey().Address().String(),
			)
			// err = suite.app.BankKeeper.MintCoins(
			// 	suite.ctx,
			// 	types.ModuleName,
			// 	sdk.Coins{sdk.NewCoin(tc.denom, sdk.NewInt(10000))},
			// )

			suite.Require().NoError(err)
			regIn, found := suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract)
			suite.Require().True(found)
			suite.Require().Zero(regIn.GasMeter)
			suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, types.NewDevEarn(contract, 1000, tc.epochs, ownerPriv1.PubKey().Address().String()))
			suite.Commit()
			// Check used gas
			regIn, found = suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract)
			suite.Require().Equal(uint32(10), regIn.Epochs)
			suite.Require().True(found)
			suite.Require().Equal(uint64(1000), regIn.GasMeter)

			params := suite.app.DevearnKeeper.GetParams(suite.ctx)
			params.EnableDevEarn = true
			err = suite.app.DevearnKeeper.SetParams(suite.ctx, params)
			suite.Require().NoError(err)
			suite.Require().Equal("week", params.RewardEpochIdentifier)

			futureCtx := suite.ctx.WithBlockTime(time.Now().Add(time.Hour))
			newHeight := suite.app.LastBlockHeight() + 1

			suite.app.EpochsKeeper.BeforeEpochStart(futureCtx, tc.epochIdentifier, newHeight)
			suite.app.EpochsKeeper.AfterEpochEnd(futureCtx, tc.epochIdentifier, newHeight)

			// Epoch hook call is working
			params = suite.app.DevearnKeeper.GetParams(suite.ctx)
			suite.Require().Equal(uint64(1000), params.TvlShare)
			regIn, found = suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract)
			suite.Require().Equal(uint32(9), regIn.Epochs)

			balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(ownerPriv1.PubKey().Address()), tc.denom)
			if tc.epochIdentifier == params.RewardEpochIdentifier {
				profit := sdk.NewDec(10000).Quo(sdk.NewDec(365)).Mul(sdk.NewDec(7)).Mul(params.DevEarnInflation_APR).TruncateInt64()
				suite.Require().Equal(profit+10000, balance.Amount.Int64())
			}
		})
	}
}
