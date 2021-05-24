package token

import (
	"bytes"
	"fmt"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/token/types"
)

type GenesisState struct {
	GenesisTokens []sdk.Token `json:"genesis_tokens"`
}

func ValidateGenesis(data GenesisState) error {
	for _, token := range data.GenesisTokens {
		if !token.IsValid() {
			return fmt.Errorf("token %s is invalid", token.GetSymbol())
		}
	}

	return nil
}

func DefaultGenesisState() GenesisState {
	baseTokens := []*sdk.BaseToken{
		{
			Name:        sdk.NativeToken,
			Symbol:      sdk.Symbol(sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    sdk.NativeTokenDecimal,
			TotalSupply: sdk.NewIntWithDecimal(21, 24),
			Weight:      types.DefaultNativeTokenWeight,
			MapToken:    true,
		},
		{
			Name:        sdk.NativeUsdtToken,
			Symbol:      types.CalSymbol("usdt", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(1, 27),
			Weight:      types.DefaultStableCoinWeight,
			MapToken:    true,
		},
		{
			Name:        "btc",
			Symbol:      types.CalSymbol("btc", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    8,
			TotalSupply: sdk.NewIntWithDecimal(21, 14),
			Weight:      types.DefaultIBCTokenWeight,
			MapToken:    true,
		},
		{
			Name:        "eth",
			Symbol:      types.CalSymbol("eth", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(1, 27),
			Weight:      types.DefaultIBCTokenWeight + 1,
			MapToken:    true,
		},
		{
			Name:        "trx",
			Symbol:      types.CalSymbol("trx", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    6,
			TotalSupply: sdk.NewIntWithDecimal(1, 17),
			Weight:      types.DefaultIBCTokenWeight + 2,
			MapToken:    true,
		},
		{
			Name:        "ht",
			Symbol:      types.CalSymbol("ht", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(5, 26),
			Weight:      types.DefaultIBCTokenWeight + 2,
			MapToken:    true,
		},
		{
			Name:        "bnb",
			Symbol:      types.CalSymbol("bnb", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(2, 26),
			Weight:      types.DefaultIBCTokenWeight + 2,
			MapToken:    true,
		},
	}
	ibcTokens := []*sdk.IBCToken{
		{
			BaseToken: sdk.BaseToken{
				Name:        "btc",
				Symbol:      sdk.Symbol("btc"),
				Issuer:      "",
				Chain:       sdk.Symbol("btc"),
				SendEnabled: true,
				Decimals:    8,
				TotalSupply: sdk.NewIntWithDecimal(21, 14),
				Weight:      types.DefaultIBCTokenWeight,
			},
			TokenType:          sdk.UtxoBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 6),  // btc
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.ZeroInt(),                // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1), // gas * 3
			OpCUSysTransferNum: sdk.NewInt(1), // SysTransferAmount * 10
			GasLimit:           sdk.NewInt(1),
			GasPrice:           sdk.NewInt(10000),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 6),
			Confirmations:      6,
			IsNonceBased:       false,
			MappingSymbol:      types.CalSymbol("btc", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "eth",
				Symbol:      sdk.Symbol("eth"),
				Issuer:      "",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 1,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 17), // 0.1eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(1),
			GasLimit:           sdk.NewInt(21000),
			GasPrice:           sdk.NewInt(10000),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 17), // 0.1eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("eth", sdk.NativeToken),
			IsNonceBased:       true,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "hbc",
				Issuer:      "0x28Da24ed20906CDE186D8B4f83412C3AE37a6269",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(21, 24),
				Weight:      types.DefaultNativeTokenWeight,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 18), // 10 usdt
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),     //
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18), //10 usdt
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      sdk.NativeToken,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "usdt",
				Issuer:      "0xdac17f958d2ee523a2206206994597c13d831ec7", // TODO (diff testnet & mainnet) (0xdAC17F958D2ee523a2206206994597C13D831ec7)
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    6,
				TotalSupply: sdk.NewIntWithDecimal(1, 17),
				Weight:      types.DefaultStableCoinWeight,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(10, 6), // 10 usdt
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),     //
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(10, 6), //10 usdt
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("usdt", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "link",
				Issuer:      "0x514910771af9ca656af840dff83e8264ecf986ca",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 4,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(5, 18),
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewInt(1000),
			DepositThreshold:   sdk.NewIntWithDecimal(5, 18),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("link", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "trx",
				Symbol:      sdk.Symbol("trx"),
				Issuer:      "",
				Chain:       sdk.Symbol("trx"),
				SendEnabled: true,
				Decimals:    6,
				TotalSupply: sdk.NewIntWithDecimal(1, 17),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(100, 6), // 100 trx
			OpenFee:            sdk.NewIntWithDecimal(1, 16),  // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17),  // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),       //1x gas
			OpCUSysTransferNum: sdk.NewInt(5),       //5x gas
			GasLimit:           sdk.NewInt(1000000), //  1tron
			GasPrice:           sdk.NewInt(1),
			DepositThreshold:   sdk.NewIntWithDecimal(100, 6), // same as btc
			Confirmations:      20,
			IsNonceBased:       false,
			MappingSymbol:      types.CalSymbol("trx", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "usdt",
				Issuer:      "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
				Chain:       sdk.Symbol("trx"),
				SendEnabled: true,
				Decimals:    6,
				TotalSupply: sdk.NewIntWithDecimal(1, 17),
				Weight:      types.DefaultStableCoinWeight,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(10, 6), // 10 usdt
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),     // 1 tron
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),       //1x gas
			OpCUSysTransferNum: sdk.NewInt(5),       //5x gas
			GasLimit:           sdk.NewInt(6000000), //  6trx
			GasPrice:           sdk.NewInt(1),
			DepositThreshold:   sdk.NewIntWithDecimal(10, 6), // 10 TRXUSDT
			Confirmations:      20,
			IsNonceBased:       false,
			MappingSymbol:      types.CalSymbol("usdt", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		//heco
		{
			BaseToken: sdk.BaseToken{
				Name:        "ht",
				Symbol:      sdk.Symbol("ht"),
				Issuer:      "",
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(5, 26),
				Weight:      types.DefaultIBCTokenWeight + 1,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(2, 17), // // TODO (diff testnet & mainnet)
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(21000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(2, 17),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("ht", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "hbc",
				Issuer:      "0x894b2917c783514c9e4c24229bf60f3cb4c9c905",
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(21, 24),
				Weight:      types.DefaultNativeTokenWeight,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 18),
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(80000),           //  ht
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      sdk.NativeToken,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "usdt",
				Issuer:      "0xa71edc38d189767582c38a3145b5873052c3e47a", // TODO (diff testnet & mainnet)
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(5, 26),
				Weight:      types.DefaultStableCoinWeight,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 18), // TODO (diff testnet & mainnet)
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(80000),           //  ht
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("usdt", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		//bsc
		{
			BaseToken: sdk.BaseToken{
				Name:        "bnb",
				Symbol:      sdk.Symbol("bnb"),
				Issuer:      "",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(2, 26),
				Weight:      types.DefaultIBCTokenWeight + 1,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 18), // 1 bnb for test
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(21000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("bnb", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "usdt",
				Issuer:      "0x55d398326f99059ff775485246999027b3197955",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(5, 26),
				Weight:      types.DefaultStableCoinWeight,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 18), // 1 dai for test
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(80000),           //  bnb
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("usdt", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "doge",
				Symbol:      sdk.Symbol("doge"),
				Issuer:      "",
				Chain:       sdk.Symbol("doge"),
				SendEnabled: true,
				Decimals:    8,
				TotalSupply: sdk.NewIntWithDecimal(21, 14),
				Weight:      types.DefaultIBCTokenWeight + 3,
			},
			TokenType:          sdk.UtxoBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(10, 8),
			OpenFee:            sdk.NewIntWithDecimal(1, 16),
			SysOpenFee:         sdk.ZeroInt(),
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(2),
			GasLimit:           sdk.NewInt(10000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(10, 8),
			Confirmations:      6,
			IsNonceBased:       false,
			MappingSymbol:      types.CalSymbol("doge", sdk.NativeToken),
			NeedCollectFee:     true,
		},
	}
	for _, ibcToken := range ibcTokens {
		if ibcToken.Symbol != ibcToken.Chain {
			ibcToken.Symbol = types.CalSymbol(ibcToken.Issuer, ibcToken.Chain)
		}
	}

	genTokens := make([]sdk.Token, 0, len(baseTokens)+len(ibcTokens))
	for _, baseToken := range baseTokens {
		genTokens = append(genTokens, baseToken)
	}
	for _, ibcToken := range ibcTokens {
		genTokens = append(genTokens, ibcToken)
	}

	return GenesisState{
		GenesisTokens: genTokens,
	}
}

func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, token := range data.GenesisTokens {
		err := k.CreateToken(ctx, token)
		if err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var tokens []sdk.Token
	iter := k.GetSymbolIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var tokenInfo sdk.Token
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &tokenInfo)
		tokens = append(tokens, tokenInfo)
	}
	return GenesisState{GenesisTokens: tokens}
}

// Checks whether 2 GenesisState structs are equivalent.
func (g GenesisState) Equal(g2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(g)
	b2 := ModuleCdc.MustMarshalBinaryBare(g2)
	return bytes.Equal(b1, b2)
}

// Returns if a GenesisState is empty or has data in it
func (g GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return g.Equal(emptyGenState)
}

func (g GenesisState) String() string {
	var b strings.Builder

	for _, token := range g.GenesisTokens {
		b.WriteString(token.String())
		b.WriteString("\n")
	}

	return b.String()
}
