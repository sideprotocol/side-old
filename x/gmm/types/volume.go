package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxObservations = 100
const SecondsIn24Hours = int64(86400)

type VolumeData struct {
	PoolID    string
	BlockTime int64
	Volume    sdk.Coins
}

type VolumeStack struct {
	Data        []VolumeData
	TotalVolume sdk.Coins
	Top         int
}

func NewVolumeStack() *VolumeStack {
	return &VolumeStack{
		Data:        make([]VolumeData, MaxObservations),
		TotalVolume: sdk.NewCoins(),
		Top:         0,
	}
}

func (v *VolumeStack) Encode() ([]byte, error) {
	return json.Marshal(v)
}

func (v *VolumeStack) Decode(data []byte) error {
	return json.Unmarshal(data, v)
}

// Ob adds a new observation or updates an existing one.
func (vs *VolumeStack) Observe(ctx sdk.Context, poolID string, newSwap sdk.Coins) {
	blockTime := ctx.BlockTime().Unix()

	// Update total volume
	vs.TotalVolume = vs.TotalVolume.Add(newSwap...)

	// Check and update if a recent entry (24 hours) exists
	updated := false
	for i := range vs.Data {
		if vs.Data[i].PoolID == poolID && blockTime-vs.Data[i].BlockTime < SecondsIn24Hours {
			vs.Data[i].Volume = vs.Data[i].Volume.Add(newSwap...)
			updated = true
			break
		}
	}

	if !updated {
		// Add new data or replace the oldest if stack is full
		newData := VolumeData{
			PoolID:    poolID,
			BlockTime: blockTime,
			Volume:    newSwap,
		}
		vs.Data[vs.Top] = newData
		vs.Top = (vs.Top + 1) % MaxObservations
	}
}

// Calculate24HourVolume calculates the volume for the last 24 hours.
func (vs *VolumeStack) Calculate24HourVolume(ctx sdk.Context, poolID string) sdk.Coins {
	currentTime := ctx.BlockTime().Unix()
	var volume24h sdk.Coins

	for _, data := range vs.Data {
		if data.PoolID == poolID && currentTime-data.BlockTime < SecondsIn24Hours {
			volume24h = volume24h.Add(data.Volume...)
		}
	}
	return volume24h
}

// GetTotalVolume returns the total volume for the pool.
func (vs *VolumeStack) GetTotalVolume() sdk.Coins {
	return vs.TotalVolume
}
