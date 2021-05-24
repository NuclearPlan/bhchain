package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
	cutypes "github.com/bluehelix-chain/bhchain/x/custodianunit/types"
)

var ModuleCdc = codec.New()

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(cutypes.StdSignature{}, "bhexchain/keygen/StdSignature", nil)
	cdc.RegisterConcrete(MsgKeyGen{}, "bhexchain/keygen/MsgKeyGen", nil)
	cdc.RegisterConcrete(MsgKeyGenWaitSign{}, "bhexchain/keygen/MsgKeyGenWaitSign", nil)
	cdc.RegisterConcrete(MsgKeyGenFinish{}, "bhexchain/keygen/MsgKeyGenFinish", nil)
	cdc.RegisterConcrete(MsgPreKeyGen{}, "bhexchain/keygen/MsgPreKeyGen", nil)
	cdc.RegisterConcrete(MsgOpcuMigrationKeyGen{}, "bhexchain/keygen/MsgOpcuMigrationKeyGen", nil)
	cdc.RegisterConcrete(MsgNewOpCU{}, "bhexchain/keygen/MsgNewOpCU", nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
