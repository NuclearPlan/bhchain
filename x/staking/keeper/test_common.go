package keeper // noalias

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/bluehelix-chain/bhchain/codec"
	"github.com/bluehelix-chain/bhchain/store"
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/custodianunit"
	"github.com/bluehelix-chain/bhchain/x/params"
	"github.com/bluehelix-chain/bhchain/x/receipt"
	"github.com/bluehelix-chain/bhchain/x/staking/types"
	"github.com/bluehelix-chain/bhchain/x/supply"
	"github.com/bluehelix-chain/bhchain/x/supply/exported"
	"github.com/bluehelix-chain/bhchain/x/transfer"
)

// dummy addresses used for testing
// nolint: unused deadcode
var (
	Addrs = createTestAddrs(500)
	PKs   = createTestPubKeys(500)

	addrDels = []sdk.CUAddress{
		Addrs[0],
		Addrs[1],
	}
	addrVals = []sdk.CUAddress{
		Addrs[2],
		Addrs[3],
		Addrs[4],
		Addrs[5],
		Addrs[6],
		Addrs[7],
		Addrs[8],
		Addrs[9],
		Addrs[10],
		Addrs[11],
	}
)

type testCU struct {
	exported.ModuleAccountI
	ctx       sdk.Context
	trk       transfer.Keeper
	coinCache *sdk.Coins
}

func NewTestCU(ctx sdk.Context, trk transfer.Keeper, cu exported.ModuleAccountI) *testCU {
	t := &testCU{ModuleAccountI: cu, ctx: ctx, trk: trk}
	t.GetCoins()
	return t
}

func (t *testCU) SetCoins(coins sdk.Coins) error {
	curCoins := t.trk.GetAllBalance(t.ctx, t.ModuleAccountI.GetAddress())
	t.trk.SubCoins(t.ctx, t.ModuleAccountI.GetAddress(), curCoins)
	t.trk.AddCoins(t.ctx, t.ModuleAccountI.GetAddress(), coins)
	newCoins := coins
	t.coinCache = &newCoins
	return nil
}

func (t *testCU) GetCoins() sdk.Coins {
	if t.coinCache == nil {
		c := t.trk.GetAllBalance(t.ctx, t.ModuleAccountI.GetAddress())
		t.coinCache = &c
	}

	return *t.coinCache
}

//_______________________________________________________________________________________

// intended to be used with require/assert:  require.True(ValEq(...))
func ValEq(t *testing.T, exp, got types.Validator) (*testing.T, bool, string, types.Validator, types.Validator) {
	return t, exp.TestEquivalent(got), "expected:\t%v\ngot:\t\t%v", exp, got
}

//_______________________________________________________________________________________

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(transfer.MsgSend{}, "test/staking/Send", nil)
	cdc.RegisterConcrete(types.MsgCreateValidator{}, "test/staking/CreateValidator", nil)
	cdc.RegisterConcrete(types.MsgEditValidator{}, "test/staking/EditValidator", nil)
	cdc.RegisterConcrete(types.MsgUndelegate{}, "test/staking/Undelegate", nil)
	cdc.RegisterConcrete(types.MsgBeginRedelegate{}, "test/staking/BeginRedelegate", nil)
	cdc.RegisterConcrete(testCU{}, "test/staking/testCU", nil)

	// Register AppAccount
	//cdc.RegisterInterface((*custodianunit.CU)(nil), nil)
	//cdc.RegisterConcrete(&custodianunit.BaseCU{}, "test/staking/BaseCU", nil)
	receipt.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	custodianunit.RegisterCodec(cdc)

	return cdc
}

// Hogpodge of all sorts of input required for testing.
// `initPower` is converted to an amount of tokens.
// If `initPower` is 0, no addrs get created.
func CreateTestInput(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, custodianunit.CUKeeper, Keeper, types.SupplyKeeper) {
	ctx, ck, sk, stk, _ := CreateTestInputEx(t, isCheckTx, initPower)
	return ctx, ck, sk, stk
}

