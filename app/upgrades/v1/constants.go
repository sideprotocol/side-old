package v01

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/sideprotocol/side/app/upgrades"
	gmmypes "github.com/sideprotocol/side/x/gmm/types"
	yieldmoduletypes "github.com/sideprotocol/side/x/yield/types"
)

const (
	UpgradeName = "v01"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			gmmypes.ModuleName,
			yieldmoduletypes.ModuleName,
		},
		Deleted: []string{},
	},
}
