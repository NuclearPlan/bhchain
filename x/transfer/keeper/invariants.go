package keeper

import (
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/transfer/types"
)

// RegisterInvariants registers the bank module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, ak types.CUKeeper) {
	/*  ir.RegisterRoute(types.ModuleName, "nonnegative-outstanding", */
	/* NonnegativeBalanceInvariant(ak)) */
}
