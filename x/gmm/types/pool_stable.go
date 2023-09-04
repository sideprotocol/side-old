package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) estimateShareInStablePool(coins types.Coins) ([]types.Coin, error) {
	panic("did not implement still")
}

func (p *Pool) estimateSwapInStablePool(amountIn types.Coin, denomOut string) (types.Coin, error) {
	return types.Coin{}, nil
}

func (p *Pool) estimateWithdrawalsFromStablePool(share types.Coin) ([]types.Coin, error) {

	return []types.Coin{}, nil
}
