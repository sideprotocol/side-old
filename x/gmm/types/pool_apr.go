package types

import (
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PoolAPR struct {
	Fees      sdk.Coins
	CreatedAt int64
}

func (pa *PoolAPR) Decode(ctx sdk.Context, data []byte) error {
	if len(data) == 0 {
		pa = NewPoolAPR(ctx)
		return nil
	}
	return json.Unmarshal(data, &pa)
}

func (pa *PoolAPR) Encode() ([]byte, error) {
	return json.Marshal(pa)
}

func NewPoolAPR(ctx sdk.Context) *PoolAPR {
	return &PoolAPR{
		Fees:      sdk.NewCoins(),
		CreatedAt: ctx.BlockTime().Unix(),
	}
}

// Please multiply market price of every assets after getting to display as USD.
func (pa *PoolAPR) CalcAPR(ctx sdk.Context, tvl map[string]PoolAsset) sdk.Coins {
	oneYearAsSec := int64(60 * 60 * 24 * 365)
	var apr sdk.Coins
	for _, coin := range pa.Fees {
		interval := ctx.BlockTime().Unix() - pa.CreatedAt

		if _, found := tvl[coin.Denom]; found && !tvl[coin.Denom].Token.Amount.IsZero() {
			if interval <= 0 {
				apr = apr.Add(sdk.NewCoin(coin.Denom, sdk.NewInt(0)))
				continue
			}
			avg := coin.Amount.Mul(sdkmath.NewInt(oneYearAsSec)).Mul(sdk.NewInt(1e10)).Quo(sdkmath.NewInt(interval))
			avg = avg.Quo(tvl[coin.Denom].Token.Amount)
			apr = apr.Add(sdk.NewCoin(coin.Denom, avg))
		}
	}
	return apr
}
