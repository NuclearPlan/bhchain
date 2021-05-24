package client

import (
	govclient "github.com/bluehelix-chain/bhchain/x/gov/client"
	"github.com/bluehelix-chain/bhchain/x/mapping/client/cli"
)

var AddMappingProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitAddMappingProposal, nil)
var SwitchMappingProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitSwitchMappingProposal, nil)
