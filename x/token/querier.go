package token

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/token/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryToken:
			return queryToken(ctx, req, keeper)
		case types.QueryIBCTokens:
			return queryIBCTokens(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown nameservice query endpoint")
		}
	}
}

//=====
func queryToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var ti QueryTokenInfoParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &ti); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	symbol := sdk.Symbol(ti.Symbol)
	token := keeper.GetToken(ctx, symbol)
	if token == nil {
		return nil, sdk.ErrUnknownRequest("Non-exits symbol")
	}

	collectFee := sdk.ZeroInt()
	withDrawalFee := sdk.ZeroInt()
	gasPrice := sdk.ZeroInt()
	if token.IsIBCToken() {
		collectFee = keeper.GetCollectFee(ctx, symbol).Amount
		withDrawalFee = keeper.GetWithDrawalFee(ctx, symbol).Amount
		if token.GetChain() != token.GetSymbol() {
			chainToken := keeper.GetIBCToken(ctx, token.GetChain())
			gasPrice = chainToken.GasPrice
		} else {
			tokenInfo := token.(*sdk.IBCToken)
			gasPrice = tokenInfo.GasPrice
		}
	}

	bz := keeper.cdc.MustMarshalJSON(types.NewResToken(token, collectFee, withDrawalFee, gasPrice))
	return bz, nil
}

func queryIBCTokens(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	tokens := keeper.GetIBCTokenList(ctx)
	res := make([]*ResToken, len(tokens))
	for i, token := range tokens {
		collectFee := sdk.ZeroInt()
		withDrawalFee := sdk.ZeroInt()
		if token.IsIBCToken() {
			collectFee = keeper.GetCollectFee(ctx, token.Symbol).Amount
			withDrawalFee = keeper.GetWithDrawalFee(ctx, token.Symbol).Amount
		}
		if token.Chain != token.Symbol {
			chainToken := keeper.GetIBCToken(ctx, token.Chain)
			res[i] = types.NewResToken(token, collectFee, withDrawalFee, chainToken.GasPrice)
		} else {
			res[i] = types.NewResToken(token, collectFee, withDrawalFee, token.GasPrice)
		}
	}
	bz, err := keeper.cdc.MarshalJSON(res)
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
