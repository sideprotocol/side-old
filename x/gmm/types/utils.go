package types

import (
	"crypto/sha256"
	"encoding/hex"
	fmt "fmt"
	"sort"
	"strings"
)

func GetPoolID(denoms []string) string {
	// Generate poolID
	sort.Strings(denoms)
	poolIDHash := sha256.New()

	poolIDHash.Write([]byte(strings.Join(denoms, "")))
	poolID := "pool" + fmt.Sprintf("%v", hex.EncodeToString(poolIDHash.Sum(nil)))
	return poolID
}
