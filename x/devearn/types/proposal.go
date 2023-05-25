package types

import (
	"errors"
	"fmt"
	"sidechain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// constants
const (
	ProposalTypeRegisterDevEarnInfo      string = "RegisterDevEarnInfo"
	ProposalTypeCancelDevEarnInfo        string = "CancelDevEarnInfo"
	ProposalTypeAddAsssetToWhitelist     string = "AddAssetToWhitelist"
	ProposalTypeRemoveAssetFromWhitelist string = "RemoveAssetFromWhitelist"
)

// Implements Proposal Interface
var (
	_ govv1beta1.Content = &RegisterDevEarnInfoProposal{}
	_ govv1beta1.Content = &CancelDevEarnInfoProposal{}
	_ govv1beta1.Content = &AddAssetToWhitelistProposal{}
	_ govv1beta1.Content = &RemoveAssetFromWhitelistProposal{}
)

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeRegisterDevEarnInfo)
	govv1beta1.RegisterProposalType(ProposalTypeCancelDevEarnInfo)
	govv1beta1.RegisterProposalType(ProposalTypeAddAsssetToWhitelist)
	govv1beta1.RegisterProposalType(ProposalTypeRemoveAssetFromWhitelist)
	govv1beta1.ModuleCdc.Amino.RegisterConcrete(&RegisterDevEarnInfoProposal{}, "devearn/RegisterDevEarnInfoProposal", nil)
	govv1beta1.ModuleCdc.Amino.RegisterConcrete(&CancelDevEarnInfoProposal{}, "devearn/CancelDevEarnInfoProposal", nil)
	govv1beta1.ModuleCdc.Amino.RegisterConcrete(&AddAssetToWhitelistProposal{}, "devearn/AddAssetToWhitelistProposal", nil)
	govv1beta1.ModuleCdc.Amino.RegisterConcrete(&RemoveAssetFromWhitelistProposal{}, "devearn/RemoveAssetFromWhitelistProposal", nil)
}

// NewRegisterDevEarnInfoProposal returns new instance of RegisterIncentiveProposal
func NewRegisterDevEarnInfoProposal(
	title, description, contract, ownerAddr string,
	epochs uint32,
) govv1beta1.Content {
	return &RegisterDevEarnInfoProposal{
		Title:        title,
		Description:  description,
		Contract:     contract,
		OwnerAddress: ownerAddr,
		Epochs:       epochs,
	}
}

// ProposalRoute returns router key for this proposal
func (*RegisterDevEarnInfoProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RegisterDevEarnInfoProposal) ProposalType() string {
	return ProposalTypeRegisterDevEarnInfo
}

// ValidateBasic performs a stateless check of the proposal fields
func (rip *RegisterDevEarnInfoProposal) ValidateBasic() error {
	if err := types.ValidateAddress(rip.Contract); err != nil {
		return err
	}

	if err := validateEpochs(rip.Epochs); err != nil {
		return err
	}

	return govv1beta1.ValidateAbstract(rip)
}

// validateAllocations checks if each allocation has
// - a valid denom
// - a valid amount representing the percentage of allocation
func validateAllocations(allocations sdk.DecCoins) error {
	if allocations.Empty() {
		return errors.New("incentive allocations cannot be empty")
	}

	for _, al := range allocations {
		if err := validateAmount(al.Amount); err != nil {
			return err
		}
	}

	return allocations.Validate()
}

func validateAmount(amount sdk.Dec) error {
	if amount.GT(sdk.OneDec()) || amount.LTE(sdk.ZeroDec()) {
		return fmt.Errorf("invalid amount for allocation: %s", amount)
	}
	return nil
}

func validateEpochs(epochs uint32) error {
	if epochs == 0 {
		return fmt.Errorf("epochs value (%d) cannot be 0", epochs)
	}
	return nil
}

// NewCancelDevEarnProposal returns new instance of RegisterIncentiveProposal
func NewCancelDevEarnProposal(
	title, description, contract string,
) govv1beta1.Content {
	return &CancelDevEarnInfoProposal{
		Title:       title,
		Description: description,
		Contract:    contract,
	}
}

// ProposalRoute returns router key for this proposal
func (*CancelDevEarnInfoProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*CancelDevEarnInfoProposal) ProposalType() string {
	return ProposalTypeCancelDevEarnInfo
}

// ValidateBasic performs a stateless check of the proposal fields
func (rip *CancelDevEarnInfoProposal) ValidateBasic() error {
	if err := types.ValidateAddress(rip.Contract); err != nil {
		return err
	}

	return govv1beta1.ValidateAbstract(rip)
}

// NewAddAssetToWhitelistProposal add asset to whitelist
func NewAddAssetToWhitelistProposal(
	title, description, denom string,
) govv1beta1.Content {
	return &AddAssetToWhitelistProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}

// ProposalRoute returns router key for this proposal
func (*AddAssetToWhitelistProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*AddAssetToWhitelistProposal) ProposalType() string {
	return ProposalTypeAddAsssetToWhitelist
}

// ValidateBasic performs a stateless check of the proposal fields
func (rip *AddAssetToWhitelistProposal) ValidateBasic() error {
	return govv1beta1.ValidateAbstract(rip)
}

// RemoveAssetFromWhitelistProposal returns new instance of RegisterIncentiveProposal
func NewRemoveAssetFromWhitelistProposal(
	title, description, denom string,
) govv1beta1.Content {
	return &RemoveAssetFromWhitelistProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}

// ProposalRoute returns router key for this proposal
func (*RemoveAssetFromWhitelistProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns proposal type for this proposal
func (*RemoveAssetFromWhitelistProposal) ProposalType() string {
	return ProposalTypeRemoveAssetFromWhitelist
}

// ValidateBasic performs a stateless check of the proposal fields
func (rip *RemoveAssetFromWhitelistProposal) ValidateBasic() error {
	return govv1beta1.ValidateAbstract(rip)
}
