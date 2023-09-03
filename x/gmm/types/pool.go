package types

import (
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) GetAssetDenoms() []string {
	denoms := []string{}
	for _, asset := range p.Assets {
		denoms = append(denoms, asset.Token.Denom)
	}
	return denoms
}

func GetEscrowAddress(poolId string) sdk.AccAddress {
	// a slash is used to create domain separation between port and channel identifiers to
	// prevent address collisions between escrow addresses created for different channels

	// ADR 028 AddressHash construction
	preImage := []byte(Version)
	preImage = append(preImage, 0)
	preImage = append(preImage, poolId...)
	hash := sha256.Sum256(preImage)
	return hash[:20]
}
