package types

const (
	EventValueActionCreatePool = "create_pool"
	EventValueActionCancelPool = "cancel_pool"
	EventValueActionDeposit    = "deposit"
	EventValueActionWithdraw   = "withdraw"
	EventValueActionSwap       = "swap_pool"
)

const (
	AttributeKeyAction              = "action"
	AttributeKeyPoolId              = "pool_id"
	AttributeKeyMultiDepositOrderId = "order_id"
	AttributeKeyMultiDeposits       = "deposit"
	AttributeKeyTokenIn             = "token_in"
	AttributeKeyTokenOut            = "token_out"
	AttributeKeyLpToken             = "liquidity_pool_token"
	AttributeKeyLpSupply            = "liquidity_pool_token_supply"
	AttributeKeyPoolCreator         = "pool_creator"
	AttributeKeyOrderCreator        = "order_creator"
	AttributeKeyName                = "name"
	AttributeKeyPoolStatus          = "pool_status"
	AttributeKeyMsgSender           = "msg_sender"
)
