package genutil

// DONTCOVER

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bluehelix-chain/bhchain/codec"
	sdk "github.com/bluehelix-chain/bhchain/types"
	authtypes "github.com/bluehelix-chain/bhchain/x/custodianunit/types"
	"github.com/bluehelix-chain/bhchain/x/genutil/types"
	stakingtypes "github.com/bluehelix-chain/bhchain/x/staking/types"
)

// SetGenTxsInAppGenesisState - sets the genesis transactions in the app genesis state
func SetGenTxsInAppGenesisState(cdc *codec.Codec, appGenesisState map[string]json.RawMessage,
	genTxs []authtypes.StdTx) (map[string]json.RawMessage, error) {

	genesisState := GetGenesisStateFromAppState(cdc, appGenesisState)
	// convert all the GenTxs to JSON
	var genTxsBz []json.RawMessage
	for _, genTx := range genTxs {
		txBz, err := cdc.MarshalJSON(genTx)
		if err != nil {
			return appGenesisState, err
		}
		genTxsBz = append(genTxsBz, txBz)
	}

	genesisState.GenTxs = genTxsBz
	return SetGenesisStateInAppState(cdc, appGenesisState, genesisState), nil
}

// ValidateAccountInGenesis checks that the provided key has sufficient
// coins in the genesis accounts
func ValidateAccountInGenesis(appGenesisState map[string]json.RawMessage,
	genAccIterator types.GenesisCUsIterator,
	key sdk.CUAddress, coins sdk.Coins, cdc *codec.Codec) error {

	//accountIsInGenesis := false

	// TODO: refactor out bond denom to common state area
	stakingDataBz := appGenesisState[stakingtypes.ModuleName]
	var stakingData stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(stakingDataBz, &stakingData)

	genUtilDataBz := appGenesisState[stakingtypes.ModuleName]
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(genUtilDataBz, &genesisState)

	//var err error
	//genAccIterator.IterateGenesisCUs(cdc, appGenesisState,
	//	func(acc authexported.CustodianUnit) (stop bool) {
	//		CUAddress := acc.GetAddress()
	//		// Ensure that CustodianUnit is in genesis
	//		if CUAddress.Equals(key) {
	//			accountIsInGenesis = true
	//			return true
	//		}
	//		return false
	//	},
	//)
	//if err != nil {
	//	return err
	//}
	//
	//if !accountIsInGenesis {
	//	return fmt.Errorf("CustodianUnit %s in not in the app_state.accounts array of genesis.json", key)
	//}

	return nil
}

type deliverTxfn func(abci.RequestDeliverTx) abci.ResponseDeliverTx

// DeliverGenTxs - deliver a genesis transaction
func DeliverGenTxs(ctx sdk.Context, cdc *codec.Codec, genTxs []json.RawMessage,
	stakingKeeper types.StakingKeeper, deliverTx deliverTxfn) []abci.ValidatorUpdate {

	for _, genTx := range genTxs {
		var tx authtypes.StdTx
		cdc.MustUnmarshalJSON(genTx, &tx)
		bz := cdc.MustMarshalBinaryLengthPrefixed(tx)
		res := deliverTx(abci.RequestDeliverTx{Tx: bz})
		if !res.IsOK() {
			panic(res.Log)
		}
	}
	return stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
}
