package v01

import (
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
		// versionSetter := keepers.UpgradeKeeper.GetVersionSetter()
		// versionSetter.SetProtocolVersion(2672694)
		// _, correctTypecast := mm.Modules[types.ModuleName].(gmmmodule.AppModule)
		// if !correctTypecast {
		// 	panic("mm.Modules[gmm.ModuleName] is not of type ica.AppModule")
		// }

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
