package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sideprotocol/side/x/yield/keeper"
	"github.com/sideprotocol/side/x/yield/types"
)

func SimulateMsgRegisterHostChain(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRegisterHostChain{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the RegisterHostChain simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "RegisterHostChain simulation not implemented"), nil, nil
	}
}
