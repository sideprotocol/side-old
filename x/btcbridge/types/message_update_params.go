package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateParams = "update_params"

// Route returns the route of MsgUpdateParams.
func (msg *MsgUpdateParamsRequest) Route() string {
	return RouterKey
}

// Type returns the type of MsgUpdateParams.
func (msg *MsgUpdateParamsRequest) Type() string {
	return TypeMsgUpdateParams
}

// GetSignBytes implements the LegacyMsg interface.
func (m MsgUpdateParamsRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUpdateParamsRequest) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic performs basic MsgUpdateParams message validation.
func (m *MsgUpdateParamsRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	if err := m.Params.Validate(); err != nil {
		return err
	}

	return nil
}
