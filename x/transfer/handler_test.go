package transfer

import (
	"strings"
	"testing"

	sdk "github.com/bluehelix-chain/bhchain/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"
)

func TestInvalidMsg(t *testing.T) {
	h := NewHandler(BaseKeeper{})

	res := h(sdk.NewContext(nil, abci.Header{}, false, nil), sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "unrecognized bank message type"))
}
