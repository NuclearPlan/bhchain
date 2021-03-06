package internal

import (
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/evidence/exported"
)

type StakingKeeper interface {
	IsActiveKeyNode(ctx sdk.Context, addr sdk.CUAddress) (bool, int)
}

type EvidenceKeeper interface {
	VoteWithCustomBox(ctx sdk.Context, voteID string, voter sdk.CUAddress, vote exported.Vote, height uint64, newVoteBox exported.NewVoteBox) (bool, bool, []*exported.VoteItem)
}
