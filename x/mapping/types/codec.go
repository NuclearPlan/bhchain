package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
)

var ModuleCdc = codec.New()

func init() {
	codec.RegisterCrypto(ModuleCdc)
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MappingInfo{}, "bhexchain/mapping/MappingInfo", nil)
	cdc.RegisterConcrete(AddMappingProposal{}, "bhexchain/mapping/AddMappingProposal", nil)
	cdc.RegisterConcrete(SwitchMappingProposal{}, "bhexchain/mapping/SwitchMappingProposal", nil)
	cdc.RegisterConcrete(MsgMappingSwap{}, "bhexchain/mapping/MsgMappingSwap", nil)
	cdc.RegisterConcrete(MsgCreateDirectSwap{}, "bhexchain/mapping/MsgCreateDirectSwap", nil)
	cdc.RegisterConcrete(MsgCreateFreeSwap{}, "bhexchain/mapping/MsgCreateFreeSwap", nil)
	cdc.RegisterConcrete(MsgSwapSymbol{}, "bhexchain/mapping/MsgSwapSymbol", nil)
	cdc.RegisterConcrete(MsgCancelSwap{}, "bhexchain/mapping/MsgCancelSwap", nil)
	cdc.RegisterConcrete(FreeSwapInfo{}, "bhexchain/mapping/FreeSwapInfo", nil)
	cdc.RegisterConcrete(FreeSwapOrder{}, "bhexchain/mapping/FreeSwapOrder", nil)
	cdc.RegisterConcrete(DirectSwapInfo{}, "bhexchain/mapping/DirectSwapInfo", nil)
	cdc.RegisterConcrete(DirectSwapOrder{}, "bhexchain/mapping/DirectSwapOrder", nil)
}
