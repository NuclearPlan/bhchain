package keeper

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/bluehelix-chain/bhchain/codec"
	"github.com/bluehelix-chain/bhchain/store"
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/custodianunit"
	"github.com/bluehelix-chain/bhchain/x/distribution/types"
	"github.com/bluehelix-chain/bhchain/x/params"
	"github.com/bluehelix-chain/bhchain/x/receipt"
	"github.com/bluehelix-chain/bhchain/x/staking"
	"github.com/bluehelix-chain/bhchain/x/supply"
	"github.com/bluehelix-chain/bhchain/x/transfer"
)

//nolint: deadcode unused
var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delPk2   = ed25519.GenPrivKey().PubKey()
	delPk3   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.CUAddress(delPk1.Address())
	delAddr2 = sdk.CUAddress(delPk2.Address())
	delAddr3 = sdk.CUAddress(delPk3.Address())

	valOpPk1    = ed25519.GenPrivKey().PubKey()
	valOpPk2    = ed25519.GenPrivKey().PubKey()
	valOpPk3    = ed25519.GenPrivKey().PubKey()
	valOpAddr1  = sdk.ValAddress(valOpPk1.Address())
	valOpAddr2  = sdk.ValAddress(valOpPk2.Address())
	valOpAddr3  = sdk.ValAddress(valOpPk3.Address())
	valAccAddr1 = sdk.CUAddress(valOpPk1.Address()) // generate acc addresses for these validator keys too
	valAccAddr2 = sdk.CUAddress(valOpPk2.Address())
	valAccAddr3 = sdk.CUAddress(valOpPk3.Address())

	valConsPk1   = ed25519.GenPrivKey().PubKey()
	valConsPk2   = ed25519.GenPrivKey().PubKey()
	valConsPk3   = ed25519.GenPrivKey().PubKey()
	valConsAddr1 = sdk.ConsAddress(valConsPk1.Address())
	valConsAddr2 = sdk.ConsAddress(valConsPk2.Address())
	valConsAddr3 = sdk.ConsAddress(valConsPk3.Address())

	// TODO move to common testing package for all modules
	// test addresses
	TestAddrs = []sdk.CUAddress{
		delAddr1, delAddr2, delAddr3,
		valAccAddr1, valAccAddr2, valAccAddr3,
	}

	emptyDelAddr sdk.CUAddress
	emptyValAddr sdk.ValAddress
	emptyPubkey  crypto.PubKey

	distrAcc = supply.NewEmptyModuleAccount(types.ModuleName)
)

type mockStakingKeeper struct {
	mock.Mock
	staking.Keeper
}

func (m *mockStakingKeeper) GetCurrentEpoch(ctx sdk.Context) sdk.Epoch {
	args := m.Called(ctx)
	return args.Get(0).(sdk.Epoch)
}

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	transfer.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	custodianunit.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	receipt.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc) // distr
	return cdc
}

// test input with default values
func CreateTestInputDefault(t *testing.T, isCheckTx bool, initPower int64) (
	sdk.Context, custodianunit.CUKeeper, transfer.Keeper, Keeper, *mockStakingKeeper, types.SupplyKeeper) {

	communityTax := sdk.NewDecWithPrec(2, 2)

	ctx, ak, tk, dk, sk, _, supplyKeeper := CreateTestInputAdvanced(t, isCheckTx, initPower, communityTax)
	return ctx, ak, tk, dk, sk, supplyKeeper
}

// hogpodge of all sorts of input required for testing
func CreateTestInputAdvanced(t *testing.T, isCheckTx bool, initPower int64,
	communityTax sdk.Dec) (sdk.Context, custodianunit.CUKeeper, transfer.Keeper,
	Keeper, *mockStakingKeeper, params.Keeper, types.SupplyKeeper) {

	initTokens := sdk.TokensFromConsensusPower(initPower)

	keyDistr := sdk.NewKVStoreKey(types.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(custodianunit.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyTransfer := sdk.NewKVStoreKey(transfer.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyTransfer, sdk.StoreTypeIAVL, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	feeCollectorAcc := supply.NewEmptyModuleAccount(custodianunit.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[notBondedPool.String()] = true
	blacklistedAddrs[bondPool.String()] = true
	blacklistedAddrs[distrAcc.String()] = true

	cdc := MakeTestCodec()
	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	rk := receipt.NewKeeper(cdc)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	cuKeeper := custodianunit.NewCUKeeper(cdc, keyAcc, pk.Subspace(custodianunit.DefaultParamspace), custodianunit.ProtoBaseCU)
	bankKeeper := transfer.NewBaseKeeper(cdc, keyTransfer, cuKeeper, nil, nil, nil, rk, nil, nil, pk.Subspace(transfer.DefaultParamspace), transfer.DefaultCodespace, blacklistedAddrs)
	maccPerms := map[string][]string{
		custodianunit.FeeCollectorName: nil,
		types.ModuleName:               nil,
		staking.NotBondedPoolName:      []string{supply.Burner, supply.Staking},
		staking.BondedPoolName:         []string{supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, cuKeeper, bankKeeper, maccPerms)

	sk := staking.NewKeeper(cdc, keyStaking, tkeyStaking, supplyKeeper, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)

	mockSk := &mockStakingKeeper{
		Keeper: sk,
	}
	mockSk.SetParams(ctx, staking.DefaultParams())

	keeper := NewKeeper(cdc, keyDistr, pk.Subspace(DefaultParamspace), mockSk, supplyKeeper, bankKeeper, types.DefaultCodespace, custodianunit.FeeCollectorName, blacklistedAddrs)

	initCoins := sdk.NewCoins(sdk.NewCoin(mockSk.BondDenom(ctx), initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(mockSk.BondDenom(ctx), initTokens.MulRaw(int64(len(TestAddrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range TestAddrs {
		_, _, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		require.Nil(t, err)
	}

	// set module accounts
	keeper.supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, bondPool)
	keeper.supplyKeeper.SetModuleAccount(ctx, distrAcc)

	// set the distribution hooks on staking
	mockSk.SetHooks(keeper.Hooks())

	// set genesis items required for distribution
	keeper.SetFeePool(ctx, types.InitialFeePool())
	keeper.SetCommunityTax(ctx, communityTax)
	keeper.SetBaseProposerReward(ctx, sdk.NewDecWithPrec(1, 2))
	keeper.SetBonusProposerReward(ctx, sdk.NewDecWithPrec(4, 2))
	keeper.SetKeyNodeReward(ctx, sdk.NewDecWithPrec(5, 2))

	params := mockSk.GetParams(ctx)
	mockSk.SetParams(ctx, params)

	return ctx, cuKeeper, bankKeeper, keeper, mockSk, pk, supplyKeeper
}
