package types

import (
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
