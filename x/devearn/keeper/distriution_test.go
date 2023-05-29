package keeper_test

import (
	"fmt"
	"math/big"
	"sidechain/contracts"
	"sidechain/x/devearn/types"

	erc20types "sidechain/x/erc20/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Distribute incentives on basis of gas used only
func (suite *KeeperTestSuite) TestDistributeIncentivesGas() {
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

			//Mint tokens in module account
			err := suite.app.BankKeeper.MintCoins(
				suite.ctx,
				types.ModuleName,
				sdk.Coins{sdk.NewInt64Coin(tc.denom, tc.mintAmount)},
			)
			suite.Require().NoError(err)

			// create incentive
			_, err = suite.app.DevearnKeeper.RegisterDevEarnInfo(
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
				params := suite.app.DevearnKeeper.GetParams(suite.ctx)
				// distributes the rewards to all participants
				sdkParticipant := sdk.AccAddress(ownerPriv1.PubKey().Address().Bytes())
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdkParticipant, tc.denom)
				gasRatio := sdk.NewDec(int64(gasUsed)).QuoInt64(int64(totalGasUsed))
				tvlAllocation := sdk.NewDec(tc.mintAmount).Mul(sdk.NewDecFromBigInt(new(big.Int).SetUint64(params.TvlShare)))
				tvlAllocation = tvlAllocation.Quo(sdk.NewDec(10000))
				gasAllocation := sdk.NewDec(tc.mintAmount).Sub(tvlAllocation)
				expBalance := gasAllocation.Mul(gasRatio)
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

// Distribute incentives on basis of gas used and TVL
// Register erc20 to erc20module
// Get denom from registered pair
// Use denom in oracle
// Add denom in asset whitelist
// TVL should account all of these
func (suite *KeeperTestSuite) TestDistributeIncentivesTVLandGas() {
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

			//Mint tokens in module account
			err := suite.app.BankKeeper.MintCoins(
				suite.ctx,
				types.ModuleName,
				sdk.Coins{sdk.NewInt64Coin(tc.denom, tc.mintAmount)},
			)
			suite.Require().NoError(err)

			// Register first contract
			_, err = suite.app.DevearnKeeper.RegisterDevEarnInfo(
				suite.ctx,
				contract,
				tc.epochs,
				suite.priv.PubKey().Address().String(),
			)
			suite.Require().NoError(err)

			// Register second contract
			_, err = suite.app.DevearnKeeper.RegisterDevEarnInfo(
				suite.ctx,
				contract2,
				tc.epochs,
				ownerPriv1.PubKey().Address().String(),
			)
			suite.Require().NoError(err)

			regIn, found := suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract)
			suite.Require().True(found)
			suite.Require().Zero(regIn.GasMeter)

			regIn, found = suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract2)
			suite.Require().True(found)
			suite.Require().Zero(regIn.GasMeter)

			balance := suite.app.BankKeeper.GetBalance(suite.ctx, moduleAddr, tc.denom)
			suite.Require().True(balance.IsPositive())

			// Set total gas meter
			suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, types.NewDevEarn(contract, gasUsed, tc.epochs, ownerPriv1.PubKey().Address().String()))
			suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, types.NewDevEarn(contract2, gasUsed, tc.epochs, ownerPriv2.PubKey().Address().String()))

			// Mint tokens to contracts
			// Mint tokens and send to receiver
			erc20Abi := contracts.ERC20MinterBurnerDecimalsContract.ABI
			_, err = suite.app.Erc20Keeper.CallEVM(suite.ctx, erc20Abi, erc20types.ModuleAddress, contract, true, "mint", contract2, 100000)
			_, err = suite.app.Erc20Keeper.CallEVM(suite.ctx, erc20Abi, erc20types.ModuleAddress, contract2, true, "mint", contract, 100000)

			// Register erc20 token
			// contract reg
			expPair := erc20types.NewTokenPair(contract, "coin", true, erc20types.OWNER_MODULE)
			id := expPair.GetID()
			suite.app.Erc20Keeper.SetTokenPair(suite.ctx, expPair)
			suite.app.Erc20Keeper.SetDenomMap(suite.ctx, expPair.Denom, id)
			suite.app.Erc20Keeper.SetERC20Map(suite.ctx, expPair.GetERC20Contract(), id)
			// contract2 reg
			expPair = erc20types.NewTokenPair(contract2, "coin2", true, erc20types.OWNER_MODULE)
			id = expPair.GetID()
			suite.app.Erc20Keeper.SetTokenPair(suite.ctx, expPair)
			suite.app.Erc20Keeper.SetDenomMap(suite.ctx, expPair.Denom, id)
			suite.app.Erc20Keeper.SetERC20Map(suite.ctx, expPair.GetERC20Contract(), id)

			// Add denom to oracle
			// Add assets in devearn
			suite.app.DevearnKeeper.AddAssetToWhitelist(suite.ctx, "coin")
			suite.app.DevearnKeeper.AddAssetToWhitelist(suite.ctx, "coin2")

			suite.Commit()
			err = suite.app.DevearnKeeper.DistributeRewards(suite.ctx)

			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				params := suite.app.DevearnKeeper.GetParams(suite.ctx)
				// distributes the rewards to all participants
				sdkParticipant := sdk.AccAddress(ownerPriv1.PubKey().Address().Bytes())
				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdkParticipant, tc.denom)
				gasRatio := sdk.NewDec(int64(gasUsed)).QuoInt64(int64(totalGasUsed))
				tvlAllocation := sdk.NewDec(tc.mintAmount).Mul(sdk.NewDecFromBigInt(new(big.Int).SetUint64(params.TvlShare)))
				tvlAllocation = tvlAllocation.Quo(sdk.NewDec(10000))
				gasAllocation := sdk.NewDec(tc.mintAmount).Sub(tvlAllocation)
				expBalance := gasAllocation.Mul(gasRatio)
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
