package types

import (
	"crypto/sha256"
	"encoding/hex"
	fmt "fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Alice = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	Bob   = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	Carol = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
)

func GetPoolID(denoms []string) string {
	// Generate poolID
	sort.Strings(denoms)
	poolIDHash := sha256.New()

	poolIDHash.Write([]byte(strings.Join(denoms, "")))
	poolID := "pool" + fmt.Sprintf("%v", hex.EncodeToString(poolIDHash.Sum(nil)))
	return poolID
}

func GetEventAttrOfAsset(assets []sdk.Coin) []sdk.Attribute {
	var attr []sdk.Attribute
	for index, asset := range assets {
		attr = append(attr, sdk.NewAttribute(
			fmt.Sprintf("%d", index),
			asset.String(),
		))
	}
	return attr
}
