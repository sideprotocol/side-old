package types

import (
	"encoding/json"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VolumeData struct {
	PoolID      string
	AssetID     string
	BlockHeight int64
	Volume      sdkmath.Int // Volume data type depending on your requirements
}

type VolumeStack struct {
	Data []VolumeData
	Top  int
}

func (v *VolumeStack) Encode() ([]byte, error) {
	return json.Marshal(v.Data)
}
func (v *VolumeStack) Decode(data []byte) error {
	return json.Unmarshal(data, &v)
}

func (vs *VolumeStack) PushOrUpdate(ctx sdk.Context, poolID string, assetID string, volume sdkmath.Int) {
	blockHeight := ctx.BlockHeight()

	// Check for existing data at the same block height
	for i := range vs.Data {
		if vs.Data[i].BlockHeight == blockHeight && vs.Data[i].PoolID == poolID && vs.Data[i].AssetID == assetID {
			vs.Data[i].Volume = volume
			return
		}
	}

	// Add new data or replace the oldest if stack is full
	newData := VolumeData{
		PoolID:      poolID,
		AssetID:     assetID,
		BlockHeight: blockHeight,
		Volume:      volume,
	}

	if vs.Top >= 100 {
		// Stack is full, replace the oldest
		vs.Top = 0
	}
	vs.Data[vs.Top] = newData
	vs.Top++
}
