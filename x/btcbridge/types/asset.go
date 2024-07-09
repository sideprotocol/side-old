package types

import (
	"fmt"
	"strings"
)

// AssetTypeFromDenom returns the asset type according to the denom
func AssetTypeFromDenom(denom string, p Params) AssetType {
	if denom == p.BtcVoucherDenom {
		return AssetType_ASSET_TYPE_BTC
	}

	if strings.HasPrefix(denom, fmt.Sprintf("%s/", RunesProtocolName)) {
		return AssetType_ASSET_TYPE_RUNE
	}

	return AssetType_ASSET_TYPE_UNSPECIFIED
}
