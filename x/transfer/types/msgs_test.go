package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/bluehelix-chain/bhchain/types"
)

func TestMsgSendRoute(t *testing.T) {
	addr1 := sdk.CUAddress([]byte("HBCYzvRB1WFY6nUDmCSxWrSdzBJ7wvyG1Bdw"))
	addr2 := sdk.CUAddress([]byte("HBCW1KVVAy28Fb1HLz38jAJ7pg7Ga8N91MQm"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("hbc", 10))
	var msg = NewMsgSend(addr1, addr2, coins)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "send")
}

func TestMsgSendValidation(t *testing.T) {
	addr1, _ := sdk.CUAddressFromBase58("HBCYzvRB1WFY6nUDmCSxWrSdzBJ7wvyG1Bdw")
	addr2, _ := sdk.CUAddressFromBase58("HBCW1KVVAy28Fb1HLz38jAJ7pg7Ga8N91MQm")
	hbc123 := sdk.NewCoins(sdk.NewInt64Coin("hbc", 123))
	hbc0 := sdk.NewCoins(sdk.NewInt64Coin("hbc", 0))
	hbc123eth123 := sdk.NewCoins(sdk.NewInt64Coin("hbc", 123), sdk.NewInt64Coin("eth", 123))
	hbc123eth0 := sdk.Coins{sdk.NewInt64Coin("hbc", 123), sdk.NewInt64Coin("eth", 0)}

	var emptyAddr sdk.CUAddress

	cases := []struct {
		valid bool
		tx    MsgSend
	}{
		{true, NewMsgSend(addr1, addr2, hbc123)},       // valid send
		{true, NewMsgSend(addr1, addr2, hbc123eth123)}, // valid send with multiple coins
		{false, NewMsgSend(addr1, addr2, hbc0)},        // non positive coin
		{false, NewMsgSend(addr1, addr2, hbc123eth0)},  // non positive coin in multicoins
		{false, NewMsgSend(emptyAddr, addr2, hbc123)},  // empty from addr
		{false, NewMsgSend(addr1, emptyAddr, hbc123)},  // empty to addr
	}

	for _, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err)
		} else {
			require.NotNil(t, err)
		}
	}
}

