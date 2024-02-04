package wasmbinding

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sideprotocol/side/wasmbinding/bindings"
	gmmkeeper "github.com/sideprotocol/side/x/gmm/keeper"
)

type QueryPlugin struct {
	gmmkeeper *gmmkeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(tfk *gmmkeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		gmmkeeper: tfk,
	}
}

// GetParams is a query to get params.
func (qp QueryPlugin) GetParams(ctx sdk.Context) (*bindings.ParamsResponse, error) {
	params := qp.gmmkeeper.GetParams(ctx)

	return &bindings.ParamsResponse{Params: (*bindings.ParamsRes)(&params)}, nil
}
