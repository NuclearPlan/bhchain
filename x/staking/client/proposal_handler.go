package client

import (
	govclient "github.com/bluehelix-chain/bhchain/x/gov/client"
	"github.com/bluehelix-chain/bhchain/x/staking/client/cli"
)

var UpdateKeyNodesProposalHandler = govclient.NewProposalHandler(cli.NewCmdUpdateKeyNodesProposal, nil)
