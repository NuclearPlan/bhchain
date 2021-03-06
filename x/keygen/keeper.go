package keygen

import (
	"github.com/bluehelix-chain/bhchain/chainnode"
	"github.com/bluehelix-chain/bhchain/codec"
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/keygen/internal"
	"github.com/bluehelix-chain/bhchain/x/keygen/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	tk       internal.TokenKeeper
	ck       internal.CUKeeper
	ik       internal.IBCAssetKeeper
	ok       internal.OrderKeeper
	rk       internal.ReceiptKeeper
	vk       internal.StakingKeeper
	dk       internal.DistributionKeeper
	trk      internal.TransferKeeper
	cn       chainnode.Chainnode
}

func NewKeeper(store sdk.StoreKey, cdc *codec.Codec, tk internal.TokenKeeper, ck internal.CUKeeper, ik internal.IBCAssetKeeper, ok internal.OrderKeeper,
	rk internal.ReceiptKeeper, vk internal.StakingKeeper, dk internal.DistributionKeeper, trk internal.TransferKeeper, cn chainnode.Chainnode) Keeper {
	return Keeper{
		storeKey: store,
		cdc:      cdc,
		tk:       tk,
		ck:       ck,
		ik:       ik,
		ok:       ok,
		rk:       rk,
		vk:       vk,
		dk:       dk,
		trk:      trk,
		cn:       cn,
	}
}

func (k *Keeper) GetWaitAssignKeyGenOrderIDs(ctx sdk.Context) []string {
	orderIDs := []string{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(waitAssignKey())
	if bz != nil {
		k.cdc.MustUnmarshalBinaryBare(bz, &orderIDs)
	}
	return orderIDs
}

func (k *Keeper) AddWaitAssignKeyGenOrderID(ctx sdk.Context, orderID string) error {
	orderIDs := k.GetWaitAssignKeyGenOrderIDs(ctx)
	orderIDs = append(orderIDs, orderID)
	return k.setWaitAssignKeyGenOrderIDs(ctx, orderIDs)
}

func (k *Keeper) DelWaitAssignKeyGenOrderID(ctx sdk.Context, orderID string) error {
	orderIDs := k.GetWaitAssignKeyGenOrderIDs(ctx)
	index := sdk.StringsIndex(orderIDs, orderID)
	if index >= 0 {
		orderIDs = append(orderIDs[:index], orderIDs[index+1:]...)
		return k.setWaitAssignKeyGenOrderIDs(ctx, orderIDs)
	}
	return nil
}

func (k *Keeper) delAllWaitAssignKeyGenOrderIDs(ctx sdk.Context) error {
	orderIDs := k.GetWaitAssignKeyGenOrderIDs(ctx)
	for _, id := range orderIDs {
		order := k.ok.GetOrder(ctx, id)
		keygenOrder, ok := order.(*sdk.OrderKeyGen)
		if !ok {
			continue
		}
		keygenOrder.SetOrderStatus(sdk.OrderStatusFinish)
		k.ok.SetOrder(ctx, keygenOrder)
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(waitAssignKey())
	return nil
}

func (k *Keeper) setWaitAssignKeyGenOrderIDs(ctx sdk.Context, orderIDs []string) error {
	store := ctx.KVStore(k.storeKey)
	if len(orderIDs) == 0 {
		store.Delete(waitAssignKey())
		return nil
	}
	bz := k.cdc.MustMarshalBinaryBare(orderIDs)
	store.Set(waitAssignKey(), bz)
	return nil
}

func (k *Keeper) resetKeyGenOrders(ctx sdk.Context, epoch sdk.Epoch) {
	keyNodes := make([]sdk.CUAddress, len(epoch.KeyNodeSet))
	for i, val := range epoch.KeyNodeSet {
		keyNodes[i] = val
	}
	threshold := uint64(sdk.Majority23(len(keyNodes)))
	for _, orderID := range k.ok.GetProcessOrderList(ctx) {
		order := k.ok.GetOrder(ctx, orderID)
		if order == nil || order.GetOrderType() != sdk.OrderTypeKeyGen {
			continue
		}
		if order.GetOrderStatus() != sdk.OrderStatusBegin && order.GetOrderStatus() != sdk.OrderStatusWaitSign {
			continue
		}
		keygenOrder := order.(*sdk.OrderKeyGen)
		keygenOrder.KeyNodes = keyNodes
		keygenOrder.SignThreshold = threshold
		keygenOrder.SetOrderStatus(sdk.OrderStatusBegin)
		k.ok.SetOrder(ctx, keygenOrder)
	}
}

func (k *Keeper) getExcludedKeyNode(ctx sdk.Context, keyNodes []sdk.CUAddress) sdk.CUAddress {
	var excluded sdk.CUAddress
	var longest int64
	blkHeight := ctx.BlockHeight()
	for _, keyNode := range keyNodes {
		val, _ := k.vk.GetValidator(ctx, sdk.ValAddress(keyNode))
		if longest < blkHeight-int64(val.LastKeyNodeHeartbeatHeight) {
			longest = blkHeight - int64(val.LastKeyNodeHeartbeatHeight)
			excluded = keyNode
		}
	}
	if longest >= MaxKeyNodeHeartbeat {
		return excluded
	}
	return nil
}

func waitAssignKey() []byte {
	return types.WaitAssignKey
}
