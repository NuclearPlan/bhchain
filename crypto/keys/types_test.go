package keys

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bluehelix-chain/bhchain/crypto/keys/hd"
	"github.com/bluehelix-chain/bhchain/types"
)

func Test_writeReadLedgerInfo(t *testing.T) {
	var tmpKey secp256k1.PubKeySecp256k1
	bz, _ := hex.DecodeString("035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A")
	copy(tmpKey[:], bz)

	lInfo := newLedgerInfo("some_name", tmpKey, *hd.NewFundraiserParams(5, types.CoinType, 1))
	assert.Equal(t, TypeLedger, lInfo.GetType())

	path, err := lInfo.GetPath()
	assert.NoError(t, err)
	assert.Equal(t, "44'/496'/5'/0/1", path.String())
	assert.Equal(t,
		"hbcpub1addwnpepqddddqg2glc8x4fl7vxjlnr7p5a3czm5kcdp4239sg6yqdc4rc2r5swxthh",
		types.MustBech32ifyAccPub(lInfo.GetPubKey()))

	// Serialize and restore
	serialized := writeInfo(lInfo)
	restoredInfo, err := readInfo(serialized)
	assert.NoError(t, err)
	assert.NotNil(t, restoredInfo)

	// Check both keys match
	assert.Equal(t, lInfo.GetName(), restoredInfo.GetName())
	assert.Equal(t, lInfo.GetType(), restoredInfo.GetType())
	assert.Equal(t, lInfo.GetPubKey(), restoredInfo.GetPubKey())

	restoredPath, err := restoredInfo.GetPath()
	assert.NoError(t, err)

	assert.Equal(t, path, restoredPath)
}
