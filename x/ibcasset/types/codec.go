package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
	"github.com/bluehelix-chain/bhchain/x/ibcasset/exported"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.CUIBCAsset)(nil), nil)
	cdc.RegisterConcrete(&CUIBCAsset{}, "bhexchain/CUIBCAsset", nil)
}

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
