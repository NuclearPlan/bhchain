package token

import (
	"fmt"

	"github.com/bluehelix-chain/bhchain/codec"
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/evidence/exported"
	"github.com/bluehelix-chain/bhchain/x/token/internal"
	"github.com/bluehelix-chain/bhchain/x/token/types"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	storeKey       sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc            *codec.Codec // The wire codec for binary encoding/decoding
	sk             internal.StakingKeeper
	evidenceKeeper internal.EvidenceKeeper
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) SetStakingKeeper(sk internal.StakingKeeper) {
	k.sk = sk
}

func (k *Keeper) SetEvidenceKeeper(evidenceKeeper internal.EvidenceKeeper) {
	k.evidenceKeeper = evidenceKeeper
}

//Set entire TokenInfo
func (k *Keeper) SetToken(ctx sdk.Context, tokenInfo sdk.Token) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.TokenStoreKey(tokenInfo.GetSymbol()), k.cdc.MustMarshalBinaryBare(tokenInfo))
}

func (k *Keeper) CreateToken(ctx sdk.Context, tokenInfo sdk.Token) error {
	if k.HasToken(ctx, tokenInfo.GetSymbol()) {
		return fmt.Errorf("token %s already exists", tokenInfo.GetSymbol())
	}

	if tokenInfo.IsIBCToken() {
		ibcTokenInfo := tokenInfo.(*sdk.IBCToken)
		if !ibcTokenInfo.MappingSymbol.IsValid() {
			return fmt.Errorf("invalid mapping symbol: %s", ibcTokenInfo.MappingSymbol)
		}
		mapTokenInfo := k.GetToken(ctx, ibcTokenInfo.MappingSymbol)
		if mapTokenInfo == nil {
			mapTokenInfo := &sdk.BaseToken{
				Name:        tokenInfo.GetName(),
				Symbol:      ibcTokenInfo.MappingSymbol,
				Issuer:      "",
				Chain:       sdk.NativeToken,
				SendEnabled: true,
				Decimals:    tokenInfo.GetDecimals(),
				TotalSupply: tokenInfo.GetTotalSupply(),
				Weight:      tokenInfo.GetWeight(),
				MapToken:    true,
			}

			k.SetToken(ctx, mapTokenInfo)
		} else {
			if !mapTokenInfo.IsMapToken() {
				return fmt.Errorf("add not maptoken for ibctoken, not supported")
			}
		}
	}

	k.SetToken(ctx, tokenInfo)
	if tokenInfo.IsIBCToken() {
		symbols := k.GetIBCTokenSymbols(ctx)
		symbols = append(symbols, tokenInfo.GetSymbol())
		k.SetIBCTokenSymbols(ctx, symbols)
	}
	return nil
}

func (k *Keeper) GetToken(ctx sdk.Context, symbol sdk.Symbol) sdk.Token {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenStoreKey(symbol))
	if len(bz) == 0 {
		return nil
	}
	var tokenInfo sdk.Token
	k.cdc.MustUnmarshalBinaryBare(bz, &tokenInfo)
	return tokenInfo
}

func (k *Keeper) GetIBCToken(ctx sdk.Context, symbol sdk.Symbol) *sdk.IBCToken {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TokenStoreKey(symbol))
	if len(bz) == 0 {
		return nil
	}
	var tokenInfo sdk.Token
	k.cdc.MustUnmarshalBinaryBare(bz, &tokenInfo)
	if tokenInfo.IsIBCToken() {
		return tokenInfo.(*sdk.IBCToken)
	}
	return nil
}

func (k *Keeper) GetSymbolIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.TokenStoreKeyPrefix)
}

