package main

// DONTCOVER

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmconfig "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/bluehelix-chain/bhchain/client"
	"github.com/bluehelix-chain/bhchain/client/keys"
	"github.com/bluehelix-chain/bhchain/codec"
	"github.com/bluehelix-chain/bhchain/server"
	srvconfig "github.com/bluehelix-chain/bhchain/server/config"
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/types/module"
	"github.com/bluehelix-chain/bhchain/x/custodianunit"
	"github.com/bluehelix-chain/bhchain/x/genaccounts"
	"github.com/bluehelix-chain/bhchain/x/genutil"
	genutiltypes "github.com/bluehelix-chain/bhchain/x/genutil/types"
	"github.com/bluehelix-chain/bhchain/x/staking"
	stakingtypes "github.com/bluehelix-chain/bhchain/x/staking/types"
)

var (
	flagNodeDirPrefix     = "node-dir-prefix"
	flagNumValidators     = "v"
	flagNumKeyNodes       = "key-nodes"
	flagOutputDir         = "output-dir"
	flagNodeDaemonHome    = "node-daemon-home"
	flagNodeCLIHome       = "node-cli-home"
	flagStartingIPAddress = "starting-ip-address"
	flagSameIPAddress     = "same-ip-address"
)

// get cmd to initialize all files for tendermint testnet and application
func testnetCmd(ctx *server.Context, cdc *codec.Codec,
	mbm module.BasicManager, genAccIterator genutiltypes.GenesisCUsIterator,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a hbtc chain testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	bhchain testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			config := ctx.Config

			outputDir := viper.GetString(flagOutputDir)
			chainID := viper.GetString(client.FlagChainID)
			minGasPrices := viper.GetString(server.FlagMinGasPrices)
			nodeDirPrefix := viper.GetString(flagNodeDirPrefix)
			nodeDaemonHome := viper.GetString(flagNodeDaemonHome)
			nodeCLIHome := viper.GetString(flagNodeCLIHome)
			startingIPAddress := viper.GetString(flagStartingIPAddress)
			numValidators := viper.GetInt(flagNumValidators)
			numKeyNodes := viper.GetInt(flagNumKeyNodes)
			if numKeyNodes <= 0 || numKeyNodes > numValidators {
				numKeyNodes = numValidators
			}
			if viper.GetBool(flagSameIPAddress) {
				return InitTestnetSameIPDiffPort(cmd, config, cdc, mbm, genAccIterator, outputDir, chainID,
					minGasPrices, nodeDirPrefix, nodeDaemonHome, nodeCLIHome, startingIPAddress, numValidators, numKeyNodes)
			} else {
				return InitTestnet(cmd, config, cdc, mbm, genAccIterator, outputDir, chainID,
					minGasPrices, nodeDirPrefix, nodeDaemonHome, nodeCLIHome, startingIPAddress, numValidators, numKeyNodes)
			}
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with")
	cmd.Flags().Int(flagNumKeyNodes, 0,
		"Number of key nodes to initialize the testnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./testnet",
		"Directory to store initialization data for the testnet")
	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "bhexchain",
		"Home directory of the node's daemon configuration")
	cmd.Flags().String(flagNodeCLIHome, "bhcli",
		"Home directory of the node's cli configuration")
	cmd.Flags().String(flagStartingIPAddress, "192.168.0.2",
		"Starting IP address (192.168.0.2 results in persistent peers list ID0@192.168.0.2:46656, ID1@192.168.0.3:46656, ...)")
	cmd.Flags().String(
		client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(
		server.FlagMinGasPrices, fmt.Sprintf("1000000000%s", sdk.DefaultBondDenom),
		"Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.001stake)")
	cmd.Flags().Bool(flagSameIPAddress, false, "All nodes use same ip with different port")
	return cmd
}

const nodeDirPerm = 0755

