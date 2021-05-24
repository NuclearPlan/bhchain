package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
)

var ModuleCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgNewToken{}, "bhexchain/hrc10/MsgNewToken", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
