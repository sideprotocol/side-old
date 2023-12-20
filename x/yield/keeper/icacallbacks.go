package keeper

import (
	icacallbackstypes "github.com/Stride-Labs/stride/v16/x/icacallbacks/types"
)

const (
	ICACallbackID_Delegate   = "delegate"
	ICACallbackID_Undelegate = "undelegate"
	//ICACallbackID_Redemption     = "redemption"
)

func (k Keeper) Callbacks() icacallbackstypes.ModuleCallbacks {
	return []icacallbackstypes.ICACallback{
		{CallbackId: ICACallbackID_Delegate, CallbackFunc: icacallbackstypes.ICACallbackFunction(k.DelegateCallback)},
		{CallbackId: ICACallbackID_Undelegate, CallbackFunc: icacallbackstypes.ICACallbackFunction(k.UndelegateCallback)},
		//{CallbackId: ICACallbackID_Redemption, CallbackFunc: icacallbackstypes.ICACallbackFunction(k.RedemptionCallback)},
	}
}
