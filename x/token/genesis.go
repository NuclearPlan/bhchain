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
			Name:        "uni",
			Symbol:      types.CalSymbol("uni", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(1, 27),
			Weight:      types.DefaultIBCTokenWeight + 2,
			MapToken:    true,
		},
		{
			Name:        "sushi",
			Symbol:      types.CalSymbol("sushi", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(1, 27),
			Weight:      types.DefaultIBCTokenWeight + 2,
			MapToken:    true,
		},
		{
			Name:        "aave",
			Symbol:      types.CalSymbol("aave", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(1, 27),
			Weight:      types.DefaultIBCTokenWeight + 2,
			MapToken:    true,
		},
		{
			Name:        "comp",
			Symbol:      types.CalSymbol("comp", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(1, 27),
			Weight:      types.DefaultIBCTokenWeight + 2,
			MapToken:    true,
		},
		{
			Name:        "link",
			Symbol:      types.CalSymbol("link", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
			TotalSupply: sdk.NewIntWithDecimal(1, 27),
			Weight:      types.DefaultIBCTokenWeight + 2,
			MapToken:    true,
		},
		{
			Name:        "btc",
			Symbol:      types.CalSymbol("btc", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    18,
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
			Decimals:    18,
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
		{
			Name:        "doge",
			Symbol:      types.CalSymbol("doge", sdk.NativeToken),
			Issuer:      "",
			Chain:       sdk.Symbol(sdk.NativeToken),
			SendEnabled: true,
			Decimals:    8,
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
		//eth
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
			CollectThreshold:   sdk.NewIntWithDecimal(1, 16), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(21000),
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 16), // 0.01eth
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
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
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
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
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
			CollectThreshold:   sdk.NewIntWithDecimal(1, 18),
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("link", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "wbtc",
				Issuer:      "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    8,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 6),
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 6),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("btc", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "uni",
				Issuer:      "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 4,
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
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("uni", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "aave",
				Issuer:      "0x7fc66500c84a76ad7e9c93437bfc5ac33e2ddae9",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 17),
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 17),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("aave", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "bnb",
				Issuer:      "0xB8c77482e45F1F44dE1745F52C74426C631bDD52",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 17),
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 17),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("bnb", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "sushi",
				Issuer:      "0x6b3595068778dd592e39a122f4f5a5cf09c90fe2",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
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
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("sushi", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "comp",
				Issuer:      "0xc00e94cb662c3520282e6f5717214004a7f26888",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 17),
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 17),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("comp", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "ht",
				Issuer:      "0x6f259637dcd74c767781e37bc6133cd6a68aa161",
				Chain:       sdk.Symbol("eth"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
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
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000), //  eth
			GasPrice:           sdk.NewIntWithDecimal(20, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      15,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("ht", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		//tron
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
			OpCUSysTransferNum: sdk.NewInt(10),      //5x gas
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
			OpCUSysTransferNum: sdk.NewInt(10),      //5x gas
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
			OpCUSysTransferNum: sdk.NewInt(10),
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
				Name:        "hbtc",
				Issuer:      "0x66a79d23e58475d2738179ca52cd0b41d73f0bea",
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(21, 14),
				Weight:      types.DefaultIBCTokenWeight,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 16), // btc
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.ZeroInt(),                // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),  // gas * 3
			OpCUSysTransferNum: sdk.NewInt(10), // SysTransferAmount * 10
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 16),
			Confirmations:      20,
			IsNonceBased:       false,
			MappingSymbol:      types.CalSymbol("btc", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "eth",
				Issuer:      "0x64ff637fb478863b7468bc97d30a5bf3a428a1fd",
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 1,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 16), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 16), // 0.1eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("eth", sdk.NativeToken),
			IsNonceBased:       true,
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
			OpCUSysTransferNum: sdk.NewInt(10),
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
				Name:        "doge",
				Issuer:      "0x40280e26a572745b1152a54d1d44f365daa51618",
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    8,
				TotalSupply: sdk.NewIntWithDecimal(1, 16),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(10, 8),
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),           //  ht
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(10, 8),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("doge", sdk.NativeToken),
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
			CollectThreshold:   sdk.NewIntWithDecimal(10, 18), // TODO (diff testnet & mainnet)
			OpenFee:            sdk.NewIntWithDecimal(1, 16),  // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17),  // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),           //  ht
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(10, 18),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("usdt", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "uni",
				Issuer:      "0x22c54ce8321a4015740ee1109d9cbc25815c46e6", // TODO (diff testnet & mainnet)
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(5, 26),
				Weight:      types.DefaultIBCTokenWeight + 2,
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
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),           //  ht
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("uni", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "sushi",
				Issuer:      "0x52e00b2da5bd7940ffe26b609a42f957f31118d5", // TODO (diff testnet & mainnet)
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(5, 26),
				Weight:      types.DefaultIBCTokenWeight + 2,
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
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),           //  ht
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("sushi", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "aave",
				Issuer:      "0x202b4936fe1a82a4965220860ae46d7d3939bb25", // TODO (diff testnet & mainnet)
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(5, 26),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 17), // TODO (diff testnet & mainnet)
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),           //  ht
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(1, 17),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("aave", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "comp",
				Issuer:      "0xCe0A5CA134fb59402B723412994B30E02f083842", // TODO (diff testnet & mainnet)
				Chain:       sdk.Symbol("ht"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(5, 26),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 16), // TODO (diff testnet & mainnet)
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),           //  ht
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(1, 16),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("comp", sdk.NativeToken),
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
			OpCUSysTransferNum: sdk.NewInt(10),
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
			CollectThreshold:   sdk.NewIntWithDecimal(10, 18), // 1 dai for test
			OpenFee:            sdk.NewIntWithDecimal(1, 16),  // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17),  // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),           //  bnb
			GasPrice:           sdk.NewIntWithDecimal(5, 9), //5GWei
			DepositThreshold:   sdk.NewIntWithDecimal(10, 18),
			Confirmations:      20,
			IsNonceBased:       true,
			MappingSymbol:      types.CalSymbol("usdt", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "btcb",
				Issuer:      "0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(21, 14),
				Weight:      types.DefaultIBCTokenWeight,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 16), // btc
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.ZeroInt(),                // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),  // gas * 3
			OpCUSysTransferNum: sdk.NewInt(10), // SysTransferAmount * 10
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 16),
			Confirmations:      6,
			IsNonceBased:       false,
			MappingSymbol:      types.CalSymbol("btc", sdk.NativeToken),
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "eth",
				Issuer:      "0x2170ed0880ac9a755fd29b2688956bd959f933f8",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 1,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 16), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 16), // 0.01eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("eth", sdk.NativeToken),
			IsNonceBased:       true,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "doge",
				Issuer:      "0xba2ae424d960c26247dd6c32edc70b295c744c43",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    8,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(10, 8), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(10, 8), // 0.01eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("doge", sdk.NativeToken),
			IsNonceBased:       true,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "trx",
				Issuer:      "0x85eac5ac2f758618dfa09bdbe0cf174e7d574d5b",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(100, 18), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16),   // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17),   // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(100, 18), // 0.01eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("trx", sdk.NativeToken),
			IsNonceBased:       true,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "uni",
				Issuer:      "0xbf5140a22578168fd562dccf235e5d43a02ce9b1",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 18), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18), // 0.01eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("uni", sdk.NativeToken),
			IsNonceBased:       true,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "aave",
				Issuer:      "0xfb6115445bff7b52feb98650c87f44907e58f802",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 17), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 17), // 0.01eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("aave", sdk.NativeToken),
			IsNonceBased:       true,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "suhi",
				Issuer:      "0x52e00b2da5bd7940ffe26b609a42f957f31118d5",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 18), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 18), // 0.01eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("sushi", sdk.NativeToken),
			IsNonceBased:       true,
			NeedCollectFee:     true,
		},
		{
			BaseToken: sdk.BaseToken{
				Name:        "comp",
				Issuer:      "0x52CE071Bd9b1C4B00A0b92D298c512478CaD67e8",
				Chain:       sdk.Symbol("bnb"),
				SendEnabled: true,
				Decimals:    18,
				TotalSupply: sdk.NewIntWithDecimal(1, 27),
				Weight:      types.DefaultIBCTokenWeight + 2,
			},
			TokenType:          sdk.AccountBased,
			DepositEnabled:     true,
			WithdrawalEnabled:  true,
			CollectThreshold:   sdk.NewIntWithDecimal(1, 16), // 0.01eth
			OpenFee:            sdk.NewIntWithDecimal(1, 16), // nativeToken
			SysOpenFee:         sdk.NewIntWithDecimal(1, 17), // nativeToken
			WithdrawalFeeRate:  sdk.NewDecWithPrec(1, 0),
			MaxOpCUNumber:      4,
			SysTransferNum:     sdk.NewInt(1),
			OpCUSysTransferNum: sdk.NewInt(10),
			GasLimit:           sdk.NewInt(80000),
			GasPrice:           sdk.NewIntWithDecimal(5, 9),
			DepositThreshold:   sdk.NewIntWithDecimal(1, 16), // 0.01eth
			Confirmations:      15,
			MappingSymbol:      types.CalSymbol("comp", sdk.NativeToken),
			IsNonceBased:       true,
			NeedCollectFee:     true,
		},

		//doge
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
			OpCUSysTransferNum: sdk.NewInt(10),
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
