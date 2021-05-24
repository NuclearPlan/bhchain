package types

import (
	sdk "github.com/bluehelix-chain/bhchain/types"
	authexported "github.com/bluehelix-chain/bhchain/x/custodianunit/exported"
)

// CUKeeper defines the expected CustodianUnit keeper (noalias)
type CUKeeper interface {
	NewCU(sdk.Context, authexported.CustodianUnit) authexported.CustodianUnit
	SetCU(sdk.Context, authexported.CustodianUnit)
	IterateCUs(ctx sdk.Context, process func(authexported.CustodianUnit) (stop bool))
}

type TransferKeeper interface {
	AddCoins(sdk.Context, sdk.CUAddress, sdk.Coins) (sdk.Coins, []sdk.Flow, sdk.Error)
	AddCoinsHold(sdk.Context, sdk.CUAddress, sdk.Coins) (sdk.Coins, []sdk.Flow, sdk.Error)
}
