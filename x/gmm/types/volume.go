package types

import (
	"encoding/json"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MaxObservations  = 100
	SecondsIn24Hours = int64(86400)
)

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
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

// Ob adds a new observation or updates an existing one.
func (v *VolumeStack) Observe(ctx sdk.Context, poolID string, newSwap sdk.Coins) {
	blockTime := ctx.BlockTime().Unix()

	// Update total volume
	v.TotalVolume = v.TotalVolume.Add(newSwap...)

	// Add new data
	newData := VolumeData{
		PoolID:    poolID,
		BlockTime: blockTime,
		Volume:    newSwap,
	}
	v.Data = append(v.Data, newData)

	// Sort data based on BlockTime
	sort.Slice(v.Data, func(i, j int) bool {
		return v.Data[i].BlockTime < v.Data[j].BlockTime
	})

	// Purge data older than 24 hours
	currentTime := ctx.BlockTime().Unix()
	cutoffTime := currentTime - SecondsIn24Hours
	v.Data = filter(v.Data, func(vd VolumeData) bool {
		return vd.BlockTime >= cutoffTime
	})
}

func filter(vs []VolumeData, f func(VolumeData) bool) []VolumeData {
	vsf := make([]VolumeData, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Calculate24HourVolume calculates the volume for the last 24 hours.
func (v *VolumeStack) Calculate24HourVolume(ctx sdk.Context, poolID string) sdk.Coins {
	currentTime := ctx.BlockTime().Unix()
	var volume24h sdk.Coins

	for _, data := range v.Data {
		if data.PoolID == poolID && currentTime-data.BlockTime < SecondsIn24Hours {
			volume24h = volume24h.Add(data.Volume...)
		}
	}
	return volume24h
}

// GetTotalVolume returns the total volume for the pool.
func (v *VolumeStack) GetTotalVolume() sdk.Coins {
	return v.TotalVolume
}
