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
	balanceRatiosWithFee := make(map[string]sdkmath.LegacyDec, len(coins))
	// The weighted sum of token balance ratios without fee
	invariantRatioWithFees := sdkmath.LegacyZeroDec()

	for _, asset := range coins {
		currentWeight := sdkmath.LegacyDec(asset.Amount).Quo(sdkmath.LegacyNewDecFromInt(sum))

		assetInPool, _, exist := p.GetAssetByDenom(asset.Denom)
		if !exist {
			continue
		}
		balanceRatiosWithFee[asset.Denom] = sdkmath.LegacyDec(asset.Amount.Add(assetInPool.Token.Amount)).Quo(
			sdkmath.LegacyNewDecFromInt(assetInPool.Token.Amount),
		)

		invariantRatioWithFees = invariantRatioWithFees.Add(
			balanceRatiosWithFee[asset.Denom].Mul(currentWeight),
		)
	}

	// Second loop calculates new amounts in, taking into account the fee on the percentage excess

	newBalances := sdk.NewCoins()
	for _, amountIn := range coins {
		asset, _, exist := p.GetAssetByDenom(amountIn.Denom)
		if !exist {
			continue
		}
		var amountInWithoutFee sdkmath.Int
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
	currentInvariant := calculateInvariantInStablePool(*p.PoolParams.Amp, p.GetLiquidity())

	newInvariant := calculateInvariantInStablePool(*p.PoolParams.Amp, newBalances)

	invariantRatio := sdkmath.LegacyNewDecFromInt(newInvariant).Quo(
		sdkmath.LegacyNewDecFromInt(currentInvariant),
	)

	// If the invariant didn't increase for any reason, we simply don't mint BPT
	if invariantRatio.GT(sdkmath.LegacyZeroDec()) {
		share := p.TotalShares.Amount.Mul(sdkmath.Int(invariantRatio.Sub(sdkmath.LegacyOneDec())))
		return sdk.NewCoin(p.PoolId, share), nil

	}
	return sdk.NewCoin(p.PoolId, sdkmath.NewInt(0)), nil
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

	inv := calculateInvariantInStablePool(*p.PoolParams.Amp, p.GetLiquidity())

	assets := p.Assets
	asset, index, exist := p.GetAssetByDenom(tokenIn.Denom)
	if !exist {
		return sdk.Coin{}, ErrNotFoundAssetInPool
	}
	balance := asset.Token.Amount.Add(tokenInDec.RoundInt())
	assets[index] = PoolAsset{
		Token:   sdk.NewCoin(tokenIn.Denom, balance),
		Weight:  assets[index].Weight,
		Decimal: assets[index].Decimal,
	}

	finalBalanceOut, err := getTokenBalanceGivenInvariantAndAllOtherBalances(
		*p.PoolParams.Amp, inv, assets, tokenIn.Amount,
	)
	out := p.Assets[index].Token.Amount.Sub(finalBalanceOut).Sub(sdkmath.OneInt())
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

var AmpPrecision = sdkmath.NewInt(1000)

func calculateInvariantInStablePool(
	amp sdkmath.Int,
	assets []sdk.Coin,
) sdkmath.Int {
	sum := sdkmath.NewInt(0)

	// Number of tokens
	numTokens := sdkmath.NewInt(int64(len(assets)))

	for _, asset := range assets {
		sum = sum.Add(asset.Amount)
	}

	if sum.IsZero() {
		return sum
	}

	inv := sum
	ampTimeTotal := amp.Mul(numTokens)

	// nolint:staticcheck
	for i := 0; i < 255; i++ {
		PD := numTokens.Mul(assets[0].Amount)
		for _, asset := range assets[1:] {
			PD = PD.Mul(asset.Amount).Mul(numTokens).Quo(inv)
		}

		preInv := inv
		numerator1 := numTokens.Mul(inv).Mul(inv)
		numerator2 := ampTimeTotal.Mul(sum).Mul(PD).Quo(AmpPrecision)
		denominator := numTokens.Add(sdk.OneInt()).Mul(inv).Add(ampTimeTotal.Sub(AmpPrecision).Mul(PD).Quo(AmpPrecision))
		inv = numerator1.Add(numerator2).Quo(denominator)

		if inv.GT(preInv) {
			if inv.Sub(preInv).LTE(sdkmath.NewInt(1e18)) {
				break
			}
		} else if preInv.Sub(inv).LTE(sdkmath.NewInt(1e18)) {
			break
		}
	}
	return inv
}

func getTokenBalanceGivenInvariantAndAllOtherBalances(
	amp sdkmath.Int,
	inv sdkmath.Int,
	assets []PoolAsset,
	tokenInAmount sdkmath.Int,
) (sdkmath.Int, error) {
	numTokens := sdkmath.NewInt(int64(len(assets)))
	ampTimeTotal := amp.Mul(numTokens)
	sum := sdkmath.NewInt(0)

	PD := numTokens.Mul(tokenInAmount)

	for _, asset := range assets {
		PD = PD.Mul(asset.Token.Amount).Mul(numTokens).Quo(inv)
		sum = sum.Add(asset.Token.Amount)
	}

	sum = sum.Sub(tokenInAmount)

	inv2 := inv.Mul(inv)

	c := inv2.Quo(ampTimeTotal.Mul(PD)).Mul(AmpPrecision).Mul(tokenInAmount)
	b := sum.Add(inv.Quo(ampTimeTotal).Mul(AmpPrecision))

	tokenBalance := (inv2.Add(c)).Quo(inv.Add(b))

	for i := 0; i < 255; i++ {
		preTokenBalance := tokenBalance
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