func CreateTestInputEx(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, custodianunit.CUKeeper, Keeper, types.SupplyKeeper, transfer.Keeper) {
	keyStaking := sdk.NewKVStoreKey(types.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(types.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(custodianunit.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyTransfer := sdk.NewKVStoreKey(transfer.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyTransfer, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := MakeTestCodec()

	feeCollectorAcc := supply.NewEmptyModuleAccount(custodianunit.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(types.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(types.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[notBondedPool.String()] = true
	blacklistedAddrs[bondPool.String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	rk := receipt.NewKeeper(cdc)
	cuKeeper := custodianunit.NewCUKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(custodianunit.DefaultParamspace),
		custodianunit.ProtoBaseCU, // prototype
	)

	bk := transfer.NewBaseKeeper(cdc, keyTransfer,
		cuKeeper, nil, nil, nil, rk, nil, nil,
		pk.Subspace(transfer.DefaultParamspace),
		transfer.DefaultCodespace,
		blacklistedAddrs,
	)

	maccPerms := map[string][]string{
		custodianunit.FeeCollectorName: nil,
		types.NotBondedPoolName:        []string{supply.Burner, supply.Staking},
		types.BondedPoolName:           []string{supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, cuKeeper, bk, maccPerms)

	initTokens := sdk.TokensFromConsensusPower(initPower)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(int64(len(Addrs)))))

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	keeper := NewKeeper(cdc, keyStaking, tkeyStaking, supplyKeeper, pk.Subspace(DefaultParamspace), types.DefaultCodespace)
	keeper.SetTransferKeeper(bk)
	keeper.SetParams(ctx, types.DefaultParams())

	// set module accounts
	_, _, err = bk.AddCoins(ctx, notBondedPool.GetAddress(), totalSupply)
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range Addrs {
		_, _, err := bk.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}

	return ctx, cuKeeper, keeper, supplyKeeper, bk
}

func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

// for incode address generation
func TestAddr(addr string, bech string) sdk.CUAddress {

	res, err := sdk.CUAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	//bechres, err := sdk.CUAddressFromBase58(bech)
	bechres, err := sdk.CUAddressFromBase58(bech)

	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}

	return res
}

// nolint: unparam
func createTestAddrs(numAddrs int) []sdk.CUAddress {
	var addresses []sdk.CUAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.CUAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, TestAddr(buffer.String(), bech))
		buffer.Reset()
	}
	return addresses
}

// nolint: unparam
func createTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	//start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") //base pubkey string
		buffer.WriteString(numString)                                                       //adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

//_____________________________________________________________________________________

// does a certain by-power index record exist
func ValidatorByPowerIndexExists(ctx sdk.Context, keeper Keeper, power []byte) bool {
	store := ctx.KVStore(keeper.storeKey)
	return store.Has(power)
}

// update validator for testing
func TestingUpdateValidator(keeper Keeper, ctx sdk.Context, validator types.Validator, apply bool) types.Validator {
	keeper.SetValidator(ctx, validator)

	// Remove any existing power key for validator.
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsByPowerIndexKey)
	defer iterator.Close()
	deleted := false
	for ; iterator.Valid(); iterator.Next() {
		valAddr := types.ParseValidatorPowerRankKey(iterator.Key())
		if bytes.Equal(valAddr, validator.OperatorAddress) {
			if deleted {
				panic("found duplicate power index key")
			} else {
				deleted = true
			}
			store.Delete(iterator.Key())
		}
	}

	keeper.SetValidatorByPowerIndex(ctx, validator)
	if apply {
		keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		validator, found := keeper.GetValidator(ctx, validator.OperatorAddress)
		if !found {
			panic("validator expected but not found")
		}
		return validator
	}
	cachectx, _ := ctx.CacheContext()
	keeper.ApplyAndReturnValidatorSetUpdates(cachectx)
	validator, found := keeper.GetValidator(cachectx, validator.OperatorAddress)
	if !found {
		panic("validator expected but not found")
	}
	return validator
}

// nolint: deadcode unused
func validatorByPowerIndexExists(k Keeper, ctx sdk.Context, power []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(power)
}

// RandomValidator returns a random validator given access to the keeper and ctx
func RandomValidator(r *rand.Rand, keeper Keeper, ctx sdk.Context) types.Validator {
	vals := keeper.GetAllValidators(ctx)
	i := r.Intn(len(vals))
	return vals[i]
}

// RandomBondedValidator returns a random bonded validator given access to the keeper and ctx
func RandomBondedValidator(r *rand.Rand, keeper Keeper, ctx sdk.Context) types.Validator {
	vals := keeper.GetBondedValidatorsByPower(ctx)
	i := r.Intn(len(vals))
	return vals[i]
}
