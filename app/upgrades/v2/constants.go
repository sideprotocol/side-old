package v02

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/sideprotocol/side/app/upgrades"
)

const (
	UpgradeName = "v02"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
