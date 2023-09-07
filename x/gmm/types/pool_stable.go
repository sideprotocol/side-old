package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) estimateShareInStablePool(coins sdk.Coins) (sdk.Coin, error) {
	// BPT out, so we round down overall.

	// First loop calculates the sum of all token balances, which will be used to calculate
	// the current weights of each token, relative to this sum
	sum := p.Sum()

	// Calculate the weighted balance ratio without considering fees
	balanceRatiosWithFee := make(map[string]sdkmath.LegacyDec, len(coins)) //new Array<BigNumber>(amountsIn.length);
	// The weighted sum of token balance ratios without fee
	invariantRatioWithFees := sdkmath.LegacyZeroDec()

	for _, asset := range coins {
		currentWeight := sdkmath.LegacyDec(asset.Amount).Quo(sdkmath.LegacyNewDecFromInt(sum))

		balanceRatiosWithFee[asset.Denom] = sdkmath.LegacyDec(asset.Amount.Add(p.Assets[asset.Denom].Token.Amount)).Quo(
			sdkmath.LegacyNewDecFromInt(p.Assets[asset.Denom].Token.Amount),
		)

		invariantRatioWithFees = invariantRatioWithFees.Add(
			balanceRatiosWithFee[asset.Denom].Mul(currentWeight),
		)
	}

	// Second loop calculates new amounts in, taking into account the fee on the percentage excess

	newBalances := sdk.NewCoins()
	for _, amountIn := range coins {
		amountInWithoutFee := sdkmath.ZeroInt()
		asset := p.Assets[amountIn.Denom]
		// Check if the balance ratio is greater than the ideal ratio to charge fees or not
		if balanceRatiosWithFee[asset.Token.Denom].GT(invariantRatioWithFees) {
			nonTaxableAmount := asset.Token.Amount.Quo(sdkmath.Int(invariantRatioWithFees).Sub(sdkmath.Int(sdkmath.LegacyOneDec())))
			taxableAmount := amountIn.Amount.Sub(nonTaxableAmount)
			remainFee := sdkmath.LegacyNewDec(10000).Sub(p.PoolParams.SwapFee).RoundInt()
			amountInWithoutFee = nonTaxableAmount.Add(taxableAmount.Mul(sdk.NewInt(10000)).Quo(remainFee))

		} else {
			amountInWithoutFee = amountIn.Amount
		}
		newBalances = append(newBalances, sdk.NewCoin(
			amountIn.Denom,
			asset.Token.Amount.Add(amountInWithoutFee),
		))
	}

	// Get current and new invariants, taking swap fees into account
	currentInvariant, err := calculateInvariantInStablePool(p.PoolParams.Amp, p.GetLiquidity())
	if err != nil {
		return sdk.NewCoin(p.PoolId, sdkmath.NewInt(0)), nil
	}

	newInvariant, err := calculateInvariantInStablePool(p.PoolParams.Amp, newBalances)
	if err != nil {
		return sdk.NewCoin(p.PoolId, sdkmath.NewInt(0)), nil
	}

	invariantRatio := sdkmath.LegacyNewDecFromInt(newInvariant).Quo(
		sdkmath.LegacyNewDecFromInt(currentInvariant),
	)

	// If the invariant didn't increase for any reason, we simply don't mint BPT
	if invariantRatio.GT(sdkmath.LegacyZeroDec()) {
		share := p.TotalShares.Amount.Mul(sdkmath.Int(invariantRatio.Sub(sdkmath.LegacyOneDec())))
		return sdk.NewCoin(p.PoolId, share), nil

	} else {
		return sdk.NewCoin(p.PoolId, sdkmath.NewInt(0)), nil
	}
}

func (p *Pool) estimateSwapInStablePool(tokenIn sdk.Coin, denomOut string) (sdk.Coin, error) {
	/**************************************************************************************************************
	  // outGivenIn token x for y - polynomial equation to solve                                                   //
	  // ay = amount out to calculate                                                                              //
	  // by = balance token out                                                                                    //
	  // y = by - ay (finalBalanceOut)                                                                             //
	  // D = invariant                                               D                     D^(n+1)                 //
	  // A = amplification coefficient               y^2 + ( S - ----------  - D) * y -  ------------- = 0         //
	  // n = number of tokens                                    (A * n^n)               A * n^2n * P              //
	  // S = sum of final balances but y                                                                           //
	  // P = product of final balances but y                                                                       //
	  **************************************************************************************************************/

	// Subtract the fee from the amount in if requested

	tokenInDec := MinusFees(tokenIn.Amount, p.PoolParams.SwapFee)

	inv, err := calculateInvariantInStablePool(p.PoolParams.Amp, p.GetLiquidity())
	if err != nil {
	}

	assets := p.Assets

	balance := assets[tokenIn.Denom].Token.Amount.Add(tokenInDec.RoundInt())
	assets[tokenIn.Denom] = PoolAsset{
		Token:   sdk.NewCoin(tokenIn.Denom, balance),
		Weight:  assets[tokenIn.Denom].Weight,
		Decimal: assets[tokenIn.Denom].Decimal,
	}

	finalBalanceOut, err := getTokenBalanceGivenInvariantAndAllOtherBalances(
		p.PoolParams.Amp, inv, assets, tokenIn.Denom,
	)
	out := p.Assets[denomOut].Token.Amount.Sub(finalBalanceOut).Sub(sdkmath.OneInt())
	return sdk.NewCoin(denomOut, out), err
}

