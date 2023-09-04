package types

import (
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) GetAssetDenoms() []string {
	denoms := []string{}
	for _, asset := range p.Assets {
		denoms = append(denoms, asset.Token.Denom)
	}
	return denoms
}

func GetEscrowAddress(poolId string) sdk.AccAddress {
	// a slash is used to create domain separation between port and channel identifiers to
	// prevent address collisions between escrow addresses created for different channels

	// ADR 028 AddressHash construction
	preImage := []byte(Version)
	preImage = append(preImage, 0)
	preImage = append(preImage, poolId...)
	hash := sha256.Sum256(preImage)
	return hash[:20]
}

// EstimateShare estimate share amount when user deposit
func (p *Pool) EstimateShare(coins types.Coins) (types.Coin, error) {
	switch p.PoolParams.Type {
	case PoolType_WEIGHT:
		return p.estimateShareInWeightPool(coins)
	case PoolType_STABLE:
		return p.estimateShareInWeightPool(coins)
	}
	return types.Coin{}, ErrInvalidPoolType
}

func (p *Pool) EstimateSwap(amountIn types.Coin, denomOut string) (types.Coin, error) {
	switch p.PoolParams.Type {
	case PoolType_WEIGHT:
		return p.estimateSwapInWeightPool(amountIn, denomOut)
	case PoolType_STABLE:
		return p.estimateSwapInStablePool(amountIn, denomOut)
	}
	return types.Coin{}, ErrInvalidPoolType
}

// Withdraw tokens from pool
func (p *Pool) EstimateWithdrawals(share types.Coin) ([]types.Coin, error) {
	switch p.PoolParams.Type {
	case PoolType_WEIGHT:
		return p.estimateWithdrawalsFromWeightPool(share)
	case PoolType_STABLE:
		return p.estimateWithdrawalsFromStablePool(share)
	}
	return []types.Coin{}, ErrInvalidPoolType
}

// Helper functions
func (p *Pool) TakeFees(amount types.Int) types.Dec {
	amountDec := types.NewDecFromInt(amount)
	feeRate := p.PoolParams.SwapFee.Quo(types.NewDec(10000))
	fees := amountDec.Mul(feeRate)
	amountMinusFees := amountDec.Sub(fees)
	return amountMinusFees
}

// IncreaseShare add xx amount share to total share amount in pool
func (p *Pool) IncreaseShare(amt sdk.Int) {
	p.TotalShares.Amount = p.TotalShares.Amount.Add(amt)
}

// DecreaseShare subtract xx amount share to total share amount in pool
func (p *Pool) DecreaseShare(amt sdk.Int) {
	p.TotalShares.Amount = p.TotalShares.Amount.Sub(amt)
}

// IncreaseLiquidity adds xx amount liquidity to assets in pool
func (p *Pool) IncreaseLiquidity(coins []types.Coin) error {
	for _, coin := range coins {
		asset, exists := p.Assets[coin.Denom]
		if !exists {
			return ErrNotFoundAssetInPool
		}
		// Add liquidity logic here
		asset.Token.Amount = asset.Token.Amount.Add(coin.Amount)
		p.Assets[coin.Denom] = asset
	}
	// Update TotalShares or other fields if necessary
	return nil
}

// DecreaseLiquidity subtracts xx amount liquidity from assets in pool
func (p *Pool) DecreaseLiquidity(coins []types.Coin) error {
	for _, coin := range coins {
		asset, exists := p.Assets[coin.Denom]
		if !exists {
			return ErrNotFoundAssetInPool
		}
		// Add liquidity logic here
		asset.Token.Amount = asset.Token.Amount.Sub(coin.Amount)
		p.Assets[coin.Denom] = asset
	}
	// Update TotalShares or other fields if necessary
	return nil
}

// findAssetByDenom finds pool asset by denom
func (p *Pool) findAssetByDenom(denom string) (PoolAsset, error) {
	for _, asset := range p.Assets {
		if asset.Token.Denom == denom {
			return asset, nil
		}
	}
	return PoolAsset{}, ErrNotFoundAssetInPool
}
