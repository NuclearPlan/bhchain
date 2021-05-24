package types

import (
	"github.com/bluehelix-chain/bhchain/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Plan{}, "bhexchain/Plan", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "bhexchain/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&CancelSoftwareUpgradeProposal{}, "bhexchain/CancelSoftwareUpgradeProposal", nil)
}
