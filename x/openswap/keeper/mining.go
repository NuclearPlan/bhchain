package keeper

import (
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/openswap/types"
)

func (k Keeper) CalculateLPEarning(ctx sdk.Context, addr sdk.CUAddress, dexID uint32, tokenA, tokenB sdk.Symbol) sdk.Int {
	tokenA, tokenB, _ = k.SortTokens(ctx, tokenA, tokenB)
	liquidity := k.GetLiquidity(ctx, addr, dexID, tokenA, tokenB)
	globalMask := k.getDec(ctx, types.GlobalMaskKey(dexID, tokenA, tokenB))
	addrMask := k.getDec(ctx, types.AddrMaskKey(addr, dexID, tokenA, tokenB))
	return globalMask.Mul(liquidity.ToDec()).Sub(addrMask).TruncateInt()
}

func (k Keeper) ClaimLPEarning(ctx sdk.Context, addr sdk.CUAddress, dexID uint32, tokenA, tokenB sdk.Symbol) sdk.Result {
	tokenA, tokenB, _ = k.SortTokens(ctx, tokenA, tokenB)
	earning := k.CalculateLPEarning(ctx, addr, dexID, tokenA, tokenB)
	if !earning.IsPositive() {
		return sdk.Result{}
	}

	result, err := k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, earning)))
	if err != nil {
		return err.Result()
	}

	flows := k.getFlowFromResult(&result)

	addrMaskKey := types.AddrMaskKey(addr, dexID, tokenA, tokenB)
	addrMask := k.getDec(ctx, addrMaskKey)
	addrMask = addrMask.Add(earning.ToDec())
	k.setDec(ctx, addrMaskKey, addrMask)

	receipt := k.rk.NewReceipt(sdk.CategoryTypeOpenswap, flows)
	result = sdk.Result{}
	k.rk.SaveReceiptToResult(receipt, &result)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawLPEarning,
			sdk.NewAttribute(types.AttributeKeyAddress, addr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, earning.String()),
		),
	})
	result.Events = append(result.Events, ctx.EventManager().Events()...)

	return result
}

func (k Keeper) CalculateTradeEarning(ctx sdk.Context, addr sdk.CUAddress, dexID uint32, tokenA, tokenB sdk.Symbol) (sdk.Int, sdk.Int, *types.TradeReward) {
	addrVolume := k.getAddrTradeVolume(ctx, dexID, tokenA, tokenB, addr)
	poolReward := k.GetPairTradeReward(ctx, dexID, tokenA, tokenB)
	earning := sdk.ZeroInt()
	if poolReward.TradingVolume.IsPositive() {
		earning = mulAndDiv(addrVolume, poolReward.Reward, poolReward.TradingVolume)
	}
	return earning, addrVolume, poolReward
}

func (k Keeper) ClaimTradeEarning(ctx sdk.Context, addr sdk.CUAddress, dexID uint32, tokenA, tokenB sdk.Symbol) sdk.Result {
	earning, addrVolume, poolReward := k.CalculateTradeEarning(ctx, addr, dexID, tokenA, tokenB)
	if !earning.IsPositive() {
		return sdk.Result{}
	}

	result, err := k.sk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, earning)))
	if err != nil {
		return err.Result()
	}
	flows := k.getFlowFromResult(&result)

	k.delAddrTradeVolume(ctx, dexID, tokenA, tokenB, addr)
	poolReward.TradingVolume = poolReward.TradingVolume.Sub(addrVolume)
	k.setPairTradeReward(ctx, poolReward)

	receipt := k.rk.NewReceipt(sdk.CategoryTypeOpenswap, flows)
	result = sdk.Result{}
	k.rk.SaveReceiptToResult(receipt, &result)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdrawTradeEarning,
			sdk.NewAttribute(types.AttributeKeyAddress, addr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, earning.String()),
		),
	})
	result.Events = append(result.Events, ctx.EventManager().Events()...)

	return result
}

func (k Keeper) Mining(ctx sdk.Context) {
	bonusCollector := k.sk.GetModuleAccount(ctx, types.BonusCollectorName)
	amount := k.tk.GetBalance(ctx, bonusCollector.GetAddress(), sdk.NativeToken)
	if !amount.IsPositive() {
		return
	}

	k.sk.SendCoinsFromAccountToModule(ctx, bonusCollector.GetAddress(), types.ModuleName, sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, amount)))

	lpAmountRate := k.LPMiningRewardRate(ctx)
	lpAmount := amount.ToDec().Mul(lpAmountRate).TruncateInt()
	lpAmount = k.distributeLpReward(ctx, lpAmount)

	tradeAmount := amount.Sub(lpAmount)
	tradeAmount = k.distributeTradeReward(ctx, tradeAmount)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMining,
			sdk.NewAttribute(types.AttributeKeyLPAmount, lpAmount.String()),
			sdk.NewAttribute(types.AttributeKeyTradeAmount, tradeAmount.String()),
		),
	})
}

func (k Keeper) distributeLpReward(ctx sdk.Context, amount sdk.Int) sdk.Int {
	if amount.IsZero() {
		return sdk.ZeroInt()
	}
	totalWeight := sdk.ZeroInt()
	miningWeights := k.LPMiningWeights(ctx)
	for _, w := range miningWeights {
		totalWeight = totalWeight.Add(w.Weight)
	}
	if totalWeight.IsZero() {
		return sdk.ZeroInt()
	}

	remaining := amount
	for i, w := range miningWeights {
		distribution := remaining
		if i < len(miningWeights)-1 {
			distribution = amount.Mul(w.Weight).Quo(totalWeight)
			remaining = remaining.Sub(distribution)
		}
		tokenA, tokenB, _ := k.SortTokens(ctx, w.TokenA, w.TokenB)
		k.onMining(ctx, w.DexID, tokenA, tokenB, distribution)
	}
	return amount
}

