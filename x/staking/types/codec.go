package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "bhexchain/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "bhexchain/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgKeyNodeHeartbeat{}, "bhexchain/MsgKeyNodeHeartbeat", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "bhexchain/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "bhexchain/MsgUndelegate", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "bhexchain/MsgBeginRedelegate", nil)
	cdc.RegisterConcrete(&UpdateKeyNodesProposal{}, "bhexchain/UpdateKeyNodesProposal", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
