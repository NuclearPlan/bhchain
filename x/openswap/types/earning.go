package types

import sdk "github.com/bluehelix-chain/bhchain/types"

type Earning struct {
	DexID                  uint32     `json:"dex_id"`
	TokenA                 sdk.Symbol `json:"token_a"`
	TokenB                 sdk.Symbol `json:"token_b"`
	Amount                 sdk.Int    `json:"amount"`
	UnclaimedTradingVolume sdk.Int    `json:"unclaimed_trading_volume"`
}

func NewEarning(dexID uint32, tokenA, tokenB sdk.Symbol, amount, unclaimedTradingVolume sdk.Int) *Earning {
	return &Earning{
		DexID:                  dexID,
		TokenA:                 tokenA,
		TokenB:                 tokenB,
		Amount:                 amount,
		UnclaimedTradingVolume: unclaimedTradingVolume,
	}
}
