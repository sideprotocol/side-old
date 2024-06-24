package types

import (
	sdkmath "cosmossdk.io/math"
)

func (p *PoolI) ToPool() Pool {
	assets := []PoolAsset{} // make(map[PoolAsset])
	for _, asset := range p.Assets {
		weight := sdkmath.NewIntFromUint64(uint64(asset.Weight))
		assets = append(assets, PoolAsset{
			Token:   *asset.Balance,
			Weight:  &weight,
			Decimal: sdkmath.NewIntFromUint64(uint64(asset.Decimal)),
		})
	}
	return Pool{
		PoolId: p.Id,
		Sender: p.SourceCreator,
		PoolParams: PoolParams{
			Type:      p.PoolType,
			SwapFee:   sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(int64(p.SwapFee))),
			ExitFee:   sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(int64(p.SwapFee))),
			UseOracle: false,
			Amp:       p.Amp,
		},
		Assets:      assets,
		TotalShares: *p.Supply,
	}
}
