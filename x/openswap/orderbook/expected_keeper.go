package orderbook

import (
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/openswap/types"
)

type OpenswapKeeper interface {
	IteratorAllUnfinishedOrder(ctx sdk.Context, f func(*types.Order))
}