func (p *Pool) estimateWithdrawalsFromStablePool(share sdk.Coin) ([]sdk.Coin, error) {
	/**********************************************************************************************
	// exactBPTInForTokensOut                                                                    //
	// (per token)                                                                               //
	// aO = tokenAmountOut             /        bptIn         \                                  //
	// b = tokenBalance      a0 = b * | ---------------------  |                                 //
	// bptIn = bptAmountIn             \     bptTotalSupply    /                                 //
	// bpt = bptTotalSupply                                                                      //
	**********************************************************************************************/

	// Since we're computing an amount out, we round down overall. This means rounding down on both the
	// multiplication and division.

	bptAmountIn := sdkmath.LegacyNewDecFromInt(share.Amount)
	totalShareDec := sdkmath.LegacyNewDecFromInt(p.TotalShares.Amount)
	bptRatio := bptAmountIn.Quo(totalShareDec)

	outs := sdk.NewCoins()
	for _, asset := range p.Assets {
		amountOut := sdkmath.LegacyNewDecFromInt(asset.Token.Amount).Mul(bptRatio)
		outs = outs.Add(sdk.NewCoin(
			asset.Token.Denom,
			amountOut.RoundInt(),
		))
	}

	return outs, nil
}

/**********************************************************************************************
  // invariant                                                                                 //
  // D = invariant                                                  D^(n+1)                    //
  // A = amplification coefficient      A  n^n S + D = A D n^n + -----------                   //
  // S = sum of balances                                             n^n P                     //
  // P = product of balances                                                                   //
  // n = number of tokens                                                                      //
  **********************************************************************************************/

// Math variables (assuming these are constants or have been previously defined)
var AMP_PRECISION = sdkmath.NewInt(1000)

func calculateInvariantInStablePool(
	amp sdkmath.Int,
	assets []sdk.Coin,
) (sdkmath.Int, error) {

	// Initialize sum as zero
	sum := sdkmath.NewInt(0)

	// Number of tokens
	numTokens := sdkmath.NewInt(int64(len(assets)))

	for _, asset := range assets {
		sum = sum.Add(asset.Amount)
	}

	if sum.IsZero() {
		return sum, nil
	}

	preInv := sdkmath.NewInt(0)
	inv := sum
	ampTimeTotal := amp.Mul(numTokens)

	for i := 0; i < 255; i++ {
		P_D := numTokens.Mul(assets[0].Amount)
		for _, asset := range assets {
			P_D = P_D.Mul(asset.Amount).Mul(numTokens).Quo(inv)
		}
		preInv = inv

		inv1 := numTokens.Mul(inv).Mul(inv).Add((ampTimeTotal.Mul(sum).Mul(P_D)).Quo(AMP_PRECISION))
		inv2 := (numTokens.Add(sdk.OneInt()).Mul(inv)).Add(
			ampTimeTotal.Sub(AMP_PRECISION).Mul(P_D).Quo(AMP_PRECISION),
		)
		inv = inv1.Add(inv2)

		if inv.GT(preInv) {
			if inv.Sub(preInv).LTE(sdk.OneInt()) {
				return inv, nil
			}
		} else if preInv.Sub(inv).LTE(sdk.OneInt()) {
			return inv, nil
		}
		return sdkmath.ZeroInt(), ErrInvalidInvariantConverge
	}

	return sdkmath.ZeroInt(), nil
}

func getTokenBalanceGivenInvariantAndAllOtherBalances(
	amp sdkmath.Int,
	inv sdkmath.Int,
	assets map[string]PoolAsset,
	tokenInDenom string,
) (sdkmath.Int, error) {
	//assets := p.GetAssetList()
	numTokens := sdkmath.NewInt(int64(len(assets)))
	ampTimeTotal := amp.Mul(numTokens)
	sum := sdkmath.NewInt(0)

	P_D := numTokens.Mul(assets[tokenInDenom].Token.Amount)

	for _, asset := range assets {
		P_D = P_D.Mul(asset.Token.Amount).Mul(numTokens).Quo(inv)
		sum = sum.Add(asset.Token.Amount)
	}

	sum = sum.Sub(assets[tokenInDenom].Token.Amount)

	inv2 := inv.Mul(inv)

	c := inv2.Quo(ampTimeTotal.Mul(P_D)).Mul(AMP_PRECISION).Mul(assets[tokenInDenom].Token.Amount)
	b := sum.Add(inv.Quo(ampTimeTotal).Mul(AMP_PRECISION))

	preTokenBalance := sdkmath.NewInt(0)

	tokenBalance := (inv2.Add(c)).Quo(inv.Add(b))

	for i := 0; i < 255; i++ {
		preTokenBalance = tokenBalance
		tokenBalance = tokenBalance.Mul(tokenBalance).Add(c).Quo((tokenBalance.Mul(sdkmath.NewInt(2)).Add(b).Sub(inv)))

		if tokenBalance.GT(preTokenBalance) {
			if tokenBalance.Sub(preTokenBalance).LTE(sdkmath.OneInt()) {
				return tokenBalance, nil
			}
		} else if preTokenBalance.Sub(tokenBalance).LTE(sdkmath.OneInt()) {
			return tokenBalance, nil
		}
	}
	return sdkmath.ZeroInt(), ErrInvalidInvariantConverge
}

// Helper functions
func MinusFees(amount sdkmath.Int, swapFee sdkmath.LegacyDec) sdk.Dec {
	amountDec := sdk.NewDecFromInt(amount)
	feeRate := swapFee.Quo(sdk.NewDec(10000))
	fees := amountDec.Mul(feeRate)
	amountMinusFees := amountDec.Sub(fees)
	return amountMinusFees
}
