package types

import (
	fmt "fmt"
	math "math"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) estimateShareInWeightPool(coins sdk.Coins) (sdk.Coin, error) {
	switch len(coins) {
	// Single Asset Deposit in balancer pool
	case 1:
		return p.estimateShareWithSingleLiquidityInWeightPool(coins[0])
	// Multi Asset Deposit in balancer pool
	case len(p.Assets):
		return p.estimateShareWithTalLiquidityInWeightPool(coins)
	}
	return sdk.Coin{}, ErrInvalidNumOfAssets
}

func (p *Pool) estimateShareWithSingleLiquidityInWeightPool(coin sdk.Coin) (sdk.Coin, error) {
	asset, err := p.findAssetByDenom(coin.Denom)
	if err != nil {
		return sdk.Coin{}, err
	}

	decToken := (sdk.NewDecCoinFromCoin(coin))
	decAsset := sdk.NewDecCoinFromCoin(asset.Token)
	weight := sdkmath.LegacyNewDecFromInt(*asset.Weight).Quo(sdkmath.LegacyNewDec(100)) // divide by 100
	ratio := decToken.Amount.Quo(decAsset.Amount).Add(sdkmath.LegacyNewDec(1))
	exponent := (math.Pow(ratio.MustFloat64(), weight.MustFloat64()) - 1) * Multiplier
	factor, err := sdkmath.LegacyNewDecFromStr(fmt.Sprintf("%f", exponent/Multiplier))
	if err != nil {
		return sdk.Coin{}, err
	}
	issueAmount := p.TotalShares.Amount.Mul(factor.RoundInt()).Quo(sdkmath.NewInt(1e10))
	outputToken := sdk.Coin{
		Amount: issueAmount,
		Denom:  p.TotalShares.Denom,
	}
	return outputToken, nil
}

func (p *Pool) estimateShareWithTalLiquidityInWeightPool(coins sdk.Coins) (sdk.Coin, error) {
	share := sdkmath.NewInt(0)
	for _, coin := range coins {
		asset, err := p.findAssetByDenom(coin.Denom)
		if err != nil {
			return sdk.Coin{}, err
		}

		decToken := sdk.NewDecCoinFromCoin(coin)
		decAsset := sdk.NewDecCoinFromCoin(asset.Token)
		decSupply := sdk.NewDecCoinFromCoin(p.TotalShares)

		ratio := decToken.Amount.Quo(decAsset.Amount).Mul(sdkmath.LegacyNewDec(Multiplier))
		issueAmount := (decSupply.Amount.Mul(sdkmath.LegacyNewDecFromInt(*asset.Weight)).Mul(ratio).Quo(sdkmath.LegacyNewDec(100)).Quo(sdkmath.LegacyNewDec(Multiplier)))

		share = share.Add(issueAmount.RoundInt())
	}

	return sdk.NewCoin(p.TotalShares.Denom, share), nil
}

// Withdraws
func (p *Pool) estimateWithdrawalsFromWeightPool(share sdk.Coin) ([]sdk.Coin, error) {
	outs := []sdk.Coin{}
	if share.Amount.GT(p.TotalShares.Amount) {
		return nil, ErrOverflowShareAmount
	}
	for _, asset := range p.Assets {
		out := asset.Token.Amount.Mul(share.Amount).Quo(p.TotalShares.Amount)
		outs = append(outs, sdk.Coin{
			Denom:  asset.Token.Denom,
			Amount: out,
		})
	}
	return outs, nil
}

// Swap implements OutGivenIn
// Input how many coins you want to sell, output an amount you will receive
// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
func (p *Pool) estimateSwapInWeightPool(amountIn sdk.Coin, denomOut string) (sdk.Coin, error) {
	assetIn, err := p.findAssetByDenom(amountIn.Denom)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("left swap failed: could not find asset in by denom")
	}

	assetOut, err := p.findAssetByDenom(denomOut)
	if err != nil {
		return sdk.Coin{}, fmt.Errorf("left swap failed: could not find asset out by denom")
	}

	balanceOut := sdkmath.LegacyNewDecFromBigInt(assetOut.Token.Amount.BigInt())
	balanceIn := sdkmath.LegacyNewDecFromBigInt(assetIn.Token.Amount.BigInt())
	weightIn := sdkmath.LegacyNewDecFromInt(*assetIn.Weight).Quo(sdkmath.LegacyNewDec(100))
	weightOut := sdkmath.LegacyNewDecFromInt(*assetIn.Weight).Quo(sdkmath.LegacyNewDec(100))
	amount := p.TakeFees(amountIn.Amount)

	// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
	balanceInPlusAmount := balanceIn.Add(amount)
	ratio := balanceIn.Quo(balanceInPlusAmount)
	oneMinusRatio := sdkmath.LegacyNewDec(1).Sub(ratio)

	power := weightIn.Quo(weightOut)
	factor := math.Pow(oneMinusRatio.MustFloat64(), power.MustFloat64()) * Multiplier
	finalFactor := factor / 1e8

	amountOut := balanceOut.Mul(sdkmath.LegacyMustNewDecFromStr(fmt.Sprintf("%f", finalFactor))).Quo(sdkmath.LegacyNewDec(1e10))
	return sdk.Coin{
		Amount: amountOut.RoundInt(),
		Denom:  denomOut,
	}, nil
}
