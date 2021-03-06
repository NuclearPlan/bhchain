package keeper

import (
	sdk "github.com/bluehelix-chain/bhchain/types"
)

// RegisterInvariants register all supply invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	// ir.RegisterRoute(types.ModuleName, "total-supply", TotalSupply(k))
}

// TotalSupply checks that the total supply reflects all the coins held in accounts
/* func TotalSupply(k Keeper) sdk.Invariant { */
// return func(ctx sdk.Context) (string, bool) {
// var expectedTotal sdk.Coins
// supply := k.GetSupply(ctx)

// k.ck.IterateCUs(ctx, func(acc exported.CustodianUnit) bool {
// expectedTotal = expectedTotal.Add(acc.GetCoins())
// return false
// })

// broken := !expectedTotal.IsEqual(supply.GetTotal())

// return sdk.FormatInvariant(types.ModuleName, "total supply",
// fmt.Sprintf(
// "\tsum of accounts coins: %v\n"+
// "\tsupply.Total:          %v\n",
// expectedTotal, supply.GetTotal())), broken
// }
/* } */
