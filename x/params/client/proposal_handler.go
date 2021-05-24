package client

import (
	govclient "github.com/bluehelix-chain/bhchain/x/gov/client"
	"github.com/bluehelix-chain/bhchain/x/params/client/cli"
	"github.com/bluehelix-chain/bhchain/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
