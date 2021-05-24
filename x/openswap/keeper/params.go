package keeper

import (
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/openswap/types"
	"github.com/bluehelix-chain/bhchain/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace = types.ModuleName
)

// ParamTable for staking module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (k Keeper) MinimumLiquidity(ctx sdk.Context) (res sdk.Int) {
	k.paramstore.Get(ctx, types.KeyMinimumLiquidity, &res)
	return
}

func (k Keeper) LimitSwapMatchingGas(ctx sdk.Context) (res sdk.Uint) {
	k.paramstore.Get(ctx, types.KeyLimitSwapMatchingGas, &res)
	return
}

func (k Keeper) MaxFeeRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMaxFeeRate, &res)
	return
}

func (k Keeper) LpRewardRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyLpRewardRate, &res)
	return
}

func (k Keeper) RepurchaseRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyRepurchaseRate, &res)
	return
}

func (k Keeper) RefererTransactionBonusRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyRefererTransactionBonusRate, &res)
	return
}

func (k Keeper) RepurchaseDuration(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyRepurchaseDuration, &res)
	return
}

func (k Keeper) RepurchaseToken(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyRepurchaseToken, &res)
	return
}

func (k Keeper) RepurchaseRoutingToken(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyRepurchaseRoutingToken, &res)
	return
}

func (k Keeper) LPMiningWeights(ctx sdk.Context) (res []*types.MiningWeight) {
	k.paramstore.Get(ctx, types.KeyLPMiningWeights, &res)
	return
}

func (k Keeper) TradeMiningWeights(ctx sdk.Context) (res []*types.MiningWeight) {
	k.paramstore.Get(ctx, types.KeyTradeMiningWeights, &res)
	return
}

func (k Keeper) LPMiningRewardRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyLPMiningRewardRate, &res)
	return
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.MinimumLiquidity(ctx),
		k.LimitSwapMatchingGas(ctx),
		k.MaxFeeRate(ctx),
		k.LpRewardRate(ctx),
		k.RepurchaseRate(ctx),
		k.RefererTransactionBonusRate(ctx),
		k.RepurchaseDuration(ctx),
		k.LPMiningWeights(ctx),
		k.TradeMiningWeights(ctx),
		k.LPMiningRewardRate(ctx),
		k.RepurchaseToken(ctx),
		k.RepurchaseRoutingToken(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
