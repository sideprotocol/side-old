package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSubmitWithdrawStatus = "submit_withdraw_status"

func NewMsgSubmitWithdrawStatusRequest(
	sender string,
	txid string,
	status SigningStatus,
) *MsgSubmitWithdrawStatusRequest {
	return &MsgSubmitWithdrawStatusRequest{
		Sender: sender,
		Txid:   txid,
		Status: status,
	}
}

func (msg *MsgSubmitWithdrawStatusRequest) Route() string {
	return RouterKey
}

func (msg *MsgSubmitWithdrawStatusRequest) Type() string {
	return TypeMsgSubmitWithdrawStatus
}

func (msg *MsgSubmitWithdrawStatusRequest) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgSubmitWithdrawStatusRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitWithdrawStatusRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid Sender address (%s)", err)
	}

	if len(msg.Txid) == 0 {
		return sdkerrors.Wrap(ErrSigningRequestNotExist, "txid cannot be empty")
	}

	if msg.Status != SigningStatus_SIGNING_STATUS_BROADCASTED {
		return sdkerrors.Wrap(ErrInvalidStatus, "invalid status")
	}

	return nil
}