// Initialize the testnet
func InitTestnet(cmd *cobra.Command, config *tmconfig.Config, cdc *codec.Codec,
	mbm module.BasicManager, genAccIterator genutiltypes.GenesisCUsIterator,
	outputDir, chainID, minGasPrices, nodeDirPrefix, nodeDaemonHome,
	nodeCLIHome, startingIPAddress string, numValidators, numKeyNodes int) error {

	if chainID == "" {
		chainID = "chain-" + cmn.RandStr(6)
	}

	monikers := make([]string, numValidators)
	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]crypto.PubKey, numValidators)

	bhConfig := srvconfig.DefaultConfig()
	bhConfig.MinGasPrices = minGasPrices

	var (
		accs     []genaccounts.GenesisCU
		genFiles []string
		keyNodes []sdk.CUAddress
	)

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		clientDir := filepath.Join(outputDir, nodeDirName, nodeCLIHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")

		config.SetRoot(nodeDir)
		config.RPC.ListenAddress = "tcp://0.0.0.0:26657"

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		if err := os.MkdirAll(clientDir, nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		monikers = append(monikers, nodeDirName)
		config.Moniker = nodeDirName

		ip := ""
		var err error
		if startingIPAddress == "0.0.0.0" {
			buf := bufio.NewReader(cmd.InOrStdin())
			promt := fmt.Sprintf(
				"ip for validator '%s' (default %s):", nodeDirName, "0.0.0.0",
			)
			ip, err = client.GetString(promt, buf)
		} else {
			ip, err = getIP(i, startingIPAddress)
			if err != nil {
				_ = os.RemoveAll(outputDir)
				return err
			}
		}

		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(config)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		genFiles = append(genFiles, config.GenesisFile())

		buf := bufio.NewReader(cmd.InOrStdin())
		prompt := fmt.Sprintf(
			"Password for account '%s' (default %s):", nodeDirName, client.DefaultKeyPass,
		)

		keyPass, err := client.GetPassword(prompt, buf)
		if err != nil && keyPass != "" {
			// An error was returned that either failed to read the password from
			// STDIN or the given password is not empty but failed to meet minimum
			// length requirements.
			return err
		}

		if keyPass == "" {
			keyPass = client.DefaultKeyPass
		}
		addr, secret, err := server.GenerateSaveCoinKey(clientDir, nodeDirName, keyPass, true)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint); err != nil {
			return err
		}

		accStakingTokens := sdk.TokensFromConsensusPower(1000000)
		accs = append(accs, genaccounts.GenesisCU{
			Address: addr,
			Coins: sdk.Coins{
				sdk.NewCoin(sdk.DefaultBondDenom, accStakingTokens),
			},
		})

		valTokens := sdk.TokensFromConsensusPower(500000)
		msg := staking.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			staking.NewDescription(nodeDirName, "", "", ""),
			staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			stakingtypes.DefaultMinValidatorDelegation,
		)
		kb, err := keys.NewKeyBaseFromDir(clientDir)
		if err != nil {
			return err
		}
		tx := custodianunit.NewStdTx([]sdk.Msg{msg}, custodianunit.StdFee{}, []custodianunit.StdSignature{}, memo)
		txBldr := custodianunit.NewTxBuilderFromCLI().WithChainID(chainID).WithMemo(memo).WithKeybase(kb)

		signedTx, err := txBldr.SignStdTx(nodeDirName, client.DefaultKeyPass, tx, false)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		txBytes, err := cdc.MarshalJSON(signedTx)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		// gather gentxs folder
		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBytes); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		bhConfigFilePath := filepath.Join(nodeDir, "config/bhexchain.toml")
		srvconfig.WriteConfigFile(bhConfigFilePath, bhConfig)

		if i < numKeyNodes {
			keyNodes = append(keyNodes, addr)
		}
	}

	if err := initGenFiles(cdc, mbm, chainID, accs, genFiles, numValidators, keyNodes); err != nil {
		return err
	}

	err := collectGenFiles(
		cdc, config, chainID, monikers, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genAccIterator,
	)
	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil

}