func (k *Keeper) RawCollectFee(ctx sdk.Context, token, chainToken *sdk.IBCToken) sdk.Int {
	if !token.NeedCollectFee {
		return sdk.ZeroInt()
	}

	var baseGasFee sdk.Int
	if token.TokenType == sdk.UtxoBased {
		baseGasFee = sdk.DefaultUtxoCollectTxSize().Mul(chainToken.GasPrice).QuoRaw(sdk.KiloBytes)
	} else {
		baseGasFee = chainToken.GasPrice.Mul(chainToken.GasLimit)
	}
	if token.Chain != token.Symbol {
		collectFeeAmount := token.SysTransferAmount(chainToken.GasPrice).Add(baseGasFee)
		return collectFeeAmount
	} else {
		return baseGasFee
	}
}

func (k *Keeper) GetCollectFee(ctx sdk.Context, symbol sdk.Symbol) sdk.Coin {
	token := k.GetIBCToken(ctx, symbol)
	feeSymbol := token.MappingSymbol
	var collectFee sdk.Int

	if token.Chain != token.Symbol { //not mainnet
		chainTokenInfo := k.GetIBCToken(ctx, token.Chain)
		chainMapTokenInfo := k.GetToken(ctx, chainTokenInfo.MappingSymbol)
		feeSymbol = chainTokenInfo.MappingSymbol
		collectFee = sdk.CalcAmountWithDecimalDiff(k.RawCollectFee(ctx, token, chainTokenInfo), chainMapTokenInfo.GetDecimals(), chainTokenInfo.GetDecimals())
	} else {
		mapTokenInfo := k.GetToken(ctx, token.MappingSymbol)
		collectFee = sdk.CalcAmountWithDecimalDiff(k.RawCollectFee(ctx, token, token), mapTokenInfo.GetDecimals(), token.GetDecimals())
	}

	return sdk.NewCoin(feeSymbol.String(), collectFee)
}

func (k *Keeper) withdrawFee(ctx sdk.Context, token, chainToken *sdk.IBCToken) sdk.Int {
	var baseGasFee sdk.Int
	if token.TokenType == sdk.UtxoBased {
		baseGasFee = sdk.DefaultUtxoWithdrawTxSize().Mul(chainToken.GasPrice).QuoRaw(sdk.KiloBytes)
	} else {
		baseGasFee = chainToken.GasPrice.Mul(token.GasLimit)
	}
	withdrawalFeeAmt := token.WithdrawalFeeRate.Mul(sdk.NewDecFromInt(baseGasFee)).TruncateInt()
	return withdrawalFeeAmt
}

func (k *Keeper) GetWithDrawalFee(ctx sdk.Context, symbol sdk.Symbol) sdk.Coin {
	token := k.GetIBCToken(ctx, symbol)
	feeSymbol := token.MappingSymbol
	var withDrawalFee sdk.Int

	if token.Chain != token.Symbol { //not mainnet
		chainTokenInfo := k.GetIBCToken(ctx, token.Chain)
		chainMapTokenInfo := k.GetToken(ctx, chainTokenInfo.MappingSymbol)
		feeSymbol = chainTokenInfo.MappingSymbol
		withDrawalFee = sdk.CalcAmountWithDecimalDiff(k.withdrawFee(ctx, token, chainTokenInfo), chainMapTokenInfo.GetDecimals(), chainTokenInfo.GetDecimals())
	} else {
		mapTokenInfo := k.GetToken(ctx, token.MappingSymbol)
		withDrawalFee = sdk.CalcAmountWithDecimalDiff(k.withdrawFee(ctx, token, token), mapTokenInfo.GetDecimals(), token.GetDecimals())
	}

	return sdk.NewCoin(feeSymbol.String(), withDrawalFee)
}

func (k *Keeper) HasToken(ctx sdk.Context, symbol sdk.Symbol) bool {
	store := ctx.KVStore(k.storeKey)
	if !symbol.IsValid() {
		return false
	}
	return store.Has(types.TokenStoreKey(symbol))
}

func (k *Keeper) SetIBCTokenSymbols(ctx sdk.Context, symbols []sdk.Symbol) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(symbols)
	store.Set(types.IBCTokenListKey, bz)
}

