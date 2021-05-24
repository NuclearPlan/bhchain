package receipt

import (
	"github.com/bluehelix-chain/bhchain/codec"
	sdk "github.com/bluehelix-chain/bhchain/types"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Flow)(nil), nil)
	cdc.RegisterConcrete(BalanceFlow{}, "bhexchain/receipt/BalanceFlow", nil)
	cdc.RegisterConcrete(OrderFlow{}, "bhexchain/receipt/OrderFlow", nil)
	cdc.RegisterConcrete(DepositFlow{}, "bhexchain/receipt/DepositFlow", nil)
	cdc.RegisterConcrete(MappingBalanceFlow{}, "bhexchain/receipt/MappingBalanceFlow", nil)
	cdc.RegisterConcrete(sdk.DepositConfirmedFlow{}, "bhexchain/receipt/DepositConfrimedFlow", nil)
	cdc.RegisterConcrete(sdk.CollectWaitSignFlow{}, "bhexchain/receipt/CollectWaitSignFlow", nil)
	cdc.RegisterConcrete(sdk.CollectSignFinishFlow{}, "bhexchain/receipt/CollectSignFinishFlow", nil)
	cdc.RegisterConcrete(sdk.CollectFinishFlow{}, "bhexchain/receipt/CollectFinishFlow", nil)
	cdc.RegisterConcrete(sdk.RecollectFlow{}, "bhexchain/receipt/RecollectFlow", nil)
	cdc.RegisterConcrete(sdk.WithdrawalFlow{}, "bhexchain/receipt/WithdrawalFlow", nil)
	cdc.RegisterConcrete(sdk.WithdrawalConfirmFlow{}, "bhexchain/receipt/WithdrawalConfirmFlow", nil)
	cdc.RegisterConcrete(sdk.WithdrawalWaitSignFlow{}, "bhexchain/receipt/WithdrawalWaitSignFlow", nil)
	cdc.RegisterConcrete(sdk.WithdrawalSignFinishFlow{}, "bhexchain/receipt/WithdrawalSignFinishFlow", nil)
	cdc.RegisterConcrete(sdk.WithdrawalFinishFlow{}, "bhexchain/receipt/WithdrawalFinishFlow", nil)
	cdc.RegisterConcrete(sdk.SysTransferFlow{}, "bhexchain/receipt/SysTransferFlow", nil)
	cdc.RegisterConcrete(sdk.SysTransferWaitSignFlow{}, "bhexchain/receipt/SysTransferWaitSignFlow", nil)
	cdc.RegisterConcrete(sdk.SysTransferSignFinishFlow{}, "bhexchain/receipt/SysTransferSignFinishFlow", nil)
	cdc.RegisterConcrete(sdk.SysTransferFinishFlow{}, "bhexchain/receipt/SysTransferFinishFlow", nil)
	cdc.RegisterConcrete(sdk.OpcuAssetTransferFlow{}, "bhexchain/receipt/OpcuAssetTransferFlow", nil)
	cdc.RegisterConcrete(sdk.OpcuAssetTransferWaitSignFlow{}, "bhexchain/receipt/OpcuAssetTransferWaitSignFlow", nil)
	cdc.RegisterConcrete(sdk.OpcuAssetTransferSignFinishFlow{}, "bhexchain/receipt/OpcuAssetTransferSignFinishFlow", nil)
	cdc.RegisterConcrete(sdk.OpcuAssetTransferFinishFlow{}, "bhexchain/receipt/OpcuAssetTransferFinishFlow", nil)

	cdc.RegisterConcrete(sdk.KeyGenFlow{}, "bhexchain/receipt/KeyGenFlow", nil)
	cdc.RegisterConcrete(sdk.KeyGenWaitSignFlow{}, "bhexchain/receipt/KeyGenWaitSignFlow", nil)
	cdc.RegisterConcrete(sdk.KeyGenFinishFlow{}, "bhexchain/receipt/KeyGenFinishFlow", nil)
	cdc.RegisterConcrete(sdk.OrderRetryFlow{}, "bhexchain/receipt/OrderRetryFlow", nil)
}
