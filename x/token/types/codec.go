package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
	sdk "github.com/bluehelix-chain/bhchain/types"
)

var ModuleCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&sdk.BaseToken{}, "bhexchain/types/BaseToken", nil)
	cdc.RegisterConcrete(&sdk.IBCToken{}, "bhexchain/types/IBCToken", nil)
	cdc.RegisterInterface((*sdk.Token)(nil), nil)
	cdc.RegisterConcrete(sdk.TokensGasPrice{}, "bhexchain/types/TokensGasPrice", nil)
	cdc.RegisterConcrete(MsgSynGasPrice{}, "bhexchain/token/MsgSynGasPrice", nil)
	cdc.RegisterConcrete(AddTokenProposal{}, "bhexchain/AddTokenProposal", nil)
	cdc.RegisterConcrete(TokenParamsChangeProposal{}, "bhexchain/TokenParamsChangeProposal", nil)
	cdc.RegisterConcrete(&GasPriceVoteBox{}, "bhexchain/token/GasPriceVoteBox", nil)
	cdc.RegisterConcrete(&GasPriceVoteItem{}, "bhexchain/token/GasPriceVoteItem", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
