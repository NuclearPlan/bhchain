package types

import (
	"fmt"
	"time"

	sdk "github.com/bluehelix-chain/bhchain/types"
	params "github.com/bluehelix-chain/bhchain/x/params/subspace"
)

// Parameter store key
var (
	ParamStoreKeyDepositParams       = []byte("depositparams")
	ParamStoreKeyVotingParams        = []byte("votingparams")
	ParamStoreKeyTallyParams         = []byte("tallyparams")
	ParamStoreKeyProposalOnlyKeyNode = []byte("onlykeynode")
)

// Key declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyDepositParams, DepositParams{},
		ParamStoreKeyVotingParams, VotingParams{},
		ParamStoreKeyTallyParams, TallyParams{},
		ParamStoreKeyProposalOnlyKeyNode, true,
	)
}

// Param around deposits for governance
type DepositParams struct {
	MinInitDeposit    sdk.Coins     `json:"min_init_deposit,omitempty" yaml:"min_init_deposit,omitempty"` // Minimum initial deposit for submitting  a proposal
	MinDeposit        sdk.Coins     `json:"min_deposit,omitempty" yaml:"min_deposit,omitempty"`           //  Minimum deposit for a proposal to enter voting period.
	MinDaoInitDeposit sdk.Coins     `json:"min_dao_init_deposit,omitempty" yaml:"min_dao_init_deposit,omitempty"`
	MinDaoDeposit     sdk.Coins     `json:"min_dao_deposit,omitempty" yaml:"min_dao_deposit,omitempty"`
	MaxDepositPeriod  time.Duration `json:"max_deposit_period,omitempty" yaml:"max_deposit_period,omitempty"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
}

// NewDepositParams creates a new DepositParams object
func NewDepositParams(minInitDeposit, minDeposit, minDaoInitDeposit, minDaoDeposit sdk.Coins, maxDepositPeriod time.Duration) DepositParams {
	return DepositParams{
		MinInitDeposit:    minInitDeposit,
		MinDeposit:        minDeposit,
		MinDaoInitDeposit: minDaoInitDeposit,
		MinDaoDeposit:     minDaoDeposit,
		MaxDepositPeriod:  maxDepositPeriod,
	}
}

func (dp DepositParams) String() string {
	return fmt.Sprintf(`Deposit Params:
  Min Init Deposit:   %s
  Min Deposit:        %s
  Min Dao Init Deposit: %s
  Min Dao Deposit:    %s
  Max Deposit Period: %s`, dp.MinInitDeposit, dp.MinDeposit, dp.MinDaoInitDeposit, dp.MinDaoDeposit, dp.MaxDepositPeriod)
}

// Checks equality of DepositParams
func (dp DepositParams) Equal(dp2 DepositParams) bool {
	return dp.MinInitDeposit.IsEqual(dp2.MinInitDeposit) && dp.MinDeposit.IsEqual(dp2.MinDeposit) && dp.MaxDepositPeriod == dp2.MaxDepositPeriod
}

// Param around Tallying votes in governance
type TallyParams struct {
	Quorum    sdk.Dec `json:"quorum,omitempty" yaml:"quorum,omitempty"`       //  Minimum percentage of total stake needed to vote for a result to be considered valid
	Threshold sdk.Dec `json:"threshold,omitempty" yaml:"threshold,omitempty"` //  Minimum proportion of Yes votes for proposal to pass. Initial value: 0.5
	Veto      sdk.Dec `json:"veto,omitempty" yaml:"veto,omitempty"`           //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
}

// NewTallyParams creates a new TallyParams object
func NewTallyParams(quorum, threshold, veto sdk.Dec) TallyParams {
	return TallyParams{
		Quorum:    quorum,
		Threshold: threshold,
		Veto:      veto,
	}
}

func (tp TallyParams) String() string {
	return fmt.Sprintf(`Tally Params:
  Quorum:             %s
  Threshold:          %s
  Veto:               %s`,
		tp.Quorum, tp.Threshold, tp.Veto)
}

// Param around Voting in governance
type VotingParams struct {
	VotingPeriod time.Duration `json:"voting_period,omitempty" yaml:"voting_period,omitempty"` //  Length of the voting period.
	MinVoteTime  int64         `json:"min_vote_time,omitempty" yaml:"min_vote_time,omitempty"`
}

// NewVotingParams creates a new VotingParams object
func NewVotingParams(votingPeriod time.Duration, minVoteTime int64) VotingParams {
	return VotingParams{
		VotingPeriod: votingPeriod,
		MinVoteTime:  minVoteTime,
	}
}

func (vp VotingParams) String() string {
	return fmt.Sprintf(`Voting Params:
  Voting Period:      %s
  Min Vote Time:      %d`, vp.VotingPeriod, vp.MinVoteTime)
}

// Params returns all of the governance params
type Params struct {
	ProposalOnlyKeyNode bool          `json:"proposal_only_key_node" yaml:"proposal_only_key_node"`
	VotingParams        VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams         TallyParams   `json:"tally_params" yaml:"tally_params"`
	DepositParams       DepositParams `json:"deposit_params" yaml:"deposit_parmas"`
}

func (gp Params) String() string {
	return gp.VotingParams.String() + "\n" +
		gp.TallyParams.String() + "\n" + gp.DepositParams.String()
}

func NewParams(vp VotingParams, tp TallyParams, dp DepositParams, proposalOnlyKeyNode bool) Params {
	return Params{
		ProposalOnlyKeyNode: proposalOnlyKeyNode,
		VotingParams:        vp,
		DepositParams:       dp,
		TallyParams:         tp,
	}
}
