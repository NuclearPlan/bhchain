package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "bhexchain/transfer/MsgSend", nil)
	cdc.RegisterConcrete(MsgMultiSend{}, "bhexchain/transfer/MsgMultiSend", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "bhexchain/transfer/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgConfirmedDeposit{}, "bhexchain/transfer/MsgConfirmedDeposit", nil)
	cdc.RegisterConcrete(MsgCollectWaitSign{}, "bhexchain/transfer/MsgCollectWaitSign", nil)
	cdc.RegisterConcrete(MsgCollectSignFinish{}, "bhexchain/transfer/MsgCollectSignFinish", nil)
	cdc.RegisterConcrete(MsgCollectFinish{}, "bhexchain/transfer/MsgCollectFinish", nil)
	cdc.RegisterConcrete(MsgRecollect{}, "bhexchain/transfer/MsgRecollect", nil)
	cdc.RegisterConcrete(MsgWithdrawal{}, "bhexchain/transfer/MsgWithdrawal", nil)
	cdc.RegisterConcrete(MsgWithdrawalConfirm{}, "bhexchain/transfer/MsgWithdrawalConfirm", nil)
	cdc.RegisterConcrete(MsgWithdrawalWaitSign{}, "bhexchain/transfer/MsgWithdrawalWaitSign", nil)
	cdc.RegisterConcrete(MsgWithdrawalSignFinish{}, "bhexchain/transfer/MsgWithdrawalSignFinish", nil)
	cdc.RegisterConcrete(MsgWithdrawalFinish{}, "bhexchain/transfer/MsgWithdrawalFinish", nil)
	cdc.RegisterConcrete(MsgSysTransfer{}, "bhexchain/transfer/MsgSysTransfer", nil)
	cdc.RegisterConcrete(MsgSysTransferWaitSign{}, "bhexchain/transfer/MsgSysTransferWaitSign", nil)
	cdc.RegisterConcrete(MsgSysTransferSignFinish{}, "bhexchain/transfer/MsgSysTransferSignFinish", nil)
	cdc.RegisterConcrete(MsgSysTransferFinish{}, "bhexchain/transfer/MsgSysTransferFinish", nil)
	cdc.RegisterConcrete(MsgOpcuAssetTransfer{}, "bhexchain/transfer/MsgOpcuAssetTransfer", nil)
	cdc.RegisterConcrete(MsgOpcuAssetTransferWaitSign{}, "bhexchain/transfer/MsgOpcuAssetTransferWaitSign", nil)
	cdc.RegisterConcrete(MsgOpcuAssetTransferSignFinish{}, "bhexchain/transfer/MsgOpcuAssetTransferSignFinish", nil)
	cdc.RegisterConcrete(MsgOpcuAssetTransferFinish{}, "bhexchain/transfer/MsgOpcuAssetTransferFinish", nil)
	cdc.RegisterConcrete(MsgOrderRetry{}, "bhexchain/transfer/MsgOrderRetry", nil)
	cdc.RegisterConcrete(MsgCancelWithdrawal{}, "bhexchain/transfer/MsgCancelWithdrawal", nil)
	cdc.RegisterConcrete(MsgForceUpdateCUNonce{}, "bhexchain/transfer/MsgForceUpdateCUNonce", nil)
	cdc.RegisterConcrete(MsgForceCancelWithdrawal{}, "bhexchain/transfer/MsgForceCancelWithdrawal", nil)
	cdc.RegisterConcrete(&TxVote{}, "bhexchain/transfer/FinishTxVote", nil)
	cdc.RegisterConcrete(&OrderRetryVoteBox{}, "bhexchain/transfer/OrderRetryVoteBox", nil)
	cdc.RegisterConcrete(&OrderRetryVoteItem{}, "bhexchain/transfer/OrderRetryVoteItem", nil)
}

// ModuleCdc - module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
