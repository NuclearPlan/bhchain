// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/bluehelix-chain/bhchain/x/bank/types
package transfer

import (
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/transfer/keeper"
	"github.com/bluehelix-chain/bhchain/x/transfer/types"
)

const (
	DefaultCodespace         = types.DefaultCodespace
	CodeSendDisabled         = types.CodeSendDisabled
	CodeInvalidInputsOutputs = types.CodeInvalidInputsOutputs
	ModuleName               = types.ModuleName
	RouterKey                = types.RouterKey
	QuerierRoute             = types.QuerierRoute
	DefaultParamspace        = types.DefaultParamspace
	DefaultSendEnabled       = types.DefaultSendEnabled

	EventTypeTransfer      = types.EventTypeTransfer
	AttributeKeyRecipient  = types.AttributeKeyRecipient
	AttributeKeySender     = types.AttributeKeySender
	AttributeValueCategory = types.AttributeValueCategory
	UtxoBased              = sdk.UtxoBased
	StoreKey               = types.StoreKey
)

var (
	// functions aliases
	RegisterCodec          = types.RegisterCodec
	ErrNoInputs            = types.ErrNoInputs
	ErrNoOutputs           = types.ErrNoOutputs
	ErrInputOutputMismatch = types.ErrInputOutputMismatch
	ErrSendDisabled        = types.ErrSendDisabled
	NewBaseKeeper          = keeper.NewBaseKeeper
	NewInput               = types.NewInput
	NewOutput              = types.NewOutput
	ParamKeyTable          = types.ParamKeyTable

	// variable aliases
	ModuleCdc                = types.ModuleCdc
	ParamStoreKeySendEnabled = types.ParamStoreKeySendEnabled
)

type (
	BaseKeeper                     = keeper.BaseKeeper // ibc module depends on this
	Keeper                         = keeper.Keeper
	Input                          = types.Input
	Output                         = types.Output
	MsgSend                        = types.MsgSend
	MsgMultiSend                   = types.MsgMultiSend
	MsgDeposit                     = types.MsgDeposit
	MsgConfirmedDeposit            = types.MsgConfirmedDeposit
	MsgCollectWaitSign             = types.MsgCollectWaitSign
	MsgCollectSignFinish           = types.MsgCollectSignFinish
	MsgCollectFinish               = types.MsgCollectFinish
	MsgRecollect                   = types.MsgRecollect
	MsgWithdrawal                  = types.MsgWithdrawal
	MsgWithdrawalConfirm           = types.MsgWithdrawalConfirm
	MsgWithdrawalWaitSign          = types.MsgWithdrawalWaitSign
	MsgWithdrawalSignFinish        = types.MsgWithdrawalSignFinish
	MsgWithdrawalFinish            = types.MsgWithdrawalFinish
	MsgSysTransfer                 = types.MsgSysTransfer
	MsgSysTransferWaitSign         = types.MsgSysTransferWaitSign
	MsgSysTransferSignFinish       = types.MsgSysTransferSignFinish
	MsgSysTransferFinish           = types.MsgSysTransferFinish
	MsgOpcuAssetTransfer           = types.MsgOpcuAssetTransfer
	MsgOpcuAssetTransferWaitSign   = types.MsgOpcuAssetTransferWaitSign
	MsgOpcuAssetTransferSignFinish = types.MsgOpcuAssetTransferSignFinish
	MsgOpcuAssetTransferFinish     = types.MsgOpcuAssetTransferFinish
	MsgOrderRetry                  = types.MsgOrderRetry
	MsgCancelWithdrawal            = types.MsgCancelWithdrawal
	MsgForceUpdateCUNonce          = types.MsgForceUpdateCUNonce
	MsgForceCancelWithdrawal       = types.MsgForceCancelWithdrawal
)