func (k *Keeper) GetIBCTokenSymbols(ctx sdk.Context) []sdk.Symbol {
	store := ctx.KVStore(k.storeKey)
	tokenList := make([]sdk.Symbol, 0)
	bz := store.Get(types.IBCTokenListKey)
	k.cdc.MustUnmarshalBinaryBare(bz, &tokenList)
	return tokenList
}

func (k *Keeper) GetIBCTokenList(ctx sdk.Context) []*sdk.IBCToken {
	symbols := k.GetIBCTokenSymbols(ctx)
	tokens := make([]*sdk.IBCToken, len(symbols))
	for i, symbol := range symbols {
		token := k.GetToken(ctx, symbol)
		tokens[i] = token.(*sdk.IBCToken)
	}
	return tokens
}

func (k *Keeper) IsSubToken(ctx sdk.Context, symbol sdk.Symbol) bool {
	token := k.GetToken(ctx, symbol)
	if token == nil {
		return false
	}
	return token.GetChain() != token.GetSymbol()
}

func (k *Keeper) SynGasPrice(ctx sdk.Context, fromAddr string, height uint64, tokensGasPrice []sdk.TokensGasPrice) ([]sdk.TokensGasPrice, sdk.Result) {
	curBlockHeight := uint64(ctx.BlockHeight())
	if height >= curBlockHeight || curBlockHeight-height > sdk.GasPriceBucketWindow {
		return nil, sdk.ErrInvalidTx(fmt.Sprintf("invalid height %d, current block height is %d", height, curBlockHeight)).Result()
	}

	address, err := sdk.CUAddressFromBase58(fromAddr)
	if err != nil {
		return nil, sdk.ErrInvalidAddress(fmt.Sprintf("can't decode addr:%s", fromAddr)).Result()
	}
	bValidator, validatorNum := k.sk.IsActiveKeyNode(ctx, address)
	if validatorNum == 0 {
		return nil, sdk.ErrInsufficientValidatorNum(fmt.Sprintf("validator's number:%v", validatorNum)).Result()
	}
	if !bValidator {
		return nil, sdk.ErrInvalidTx(fmt.Sprintf("FromCu: %v is not a validator", fromAddr)).Result()
	}
	for _, item := range tokensGasPrice {
		if !k.HasToken(ctx, sdk.Symbol(item.Chain)) {
			return nil, sdk.ErrInvalidTx(fmt.Sprintf("Chain %s not exists", item.Chain)).Result()
		}
	}

	validGasPrice := make([]sdk.TokensGasPrice, 0)
	bucket := height / sdk.GasPriceBucketWindow

	for _, item := range tokensGasPrice {
		voteID := fmt.Sprintf("%s-%d", item.Chain, bucket)
		firstConfirmed, _, validVotes := k.evidenceKeeper.VoteWithCustomBox(ctx, voteID, address, item.GasPrice, curBlockHeight, types.NewGasPriceVoteBox)
		if firstConfirmed {
			k.updateGasPrice(ctx, item.Chain, validVotes)
		}
		validGasPrice = append(validGasPrice, item)
	}

	return validGasPrice, sdk.Result{}
}

func (k *Keeper) updateGasPrice(ctx sdk.Context, chain string, validVotes []*exported.VoteItem) {
	totalGasFee := sdk.ZeroInt()
	var count int64
	for _, item := range validVotes {
		price, ok := item.Vote.(sdk.Int)
		if !ok {
			continue
		}
		totalGasFee = totalGasFee.Add(price)
		count++
	}
	if count > 0 {
		averageGasPrice := totalGasFee.QuoRaw(count)
		chainSymbol := sdk.Symbol(chain)
		tokenInfos := k.GetIBCTokenList(ctx)
		for _, tokenInfo := range tokenInfos {
			if tokenInfo.Chain == chainSymbol {
				tokenInfo.GasPrice = averageGasPrice
				k.SetToken(ctx, tokenInfo)
			}
		}
	}
}
