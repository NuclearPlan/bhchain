package mint

import (
	"math/big"

	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/mint/internal/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored minter & params
	params := k.GetParams(ctx)
	currentHeight := ctx.BlockHeight()

	mintForStaking := params.MintForStakingPerBlock
	if params.MintForStakingHalveDuration.IsPositive() {
		epoch := currentHeight / params.MintForStakingHalveDuration.Int64()
		denominator := new(big.Int).Exp(big.NewInt(2), big.NewInt(epoch), nil)
		mintForStaking = mintForStaking.Quo(sdk.NewIntFromBigInt(denominator))
	}
	if mintForStaking.IsPositive() {
		mintedCoins := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, mintForStaking))
		err := k.MintCoins(ctx, mintedCoins)
		if err != nil {
			panic(err)
		}
		err = k.AddCollectedFees(ctx, sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, mintForStaking)))
		if err != nil {
			panic(err)
		}
	}

	mintForOpenswap := sdk.ZeroInt()
	if params.MintForOpenswapBonusStartBlock.IsPositive() && currentHeight >= params.MintForOpenswapBonusStartBlock.Int64() {
		mintForOpenswap = params.MintForOpenswapBonusPerBlock
		if params.MintForOpenswapBonusHalveDuration.IsPositive() {
			epoch := (currentHeight - params.MintForOpenswapBonusStartBlock.Int64()) / params.MintForOpenswapBonusHalveDuration.Int64()
			denominator := new(big.Int).Exp(big.NewInt(2), big.NewInt(epoch), nil)
			mintForOpenswap = mintForOpenswap.Quo(sdk.NewIntFromBigInt(denominator))
		}
	}
	if mintForOpenswap.IsPositive() {
		mintedCoins := sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, mintForOpenswap))
		err := k.MintCoins(ctx, mintedCoins)
		if err != nil {
			panic(err)
		}
		err = k.AddOpenswapBonus(ctx, sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, mintForOpenswap)))
		if err != nil {
			panic(err)
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(sdk.AttributeKeyMintAmount, mintForStaking.Add(mintForOpenswap).String()),
			sdk.NewAttribute(sdk.AttributeKeyStakingBonusAmount, mintForStaking.String()),
			sdk.NewAttribute(sdk.AttributeKeyOpenswapBonusAmount, mintForOpenswap.String()),
		),
	)
}
