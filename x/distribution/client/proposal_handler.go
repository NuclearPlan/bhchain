package client

import (
	"github.com/bluehelix-chain/bhchain/x/distribution/client/cli"
	"github.com/bluehelix-chain/bhchain/x/distribution/client/rest"
	govclient "github.com/bluehelix-chain/bhchain/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
