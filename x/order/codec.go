package order

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
	cdc.RegisterInterface((*sdk.Order)(nil), nil)
	cdc.RegisterConcrete(&sdk.OrderBase{}, "bhexchain/order/OrderBase", nil)
	cdc.RegisterConcrete(&sdk.OrderCollect{}, "bhexchain/order/OrderCollect", nil)
	cdc.RegisterConcrete(&sdk.OrderWithdrawal{}, "bhexchain/order/OrderWithdrawal", nil)
	cdc.RegisterConcrete(&sdk.OrderKeyGen{}, "bhexchain/order/OrderKeyGen", nil)
	cdc.RegisterConcrete(&sdk.OrderSysTransfer{}, "bhexchain/order/OrderSysTransfer", nil)
	cdc.RegisterConcrete(&sdk.OrderOpcuAssetTransfer{}, "bhexchain/order/OrderOpcuAssetTransfer", nil)
	cdc.RegisterConcrete(&sdk.TxFinishNodeData{}, "bhexchain/order/TxFinishNodeData", nil)
}