func TestMsgSendGetSignBytes(t *testing.T) {
	addr1 := sdk.CUAddress([]byte("LUSH6CPn8RXTW759FQsK"))
	addr2 := sdk.CUAddress([]byte("Ky7eQidQcEKG1nA24bdP"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("atom", 10))
	var msg = NewMsgSend(addr1, addr2, coins)
	res := msg.GetSignBytes()

	expected := "{\"type\":\"bhexchain/transfer/MsgSend\",\"value\":{\"amount\":[{\"amount\":\"10\",\"denom\":\"atom\"}],\"from_address\":\"HBCT9W8oJnZA4h71QAFyj4we4bF1AACEgxVJ\",\"to_address\":\"HBCT4xTMyxwd2ZwbDiiG1BfNbCs1gwrwHYho\"}}"
	require.Equal(t, expected, string(res))
}

func TestMsgSendGetSigners(t *testing.T) {
	cuAddr, _ := sdk.CUAddressFromBase58("HBCLmQcskpdQivEkRrh1gNPm7c9aVB8hh1fy")
	var msg = NewMsgSend(cuAddr, sdk.CUAddress{}, sdk.NewCoins())
	res := msg.GetSigners()
	// TODO: fix this !
	require.Equal(t, "HBCLmQcskpdQivEkRrh1gNPm7c9aVB8hh1fy", res[0].String())
}

func TestMsgMultiSendRoute(t *testing.T) {
	// Construct a MsgSend
	addr1 := sdk.CUAddress([]byte("input"))
	addr2 := sdk.CUAddress([]byte("output"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("atom", 10))
	var msg = MsgMultiSend{
		Inputs:  []Input{NewInput(addr1, coins)},
		Outputs: []Output{NewOutput(addr2, coins)},
	}

	// TODO some failures for bad result
	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "multisend")
}

func TestInputValidation(t *testing.T) {
	addr1 := sdk.CUAddress([]byte{1, 2})
	addr2 := sdk.CUAddress([]byte{7, 8})
	someCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	multiCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 20))

	var emptyAddr sdk.CUAddress
	emptyCoins := sdk.NewCoins()
	emptyCoins2 := sdk.NewCoins(sdk.NewInt64Coin("eth", 0))
	someEmptyCoins := sdk.Coins{sdk.NewInt64Coin("eth", 10), sdk.NewInt64Coin("atom", 0)}
	unsortedCoins := sdk.Coins{sdk.NewInt64Coin("eth", 1), sdk.NewInt64Coin("atom", 1)}

	cases := []struct {
		valid bool
		txIn  Input
	}{
		// auth works with different apps
		{true, NewInput(addr1, someCoins)},
		{true, NewInput(addr2, someCoins)},
		{true, NewInput(addr2, multiCoins)},

		{false, NewInput(emptyAddr, someCoins)},  // empty address
		{false, NewInput(addr1, emptyCoins)},     // invalid coins
		{false, NewInput(addr1, emptyCoins2)},    // invalid coins
		{false, NewInput(addr1, someEmptyCoins)}, // invalid coins
		{false, NewInput(addr1, unsortedCoins)},  // unsorted coins
	}

	for i, tc := range cases {
		err := tc.txIn.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestOutputValidation(t *testing.T) {
	addr1 := sdk.CUAddress([]byte{1, 2})
	addr2 := sdk.CUAddress([]byte{7, 8})
	someCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	multiCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 20))

	var emptyAddr sdk.CUAddress
	emptyCoins := sdk.NewCoins()
	emptyCoins2 := sdk.NewCoins(sdk.NewInt64Coin("eth", 0))
	someEmptyCoins := sdk.Coins{sdk.NewInt64Coin("eth", 10), sdk.NewInt64Coin("atom", 0)}
	unsortedCoins := sdk.Coins{sdk.NewInt64Coin("eth", 1), sdk.NewInt64Coin("atom", 1)}

	cases := []struct {
		valid bool
		txOut Output
	}{
		// auth works with different apps
		{true, NewOutput(addr1, someCoins)},
		{true, NewOutput(addr2, someCoins)},
		{true, NewOutput(addr2, multiCoins)},

		{false, NewOutput(emptyAddr, someCoins)},  // empty address
		{false, NewOutput(addr1, emptyCoins)},     // invalid coins
		{false, NewOutput(addr1, emptyCoins2)},    // invalid coins
		{false, NewOutput(addr1, someEmptyCoins)}, // invalid coins
		{false, NewOutput(addr1, unsortedCoins)},  // unsorted coins
	}

	for i, tc := range cases {
		err := tc.txOut.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestMsgMultiSendValidation(t *testing.T) {
	addr1 := sdk.CUAddress([]byte{1, 2})
	addr2 := sdk.CUAddress([]byte{7, 8})
	atom123 := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	atom124 := sdk.NewCoins(sdk.NewInt64Coin("atom", 124))
	eth123 := sdk.NewCoins(sdk.NewInt64Coin("eth", 123))
	atom123eth123 := sdk.NewCoins(sdk.NewInt64Coin("atom", 123), sdk.NewInt64Coin("eth", 123))

	input1 := NewInput(addr1, atom123)
	input2 := NewInput(addr1, eth123)
	output1 := NewOutput(addr2, atom123)
	output2 := NewOutput(addr2, atom124)
	outputMulti := NewOutput(addr2, atom123eth123)

	var emptyAddr sdk.CUAddress

	cases := []struct {
		valid bool
		tx    MsgMultiSend
	}{
		{false, MsgMultiSend{}},                           // no input or output
		{false, MsgMultiSend{Inputs: []Input{input1}}},    // just input
		{false, MsgMultiSend{Outputs: []Output{output1}}}, // just output
		{false, MsgMultiSend{
			Inputs:  []Input{NewInput(emptyAddr, atom123)}, // invalid input
			Outputs: []Output{output1}}},
		{false, MsgMultiSend{
			Inputs:  []Input{input1},
			Outputs: []Output{{emptyAddr, atom123}}}, // invalid output
		},
		{false, MsgMultiSend{
			Inputs:  []Input{input1},
			Outputs: []Output{output2}}, // amounts dont match
		},
		{true, MsgMultiSend{
			Inputs:  []Input{input1},
			Outputs: []Output{output1}},
		},
		{true, MsgMultiSend{
			Inputs:  []Input{input1, input2},
			Outputs: []Output{outputMulti}},
		},
	}

	for i, tc := range cases {
		err := tc.tx.ValidateBasic()
		if tc.valid {
			require.Nil(t, err, "%d: %+v", i, err)
		} else {
			require.NotNil(t, err, "%d", i)
		}
	}
}

func TestMsgMultiSendGetSignBytes(t *testing.T) {
	addr1 := sdk.CUAddress([]byte("LUSH6CPn8RXTW759FQsK"))
	addr2 := sdk.CUAddress([]byte("Ky7eQidQcEKG1nA24bdP"))
	coins := sdk.NewCoins(sdk.NewInt64Coin("atom", 10))
	var msg = MsgMultiSend{
		Inputs:  []Input{NewInput(addr1, coins)},
		Outputs: []Output{NewOutput(addr2, coins)},
	}
	res := msg.GetSignBytes()

	expected := "{\"type\":\"bhexchain/transfer/MsgMultiSend\",\"value\":{\"inputs\":[{\"address\":\"HBCT9W8oJnZA4h71QAFyj4we4bF1AACEgxVJ\",\"coins\":[{\"amount\":\"10\",\"denom\":\"atom\"}]}],\"max_height\":\"0\",\"outputs\":[{\"address\":\"HBCT4xTMyxwd2ZwbDiiG1BfNbCs1gwrwHYho\",\"coins\":[{\"amount\":\"10\",\"denom\":\"atom\"}]}]}}"
	require.Equal(t, expected, string(res))
}

func TestMsgMultiSendGetSigners(t *testing.T) {
	var msg = MsgMultiSend{
		Inputs: []Input{
			NewInput(sdk.CUAddress([]byte("input1")), nil),
			NewInput(sdk.CUAddress([]byte("input2")), nil),
			NewInput(sdk.CUAddress([]byte("input3")), nil),
		},
	}
	res := msg.GetSigners()
	// TODO: fix this !
	require.Equal(t, fmt.Sprintf("%v", res), "[696E70757431 696E70757432 696E70757433]")
}

/*
// what to do w/ this test?
func TestMsgSendSigners(t *testing.T) {
	signers := []sdk.CUAddress{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	someCoins := sdk.NewCoins(sdk.NewInt64Coin("atom", 123))
	inputs := make([]Input, len(signers))
	for i, signer := range signers {
		inputs[i] = NewInput(signer, someCoins)
	}
	tx := NewMsgSend(inputs, nil)

	require.Equal(t, signers, tx.Signers())
}
*/
