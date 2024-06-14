package types

// this line is used by starport scaffolding # genesis/types/import

func DefaultBestBlockHeader() *BlockHeader {
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
