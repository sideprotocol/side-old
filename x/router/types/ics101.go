package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gmmtypes "github.com/sideprotocol/side/x/gmm/types"
)

type PoolQuery struct {
	PoolID string `json:"pool_id"`
}

type WasmPoolAsset struct {
	Balance sdk.Coin    `json:"balance"`
	Decimal sdkmath.Int `json:"decimal"`
	Weight  sdkmath.Int `json:"weight"`
}

// InterchainPoolResponse mirrors the Rust struct in Go
type InterchainPoolResponse struct {
	ID                  string          `json:"id"`
	SourceCreator       string          `json:"source_creator"`
	SourceChainID       string          `json:"source_chain_id"`
	DestinationChainID  string          `json:"destination_chain_id"`
	DestinationCreator  string          `json:"destination_creator"`
	Assets              []WasmPoolAsset `json:"assets"`
	SwapFee             uint32          `json:"swap_fee"`
	Supply              sdk.Coin        `json:"supply"`
	CounterPartyPort    string          `json:"counter_party_port"`
	CounterPartyChannel string          `json:"counter_party_channel"`
}

func (iPoolRes *InterchainPoolResponse) ToGmmPool() gmmtypes.Pool {
	amp := sdkmath.NewInt(int64(0))
	assets := make(map[string]gmmtypes.PoolAsset)
	for _, asset := range iPoolRes.Assets {
		assets[asset.Balance.Denom] = gmmtypes.PoolAsset{
			Token:   asset.Balance,
			Weight:  &asset.Weight,
			Decimal: asset.Decimal,
		}
	}
	return gmmtypes.Pool{
		PoolId: iPoolRes.ID,
		Sender: iPoolRes.SourceCreator,
		PoolParams: gmmtypes.PoolParams{
			SwapFee: sdkmath.LegacyNewDecFromInt(sdk.NewInt(int64(iPoolRes.SwapFee))),
			Amp:     &amp,
		},
		Assets:      assets,
		TotalShares: iPoolRes.Supply,
	}
}
