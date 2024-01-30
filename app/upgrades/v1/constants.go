package v01

import (
	store "cosmossdk.io/store/types"
	"github.com/sideprotocol/side/app/upgrades"
	gmmypes "github.com/sideprotocol/side/x/gmm/types"
)

// UpgradeName defines the on-chain upgrade name for the Osmosis v21 upgrade.
const (
	UpgradeName = "v01"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			gmmypes.ModuleName,
		},
		Deleted: []string{},
	},
}
