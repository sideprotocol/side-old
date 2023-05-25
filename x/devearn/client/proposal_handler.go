package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"sidechain/x/devearn/client/cli"
)

var (
	RegisterDevEarnProposalHandler  = govclient.NewProposalHandler(cli.NewRegisterDevEarnProposalCmd)
	CancelDevEarnProposalHandler    = govclient.NewProposalHandler(cli.NewCancelDevEarnProposalCmd)
	AddAssetToWhitelistHandler      = govclient.NewProposalHandler(cli.NewAddAssetToWhitelistProposalCmd)
	RemoveAssetFromWhitelistHandler = govclient.NewProposalHandler(cli.NewRemoveAssetFromWhitelistProposalCmd)
)
