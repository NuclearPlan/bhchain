package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
	"github.com/bluehelix-chain/bhchain/x/supply/exported"
)

// RegisterCodec registers the CustodianUnit types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "bhexchain/ModuleAccount", nil)
	cdc.RegisterConcrete(&Supply{}, "bhexchain/Supply", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}