// Initialize the testnet
func InitTestnetSameIPDiffPort(cmd *cobra.Command, config *tmconfig.Config, cdc *codec.Codec,
	mbm module.BasicManager, genAccIterator genutiltypes.GenesisCUsIterator,
	outputDir, chainID, minGasPrices, nodeDirPrefix, nodeDaemonHome,
	nodeCLIHome, startingIPAddress string, numValidators, numKeyNodes int) error {

	if chainID == "" {
		chainID = "chain-" + cmn.RandStr(6)
	}

	monikers := make([]string, numValidators)
	nodeIDs := make([]string, numValidators)
	valPubKeys := make([]crypto.PubKey, numValidators)

	bhConfig := srvconfig.DefaultConfig()
	bhConfig.MinGasPrices = minGasPrices

	var (
		accs     []genaccounts.GenesisCU
		genFiles []string
		keyNodes []sdk.CUAddress
	)

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		clientDir := filepath.Join(outputDir, nodeDirName, nodeCLIHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")

		port26657 := strconv.Itoa(26657 + 10*i)
		port26656 := strconv.Itoa(26656 + 10*i)
		config.SetRoot(nodeDir)
		//config.RPC.ListenAddress = "tcp://0.0.0.0:26657"
		config.RPC.ListenAddress = "tcp://127.0.0.1:" + port26657

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		if err := os.MkdirAll(clientDir, nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		monikers = append(monikers, nodeDirName)
		config.Moniker = nodeDirName

		//ip, err := getIP(i, startingIPAddress)
		//if err != nil {
		//	_ = os.RemoveAll(outputDir)
		//	return err
		//}
		ip := "127.0.0.1"
		var err error

		nodeIDs[i], valPubKeys[i], err = genutil.InitializeNodeValidatorFiles(config)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		//memo := fmt.Sprintf("%s@%s:26656", nodeIDs[i], ip)
		memo := fmt.Sprintf("%s@%s:%s", nodeIDs[i], ip, port26656)
		genFiles = append(genFiles, config.GenesisFile())

		buf := bufio.NewReader(cmd.InOrStdin())
		prompt := fmt.Sprintf(
			"Password for account '%s' (default %s):", nodeDirName, client.DefaultKeyPass,
		)

		keyPass, err := client.GetPassword(prompt, buf)
		if err != nil && keyPass != "" {
			// An error was returned that either failed to read the password from
			// STDIN or the given password is not empty but failed to meet minimum
			// length requirements.
			return err
		}

		if keyPass == "" {
			keyPass = client.DefaultKeyPass
		}

		addr, secret, err := server.GenerateSaveCoinKey(clientDir, nodeDirName, keyPass, true)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint); err != nil {
			return err
		}

		accStakingTokens := sdk.TokensFromConsensusPower(1000000)
		accs = append(accs, genaccounts.GenesisCU{
			Address: addr,
			Coins: sdk.Coins{
				sdk.NewCoin(sdk.DefaultBondDenom, accStakingTokens),
			},
		})

		valTokens := sdk.TokensFromConsensusPower(500000)
		msg := staking.NewMsgCreateValidator(
			sdk.ValAddress(addr),
			valPubKeys[i],
			sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			staking.NewDescription(nodeDirName, "", "", ""),
			staking.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			stakingtypes.DefaultMinValidatorDelegation,
		)
		kb, err := keys.NewKeyBaseFromDir(clientDir)
		if err != nil {
			return err
		}
		tx := custodianunit.NewStdTx([]sdk.Msg{msg}, custodianunit.StdFee{}, []custodianunit.StdSignature{}, memo)
		txBldr := custodianunit.NewTxBuilderFromCLI().WithChainID(chainID).WithMemo(memo).WithKeybase(kb)

		signedTx, err := txBldr.SignStdTx(nodeDirName, client.DefaultKeyPass, tx, false)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		txBytes, err := cdc.MarshalJSON(signedTx)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		// gather gentxs folder
		if err := writeFile(fmt.Sprintf("%v.json", nodeDirName), gentxsDir, txBytes); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		bhConfigFilePath := filepath.Join(nodeDir, "config/bhexchain.toml")
		srvconfig.WriteConfigFile(bhConfigFilePath, bhConfig)

		if i < numKeyNodes {
			keyNodes = append(keyNodes, addr)
		}
	}

	if err := initGenFiles(cdc, mbm, chainID, accs, genFiles, numValidators, keyNodes); err != nil {
		return err
	}

	err := collectGenFiles(
		cdc, config, chainID, monikers, nodeIDs, valPubKeys, numValidators,
		outputDir, nodeDirPrefix, nodeDaemonHome, genAccIterator,
	)
	if err != nil {
		return err
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", numValidators)
	return nil
}

func initGenFiles(cdc *codec.Codec, mbm module.BasicManager, chainID string,
	accs []genaccounts.GenesisCU, genFiles []string, numValidators int, keyNodes []sdk.CUAddress) error {

	appGenState := mbm.DefaultGenesis()
	stakingState := staking.DefaultGenesisState()
	stakingState.KeyNodes = keyNodes
	appGenState[staking.ModuleName] = cdc.MustMarshalJSON(stakingState)

	// set the accounts in the genesis state
	appGenState = genaccounts.SetGenesisStateInAppState(cdc, appGenState, accs)

	appGenStateJSON, err := codec.MarshalJSONIndent(cdc, appGenState)
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}
	return nil
}

func collectGenFiles(
	cdc *codec.Codec, config *tmconfig.Config, chainID string,
	monikers, nodeIDs []string, valPubKeys []crypto.PubKey,
	numValidators int, outputDir, nodeDirPrefix, nodeDaemonHome string,
	genAccIterator genutiltypes.GenesisCUsIterator) error {

	var appState json.RawMessage
	genTime := tmtime.Now()

	for i := 0; i < numValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		gentxsDir := filepath.Join(outputDir, "gentxs")
		moniker := monikers[i]
		config.Moniker = nodeDirName

		config.SetRoot(nodeDir)

		nodeID, valPubKey := nodeIDs[i], valPubKeys[i]
		initCfg := genutil.NewInitConfig(chainID, gentxsDir, moniker, nodeID, valPubKey)

		genDoc, err := types.GenesisDocFromFile(config.GenesisFile())
		if err != nil {
			return err
		}

		nodeAppState, err := genutil.GenAppStateFromConfig(cdc, config, initCfg, *genDoc, genAccIterator)
		if err != nil {
			return err
		}

		if appState == nil {
			// set the canonical application state (they should not differ)
			appState = nodeAppState
		}

		genFile := config.GenesisFile()

		// overwrite each validator's genesis file to have a canonical genesis time
		if err := genutil.ExportGenesisFileWithTime(genFile, chainID, nil, appState, genTime); err != nil {
			return err
		}
	}

	return nil
}

func getIP(i int, startingIPAddr string) (ip string, err error) {
	if len(startingIPAddr) == 0 {
		ip, err = server.ExternalIP()
		if err != nil {
			return "", err
		}
		return ip, nil
	}
	return calculateIP(startingIPAddr, i)
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := cmn.EnsureDir(writePath, 0700)
	if err != nil {
		return err
	}

	err = cmn.WriteFile(file, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}
