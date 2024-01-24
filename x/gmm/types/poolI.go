package types

import (
	sdkmath "cosmossdk.io/math"
)

func (p *PoolI) ToPool() Pool {
	amp := sdkmath.NewInt(0)
	assets := make(map[string]PoolAsset)
	for _, asset := range p.Assets {
		weight := sdkmath.NewIntFromUint64(uint64(asset.Weight))
		assets[asset.Balance.Denom] = PoolAsset{
			Token:   *asset.Balance,
			Weight:  &weight,
			Decimal: sdkmath.NewIntFromUint64(uint64(asset.Decimal)),
		}
	}
	return Pool{
		PoolId: p.Id,
		Sender: p.SourceCreator,
		PoolParams: PoolParams{
			Type:      PoolType_WEIGHT,
			SwapFee:   sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(int64(p.SwapFee))),
			ExitFee:   sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(int64(p.SwapFee))),
			UseOracle: true,
			Amp:       &amp,
		},
		Assets:      assets,
		TotalShares: *p.Supply,
	}
}
