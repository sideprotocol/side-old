package v01

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	keepers "github.com/sideprotocol/side/app/keepers"
	//gmmmodule "github.com/sideprotocol/side/x/gmm"
	//"github.com/sideprotocol/side/x/gmm/types"
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
		fmt.Println("yield verison", vm["yield"])
		fmt.Println("gmm verison", vm["gmm"])
		if vm["yield"] == 0 {
			vm["yield"] = 2672694 // Set to expected version
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
