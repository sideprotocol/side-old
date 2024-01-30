package v00

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// CreateUpgradeHandler creates an SDK upgrade handler for v7
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
