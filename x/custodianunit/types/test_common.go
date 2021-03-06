// nolint noalias
package types

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/bluehelix-chain/bhchain/types"
)

func NewTestMsg(addrs ...sdk.CUAddress) *sdk.TestMsg {
	return sdk.NewTestMsg(addrs...)
}

func NewTestStdFee() StdFee {
	return NewStdFee(50000,
		sdk.NewCoins(sdk.NewInt64Coin("atom", 150)),
	)
}

// coins to more than cover the fee
func NewTestCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin("atom", 10000000),
	}
}

func KeyTestPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.CUAddress) {
	key := secp256k1.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.CUAddress(pub.Address())
	return key, pub, addr
}

func NewTestTx(ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, seqs []uint64, fee StdFee) sdk.Tx {
	sigs := make([]StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := StdSignBytes(ctx.ChainID(), seqs[i], fee, msgs, "")

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := NewStdTx(msgs, fee, sigs, "")
	return tx
}

func NewTestTxWithMemo(ctx sdk.Context, msgs []sdk.Msg, privs []crypto.PrivKey, seqs []uint64, fee StdFee, memo string) sdk.Tx {
	sigs := make([]StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := StdSignBytes(ctx.ChainID(), seqs[i], fee, msgs, memo)

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := NewStdTx(msgs, fee, sigs, memo)
	return tx
}

func NewTestTxWithSignBytes(msgs []sdk.Msg, privs []crypto.PrivKey, seqs []uint64, fee StdFee, signBytes []byte, memo string) sdk.Tx {
	sigs := make([]StdSignature, len(privs))
	for i, priv := range privs {
		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := NewStdTx(msgs, fee, sigs, memo)
	return tx
}
