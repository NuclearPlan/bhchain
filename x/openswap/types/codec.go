package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
)

var ModuleCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&TradingPair{}, "bhexchain/openswap/TradingPair", nil)
	cdc.RegisterConcrete(&AddrLiquidity{}, "bhexchain/openswap/AddrLiquidity", nil)
	cdc.RegisterConcrete(&TradeReward{}, "bhexchain/openswap/TradeReward", nil)
	cdc.RegisterConcrete(MsgCreateDex{}, "bhexchain/openswap/MsgCreateDex", nil)
	cdc.RegisterConcrete(MsgEditDex{}, "bhexchain/openswap/MsgEditDex", nil)
	cdc.RegisterConcrete(MsgCreateTradingPair{}, "bhexchain/openswap/MsgCreateTradingPair", nil)
	cdc.RegisterConcrete(MsgEditTradingPair{}, "bhexchain/openswap/MsgEditTradingPair", nil)
	cdc.RegisterConcrete(MsgAddLiquidity{}, "bhexchain/openswap/MsgAddLiquidity", nil)
	cdc.RegisterConcrete(MsgRemoveLiquidity{}, "bhexchain/openswap/MsgRemoveLiquidity", nil)
	cdc.RegisterConcrete(MsgSwapExactIn{}, "bhexchain/openswap/MsgSwapExactIn", nil)
	cdc.RegisterConcrete(MsgSwapExactOut{}, "bhexchain/openswap/MsgSwapExactOut", nil)
	cdc.RegisterConcrete(MsgLimitSwap{}, "bhexchain/openswap/MsgLimitSwap", nil)
	cdc.RegisterConcrete(MsgCancelLimitSwap{}, "bhexchain/openswap/MsgCancelLimitSwap", nil)
	cdc.RegisterConcrete(MsgClaimLPEarning{}, "bhexchain/openswap/MsgClaimLPEarning", nil)
	cdc.RegisterConcrete(MsgClaimTradeEarning{}, "bhexchain/openswap/MsgClaimTradeEarning", nil)
	cdc.RegisterConcrete(&Order{}, "bhexchain/openswap/Order", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
