package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bluehelix-chain/bhchain/client"
	"github.com/bluehelix-chain/bhchain/client/context"
	"github.com/bluehelix-chain/bhchain/codec"
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/custodianunit"
	"github.com/bluehelix-chain/bhchain/x/custodianunit/client/utils"
	"github.com/bluehelix-chain/bhchain/x/transfer/types"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flagOrderID = "order-id"

const (
	flagOutfile = "output-document"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Transfer transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		SendTxCmd(cdc),
		MultiSendTxCmd(cdc),
		CancelWithDrawalCmd(cdc),
		DepositCmd(cdc),
		WithDrawalCmd(cdc),
		RecollectCmd(cdc),
	)
	return txCmd
}

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [from_key_or_address] [to_address] [amount]",
		Short: "Create and sign a send tx",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := custodianunit.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.CUAddressFromBase58(args[1])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgSend(cliCtx.GetFromAddress(), to, coins)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func MultiSendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisend [address] [coin] [exchange_address] [exchange_coin] [height]",
		Short: "Create a multisend tx",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := custodianunit.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)
			cliCtx.GenerateOnly = true
			// parse coins trying to be sent
			addr, err := sdk.CUAddressFromBase58(args[0])
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			exchangeAddr, err := sdk.CUAddressFromBase58(args[2])
			if err != nil {
				return err
			}

			exchangeCoins, err := sdk.ParseCoins(args[3])
			if err != nil {
				return err
			}

			blockheight, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("Invalid blockheight:%v", args[4])
			}

			input := types.Input{
				Address: addr,
				Coins:   coins,
			}

			exchangeInput := types.Input{
				Address: exchangeAddr,
				Coins:   exchangeCoins,
			}

			output := types.Output{
				Address: exchangeAddr,
				Coins:   coins,
			}

			exchangeOutPut := types.Output{
				Address: addr,
				Coins:   exchangeCoins,
			}

			var inputs []types.Input
			var outputs []types.Output
			inputs = append(inputs, input)
			inputs = append(inputs, exchangeInput)

			outputs = append(outputs, exchangeOutPut)
			outputs = append(outputs, output)

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgMultiSend(inputs, outputs, blockheight)
			filePath := viper.GetString(flagOutfile)
			if filePath == "" {
				return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			}
			return utils.GenerateAndSaveMsgs(cliCtx, txBldr, []sdk.Msg{msg}, filePath)
		},
	}

	cmd = client.PostCommands(cmd)[0]
	cmd.Flags().String(flagOutfile, "", "The document will be written to the given file instead of STDOUT")

	return cmd
}

func DepositCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		//  fromCU, toCU, toAddr, symbol, hash, orderID, memo string, amount sdk.Int, index uint16, height uint64
		Use:   "deposit [from_key_or_address] [toCU_address] [to_address] [coin] [txhash] [index] [memo]",
		Short: "Deposit asset to HBTC Chain",
		Long: `  Deposit asset to HBTC Chain, and HBTC Chain will check it through the real chain.
  Example: bhcli tx transfer deposit alice BHPSfYjrgEgM97gpCWRd2UStvRVRqFgw6mQ  0x2e9a512fc6fea120e567ed5faef1440e4f66b5ff 1000000000000000000eth 0x1b5894be4f66eb75a63b5010db4a0a5cfe0b589956b74086bf64585939da1659 0 memo 5377064  --chain-id bhexchain`,
		Args: cobra.MinimumNArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := custodianunit.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.CUAddressFromBase58(args[1])
			if err != nil {
				return err
			}
			extAddress := args[2]
			coins, err := sdk.ParseCoins(args[3])
			if err != nil {
				return err
			}
			txHash := args[4]
			index, err := strconv.ParseUint(args[5], 10, 16)
			if err != nil {
				return fmt.Errorf("Invalid index:%v", args[5])
			}
			memo := args[6]

			orderID := viper.GetString(flagOrderID)
			if len(orderID) == 0 {
				orderID = uuid.NewV4().String()
			}
			msg := types.NewMsgDeposit(cliCtx.GetFromAddress(), to, sdk.Symbol(coins[0].Denom), extAddress, txHash, orderID, memo, coins[0].Amount, uint16(index))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd = client.PostCommands(cmd)[0]
	cmd.Flags().String(flagOrderID, "", "order ID of deposit is a uuid string. e.g. 'fc9ffd98-c99f-4a7c-b3ab-a517fed807c4'")

	return cmd
}

func CancelWithDrawalCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		//  fromCU, toCU, toAddr, symbol, hash, orderID, memo string, amount sdk.Int, index uint16, height uint64
		Use:   "cancelwithdrawal [from_address] [orderid]",
		Short: "cancel withdrawal tx",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := custodianunit.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			msg := types.NewMsgCancelWithdrawal(cliCtx.GetFromAddress().String(), args[1])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func WithDrawalCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		//fromCU, toAddr, symbol, orderID string, amount, gasFee sdk.Int)
		Use:   "withdrawal [from_key_or_address] [to_address] [coin] [gas]",
		Short: "withdrawal to sign tx for withdrawal asset from sepecified CU to withdrawal address",
		Long: `  withdrawal to sign tx for withdrawal asset from sepecified CU to withdrawal address
  Example: bhcli tx transfer withdrawal alice 0x2e9a512fc6fea120e567ed5faef1440e4f66b5ff 1000000000000000000eth 1000000000000000 --chain-id bhchain`,

		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := custodianunit.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			toExtAddress := args[1]
			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}
			gas, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("Invalid gas:%v", args[3])
			}

			orderID := uuid.NewV4()
			msg := types.NewMsgWithdrawal(cliCtx.GetFromAddress().String(), toExtAddress, coins[0].Denom, orderID.String(), coins[0].Amount, gas)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			fmt.Println("orderID:", orderID)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func RecollectCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recollect [orderIDs]",
		Short: "recollect a list of order",
		Long: `  recollect a list of order
  Example: bhcli tx transfer recollect 17e53eec-e01e-4fe0-81fc-a77f581ff18c,c51dc506-ba6b-4d91-96cc-2666cd4a60ec --chain-id bhchain --from node0`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := custodianunit.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			orders := strings.Split(args[0], ",")
			msg := types.NewMsgRecollect(cliCtx.GetFromAddress().String(), orders)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
