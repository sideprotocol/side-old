package v01

import (
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	keepers "github.com/sideprotocol/side/app/keepers"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v7
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		logger := ctx.Logger().With("upgrade", UpgradeName)
		logger.Debug("running module migrations ...")
		// Check if the yield module is new and set its version
		if vm["yield"] == 0 {
			vm["gmm"] = 2672694 // Set to expected version
		}
		_ = keepers
		return mm.RunMigrations(ctx, configurator, vm)
	}
}
