package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"sidechain/x/devearn/client/cli"
)

var (
	RegisterIncentiveProposalHandler = govclient.NewProposalHandler(cli.NewRegisterDevEarnProposalCmd)
	CancelIncentiveProposalHandler   = govclient.NewProposalHandler(cli.NewCancelDevEarnProposalCmd)
)
