package types

import (
	fmt "fmt"
	"math/big"

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
	weight := sdk.NewDecFromInt(*asset.Weight).Quo(sdk.NewDec(100)) // divide by 100
	ratio := decToken.Amount.Quo(decAsset.Amount).Add(sdk.NewDec(1))
	precision := big.NewInt(1) //sdk.MustNewDecFromStr("0.00000001")
	_ = weight
	_ = ratio
	_ = precision
	factor := sdk.NewInt(1)
	//factor := (ApproximatePow(ratio.BigInt(), weight.BigInt(), precision).Sub(sdk.OneDec()))
	issueAmount := p.TotalShares.Amount.Mul(factor).Quo(sdk.NewInt(1e10))
	outputToken := sdk.Coin{
		Amount: issueAmount,
		Denom:  p.TotalShares.Denom,
	}
	return outputToken, nil
}

func (p *Pool) estimateShareWithTalLiquidityInWeightPool(coins sdk.Coins) (sdk.Coin, error) {
	share := sdk.NewInt(0)
	for _, coin := range coins {
		asset, err := p.findAssetByDenom(coin.Denom)
		if err != nil {
			return sdk.Coin{}, err
		}

		decToken := sdk.NewDecCoinFromCoin(coin)
		decAsset := sdk.NewDecCoinFromCoin(asset.Token)
		decSupply := sdk.NewDecCoinFromCoin(p.TotalShares)

		ratio := decToken.Amount.Quo(decAsset.Amount).Mul(sdk.NewDec(Multiplier))
		issueAmount := (decSupply.Amount.Mul(sdk.NewDecFromInt(*asset.Weight)).Mul(ratio).Quo(sdk.NewDec(100)).Quo(sdk.NewDec(Multiplier)))

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

	balanceOut := sdk.NewDecFromBigInt(assetOut.Token.Amount.BigInt())
	balanceIn := sdk.NewDecFromBigInt(assetIn.Token.Amount.BigInt())
	weightIn := sdk.NewDecFromInt(*assetIn.Weight).Quo(sdk.NewDec(100))
	weightOut := sdk.NewDecFromInt(*assetIn.Weight).Quo(sdk.NewDec(100))
	amount := p.TakeFees(amountIn.Amount)

	// Ao = Bo * ((1 - Bi / (Bi + Ai)) ** Wi/Wo)
	balanceInPlusAmount := balanceIn.Add(amount)
	ratio := balanceIn.Quo(balanceInPlusAmount)
	oneMinusRatio := sdk.NewDec(1).Sub(ratio)
	power := weightIn.Quo(weightOut)
	precision := "0.00000001"                                                        //sdk.MustNewDecFromStr("0.00000001")
	factor, err := ApproximatePow(oneMinusRatio.String(), power.String(), precision) // 100 iterations for example
	if err != nil {
		return sdk.Coin{}, err
	}
	amountOut := balanceOut.Mul(sdk.MustNewDecFromStr(factor.String()))
	return sdk.Coin{
		Amount: amountOut.RoundInt(),
		Denom:  denomOut,
	}, nil
}
