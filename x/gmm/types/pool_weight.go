package types

import (
	fmt "fmt"
	math "math"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) estimateShareInWeightPool(coins types.Coins) (types.Coin, error) {
	switch len(coins) {
	// Single Asset Deposit in balancer pool
	case 1:
		return p.estimateShareWithSingleLiquidityInWeightPool(coins[0])
	// Multi Asset Deposit in balancer pool
	case len(p.Assets):
		return p.estimateShareWithTalLiquidityInWeightPool(coins)
	}
	return types.Coin{}, ErrInvalidNumOfAssets
}

func (p *Pool) estimateShareWithSingleLiquidityInWeightPool(coin types.Coin) (types.Coin, error) {
	asset, err := p.findAssetByDenom(coin.Denom)
	if err != nil {
		return types.Coin{}, err
	}

	decToken := (types.NewDecCoinFromCoin(coin))
	decAsset := types.NewDecCoinFromCoin(*&asset.Token)
	weight := types.NewDecFromInt(asset.Weight).Quo(types.NewDec(100)) // divide by 100
	ratio := decToken.Amount.Quo(decAsset.Amount).Add(types.NewDec(1))
	exponent := (math.Pow(ratio.MustFloat64(), weight.MustFloat64()) - 1) * Multiplier
	factor, err := types.NewDecFromStr(fmt.Sprintf("%f", exponent/Multiplier))
	if err != nil {
		return types.Coin{}, err
	}
	issueAmount := p.TotalShares.Amount.Mul(factor.RoundInt()).Quo(types.NewInt(1e10))
	outputToken := types.Coin{
		Amount: issueAmount,
		Denom:  p.TotalShares.Denom,
	}
	return outputToken, nil
}

func (p *Pool) estimateShareWithTalLiquidityInWeightPool(coins types.Coins) (types.Coin, error) {
	share := sdk.NewInt(0)
	for _, coin := range coins {
		asset, err := p.findAssetByDenom(coin.Denom)
		if err != nil {
			return types.Coin{}, err
		}

		decToken := types.NewDecCoinFromCoin(coin)
		decAsset := types.NewDecCoinFromCoin(asset.Token)
		decSupply := types.NewDecCoinFromCoin(p.TotalShares)

		ratio := decToken.Amount.Quo(decAsset.Amount).Mul(types.NewDec(Multiplier))
		issueAmount := (decSupply.Amount.Mul(types.NewDecFromInt(asset.Weight)).Mul(ratio).Quo(types.NewDec(100)).Quo(types.NewDec(Multiplier)))

		share = share.Add(issueAmount.RoundInt())
	}

	return types.NewCoin(p.TotalShares.Denom, share), nil
}

// Withdraws 
func (p *Pool) estimateWithdrawalsFromWeightPool(share types.Coin) ([]types.Coin, error) {
	outs := []types.Coin{}
	if share.Amount.GT(p.TotalShares.Amount) {
		return nil, ErrOverflowShareAmount
	}
	for _, asset := range p.Assets {
		out := asset.Token.Amount.Mul(share.Amount).Quo(p.TotalShares.Amount)
		outs = append(outs, types.Coin{
			Denom:  asset.Token.Denom,
			Amount: out,
		})
	}
	return outs, nil
}


// Swap implements OutGivenIn
// Input how many coins you want to sell, output an amount you will receive
// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
func (p *Pool) estimateSwapInWeightPool(amountIn types.Coin, denomOut string) (types.Coin, error) {
	assetIn, err := p.findAssetByDenom(amountIn.Denom)
	if err != nil {
		return types.Coin{}, fmt.Errorf("left swap failed: could not find asset in by denom")
	}

	assetOut, err := p.findAssetByDenom(denomOut)
	if err != nil {
		return types.Coin{}, fmt.Errorf("left swap failed: could not find asset out by denom")
	}

	balanceOut := types.NewDecFromBigInt(assetOut.Token.Amount.BigInt())
	balanceIn := types.NewDecFromBigInt(assetIn.Token.Amount.BigInt())
	weightIn := types.NewDecFromInt(assetIn.Weight).Quo(types.NewDec(100))
	weightOut := types.NewDecFromInt(assetIn.Weight).Quo(types.NewDec(100))
	amount := p.TakeFees(amountIn.Amount)

	// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
	balanceInPlusAmount := balanceIn.Add(amount)
	ratio := balanceIn.Quo(balanceInPlusAmount)
	oneMinusRatio := types.NewDec(1).Sub(ratio)

	power := weightIn.Quo(weightOut)
	factor := math.Pow(oneMinusRatio.MustFloat64(), power.MustFloat64()) * Multiplier
	finalFactor := factor / 1e8

	amountOut := balanceOut.Mul(types.MustNewDecFromStr(fmt.Sprintf("%f", finalFactor))).Quo(types.NewDec(1e10))
	return types.Coin{
		Amount: amountOut.RoundInt(),
		Denom:  denomOut,
	}, nil
}