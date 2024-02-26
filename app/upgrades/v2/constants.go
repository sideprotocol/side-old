package v02

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	packetforwardtypes "github.com/sideprotocol/packet-forward-middleware/v7/packetforward/types"
	"github.com/sideprotocol/side/app/upgrades"
)

const (
	UpgradeName = "v02"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			packetforwardtypes.ModuleName,
		},
		Deleted: []string{},
	},
}
