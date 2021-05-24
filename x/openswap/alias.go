package openswap

import (
	"github.com/bluehelix-chain/bhchain/x/openswap/keeper"
	"github.com/bluehelix-chain/bhchain/x/openswap/types"
)

const (
	ModuleName         = types.ModuleName
	RouterKey          = types.RouterKey
	StoreKey           = types.StoreKey
	QuerierKey         = types.QuerierKey
	DefaultParamspace  = types.DefaultParamspace
	BonusCollectorName = types.BonusCollectorName
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc
	NewKeeper     = keeper.NewKeeper
)

type (
	TradingPair = types.TradingPair
	Keeper      = keeper.Keeper
)
