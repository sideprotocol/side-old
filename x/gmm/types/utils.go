package types

import (
	"crypto/sha256"
	"encoding/hex"
	fmt "fmt"
	"sort"
	"strings"
)

func GetPoolId(denoms []string) string {
	//generate poolId
	sort.Strings(denoms)
	poolIdHash := sha256.New()
	//salt := GenerateRandomString(chainID, 10)
	poolIdHash.Write([]byte(strings.Join(denoms, "")))
	poolId := "pool" + fmt.Sprintf("%v", hex.EncodeToString(poolIdHash.Sum(nil)))
	return poolId
}
