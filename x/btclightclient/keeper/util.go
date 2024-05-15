package keeper

import (
	"math/big"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/sideprotocol/side/x/btclightclient/types"
)

func HeaderConvert(header *types.BlockHeader) *wire.BlockHeader {
	prehash, _ := chainhash.NewHashFromStr(header.PreviousBlockHash)
	root, _ := chainhash.NewHashFromStr(header.MerkleRoot)
	n := new(big.Int)
	n.SetString(header.Bits, 16)
	return &wire.BlockHeader{
		Version:    int32(header.Version),
		PrevBlock:  *prehash,
		MerkleRoot: *root,
		Timestamp:  time.Unix(int64(header.Time), 0),
		Bits:       uint32(n.Uint64()),
		Nonce:      uint32(header.Nonce),
	}
}

func BitsToTarget(bits string) *big.Int {
	n := new(big.Int)
	n.SetString(bits, 16)
	return n
}

func BitsToTargetUint32(bits string) uint32 {
	return uint32(BitsToTarget(bits).Uint64())
}