func (k Keeper) distributeTradeReward(ctx sdk.Context, amount sdk.Int) sdk.Int {
	if amount.IsZero() {
		return sdk.ZeroInt()
	}
	totalWeight := sdk.ZeroInt()
	miningWeights := k.TradeMiningWeights(ctx)
	for _, w := range miningWeights {
		totalWeight = totalWeight.Add(w.Weight)
	}
	if totalWeight.IsZero() {
		return sdk.ZeroInt()
	}

	remaining := amount
	for i, w := range miningWeights {
		distribution := remaining
		if i < len(miningWeights)-1 {
			distribution = amount.Mul(w.Weight).Quo(totalWeight)
			remaining = remaining.Sub(distribution)
		}
		tokenA, tokenB, _ := k.SortTokens(ctx, w.TokenA, w.TokenB)
		tradeReward := k.GetPairTradeReward(ctx, w.DexID, tokenA, tokenB)
		tradeReward.Reward = tradeReward.Reward.Add(distribution)
		k.setPairTradeReward(ctx, tradeReward)
	}
	return amount
}

func (k Keeper) onUpdateLiquidity(ctx sdk.Context, addr sdk.CUAddress, dexID uint32, tokenA, tokenB sdk.Symbol, liquidity sdk.Int) {
	// update total share
	totalShareKey := types.TotalShareKey(dexID, tokenA, tokenB)
	totalShare := k.getDec(ctx, totalShareKey)
	totalShare = totalShare.Add(liquidity.ToDec())
	k.setDec(ctx, totalShareKey, totalShare)

	// update addr mast
	globalMask := k.getDec(ctx, types.GlobalMaskKey(dexID, tokenA, tokenB))
	if globalMask.IsPositive() {
		addrMaskKey := types.AddrMaskKey(addr, dexID, tokenA, tokenB)
		addrMask := k.getDec(ctx, addrMaskKey)
		addrMask = addrMask.Add(globalMask.Mul(liquidity.ToDec()))
		k.setDec(ctx, addrMaskKey, addrMask)
	}
}

func (k Keeper) onMining(ctx sdk.Context, dexID uint32, tokenA, tokenB sdk.Symbol, amount sdk.Int) {
	totalShare := k.getDec(ctx, types.TotalShareKey(dexID, tokenA, tokenB))
	if totalShare.IsPositive() {
		globalMaskKey := types.GlobalMaskKey(dexID, tokenA, tokenB)
		globalMask := k.getDec(ctx, globalMaskKey)
		globalMask = globalMask.Add(amount.ToDec().Quo(totalShare))
		k.setDec(ctx, globalMaskKey, globalMask)
	}
}

func (k Keeper) GetPairTradeReward(ctx sdk.Context, dexID uint32, tokenA, tokenB sdk.Symbol) *types.TradeReward {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PairTradeRewardKey(dexID, tokenA, tokenB))
	if len(bz) == 0 {
		return types.NewTradeReward(dexID, tokenA, tokenB)
	}
	var reward types.TradeReward
	k.cdc.MustUnmarshalBinaryBare(bz, &reward)
	return &reward
}

func (k Keeper) setPairTradeReward(ctx sdk.Context, reward *types.TradeReward) {
	bz := k.cdc.MustMarshalBinaryBare(reward)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PairTradeRewardKey(reward.DexID, reward.TokenA, reward.TokenB), bz)
}

func (k Keeper) delAddrTradeVolume(ctx sdk.Context, dexID uint32, tokenA, tokenB sdk.Symbol, addr sdk.CUAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AddrTradeVolumeKey(addr, dexID, tokenA, tokenB))
}

func (k Keeper) getAddrTradeVolume(ctx sdk.Context, dexID uint32, tokenA, tokenB sdk.Symbol, addr sdk.CUAddress) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AddrTradeVolumeKey(addr, dexID, tokenA, tokenB))
	if len(bz) == 0 {
		return sdk.ZeroInt()
	}
	var ret sdk.Int
	k.cdc.MustUnmarshalBinaryBare(bz, &ret)
	return ret
}

func (k Keeper) setAddrTradeVolume(ctx sdk.Context, dexID uint32, tokenA, tokenB sdk.Symbol, addr sdk.CUAddress, volume sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(volume)
	store.Set(types.AddrTradeVolumeKey(addr, dexID, tokenA, tokenB), bz)
}

func (k Keeper) onTrade(ctx sdk.Context, dexID uint32, tokenA, tokenB sdk.Symbol, addr sdk.CUAddress, volume sdk.Int) {
	poolReward := k.GetPairTradeReward(ctx, dexID, tokenA, tokenB)
	if poolReward.Reward.IsZero() {
		return
	}
	poolReward.TradingVolume = poolReward.TradingVolume.Add(volume)
	k.setPairTradeReward(ctx, poolReward)

	addrVolume := k.getAddrTradeVolume(ctx, dexID, tokenA, tokenB, addr)
	addrVolume = addrVolume.Add(volume)
	k.setAddrTradeVolume(ctx, dexID, tokenA, tokenB, addr, addrVolume)
}
