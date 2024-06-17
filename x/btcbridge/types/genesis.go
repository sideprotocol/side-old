package types

import (
	"github.com/btcsuite/btcd/chaincfg"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// this line is used by starport scaffolding # genesis/types/import

func DefaultBestBlockHeader() *BlockHeader {
	config := sdk.GetConfig().GetBtcChainCfg()
	switch config.Name {
	case chaincfg.MainNetParams.Name:
		return DefaultMainNetBestBlockHeader()
	case chaincfg.SigNetParams.Name:
		return DefaultSignetBestBlockHeader()
	}
	return DefaultTestnetBestBlockHeader()
}

func DefaultSignetBestBlockHeader() *BlockHeader {
	// testnet3 block 2815023
	return &BlockHeader{
		Version:           536870912,
		Hash:              "000000a7e90dc93d4671819d34bce273f65b7f911f335d829b6a6e181cda01c8",
		Height:            3236,
		PreviousBlockHash: "0000011e57524f2abb7054ac10833926470caefda7139059a255d25aaa306343",
		MerkleRoot:        "c7edd2424b9fe2776c7cacb766efa4a308c68cb45838db185dfcf817613d92a4",
		Time:              1718665272,
		Bits:              "1e0377ae",
		Nonce:             6997531,
		Ntx:               1,
	}
}

func DefaultTestnetBestBlockHeader() *BlockHeader {
	// testnet3 block 2815023
	return &BlockHeader{
		Version:           667459584,
		Hash:              "0000000000000009fb68da72e8994f014fafb455c72978233b94580b12af778c",
		Height:            2815023,
		PreviousBlockHash: "0000000000000004a29c20eb32532718de8072665620edb4c657b22b4d463967",
		MerkleRoot:        "9e219423eadce80e882cdff04b3026c9bbc994fd08a774f34a705ca3e710a332",
		Time:              1715566066,
		Bits:              "191881b8",
		Nonce:             3913166971,
		Ntx:               6236,
	}
}

func DefaultMainNetBestBlockHeader() *BlockHeader {
	// testnet3 block 2815023
	return &BlockHeader{
		Version:           667459584,
		Hash:              "0000000000000009fb68da72e8994f014fafb455c72978233b94580b12af778c",
		Height:            2815023,
		PreviousBlockHash: "0000000000000004a29c20eb32532718de8072665620edb4c657b22b4d463967",
		MerkleRoot:        "9e219423eadce80e882cdff04b3026c9bbc994fd08a774f34a705ca3e710a332",
		Time:              1715566066,
		Bits:              "191881b8",
		Nonce:             3913166971,
		Ntx:               6236,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:          DefaultParams(),
		BestBlockHeader: DefaultBestBlockHeader(),
		BlockHeaders:    []*BlockHeader{},
		Utxos:           []*UTXO{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	// need to be improved by checking the block headers & best block header
	if gs.BestBlockHeader == nil || gs.BestBlockHeader.Hash == "" || gs.BestBlockHeader.PreviousBlockHash == "" || gs.BestBlockHeader.MerkleRoot == "" {
		return ErrInvalidHeader
	}
	return gs.Params.Validate()
}
