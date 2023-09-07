package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) estimateShareInStablePool(coins sdk.Coins) (sdk.Coin, error) {
	_ = coins
	return sdk.Coin{}, nil
}

func (p *Pool) estimateSwapInStablePool(amountIn sdk.Coin, denomOut string) (sdk.Coin, error) {
	_ = amountIn
	_ = denomOut
	return sdk.Coin{}, nil
}

func (p *Pool) estimateWithdrawalsFromStablePool(share sdk.Coin) ([]sdk.Coin, error) {
	_ = share
	return []sdk.Coin{}, nil
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

func (p *Pool) calculateInvariantInStablePool(
	amp *sdkmath.Int,
) (*sdkmath.Int, error) {

	// assets
	assets := p.GetAssetList()

	// Initialize sum as zero
	sum := sdkmath.NewInt(0)

	// Number of tokens
	numTokens := sdkmath.NewInt(int64(len(assets)))

	for _, asset := range assets {
		sum = sum.Add(asset.Token.Amount)
	}

	if sum.IsZero() {
		return &sum, nil
	}

	_ = numTokens
	// preInv := sdkmath.NewInt(0)
	// inv := sum
	// ampTimesNpowN := math.Pow(to, y)
	return nil, nil
}
