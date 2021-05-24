package client

import (
	govclient "github.com/bluehelix-chain/bhchain/x/gov/client"
	"github.com/bluehelix-chain/bhchain/x/token/client/cli"
	"github.com/bluehelix-chain/bhchain/x/token/client/rest"
)

// param change proposal handler
var (
	AddTokenProposalHandler          = govclient.NewProposalHandler(cli.GetCmdAddTokenProposal, rest.AddTokenProposalRESTHandler)
	TokenParamsChangeProposalHandler = govclient.NewProposalHandler(cli.GetCmdTokenParamsChangeProposal, rest.TokenParamsChangeProposalRESTHandler)
)
