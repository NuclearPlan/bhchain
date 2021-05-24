package types

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/bluehelix-chain/bhchain/codec"
	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/params"
	tokentypes "github.com/bluehelix-chain/bhchain/x/token/types"
)

var (
	DefaultMinimumLiquidity            = sdk.NewInt(1000)
	DefaultLimitSwapMatchingGas        = sdk.NewUint(50000)
	DefaultMaxFeeRate                  = sdk.NewDecWithPrec(1, 1)  // 0.1
	DefaultLpRewardRate                = sdk.NewDecWithPrec(13, 4) // 0.0013
	DefaultRefererTransactionBonusRate = sdk.NewDecWithPrec(1, 4)  // 0.0001
	DefaultRepurchaseRate              = sdk.NewDecWithPrec(4, 4)  // 0.0004
	DefaultRepurchaseDuration          = int64(1000)
	DefaultRepurchaseToken             = sdk.NativeToken
	DefaultRepurchaseRoutingToken      = tokentypes.CalSymbol("usdt", sdk.NativeToken).String()
	DefaultLPMiningWeights             = []*MiningWeight{}
	DefaultTradeMiningWeights          = []*MiningWeight{}
	DefaultLPMiningRewardRate          = sdk.NewDecWithPrec(9, 1) // 0.9
)

var (
	KeyMinimumLiquidity            = []byte("MinimumLiquidity")
	KeyLimitSwapMatchingGas        = []byte("LimitSwapMatchingGas")
	KeyMaxFeeRate                  = []byte("MaxFeeRate")
	KeyLpRewardRate                = []byte("LpRewardRate")
	KeyRepurchaseRate              = []byte("RepurchaseRate")
	KeyRefererTransactionBonusRate = []byte("RefererTransactionBonusRate")
	KeyRepurchaseDuration          = []byte("RepurchaseDuration")
	KeyRepurchaseToken             = []byte("RepurchaseToken")
	KeyRepurchaseRoutingToken      = []byte("RepurchaseRoutingToken")
	KeyLPMiningWeights             = []byte("LPMiningWeights")
	KeyTradeMiningWeights          = []byte("TradeMiningWeights")
	KeyLPMiningRewardRate          = []byte("LPMiningRewardRate")
)

type MiningWeight struct {
	DexID  uint32     `json:"dex_id"`
	TokenA sdk.Symbol `json:"token_a"`
	TokenB sdk.Symbol `json:"token_b"`
	Weight sdk.Int    `json:"weight"`
}

func NewMiningWeight(dexID uint32, tokenA, tokenB sdk.Symbol, weight sdk.Int) *MiningWeight {
	return &MiningWeight{
		DexID:  dexID,
		TokenA: tokenA,
		TokenB: tokenB,
		Weight: weight,
	}
}

var _ params.ParamSet = (*Params)(nil)

// Params defines the high level settings for staking
type Params struct {
	MinimumLiquidity            sdk.Int         `json:"minimum_liquidity"`
	LimitSwapMatchingGas        sdk.Uint        `json:"limit_swap_matching_gas"`
	MaxFeeRate                  sdk.Dec         `json:"max_fee_rate"`
	LpRewardRate                sdk.Dec         `json:"lp_reward_rate"`
	RepurchaseRate              sdk.Dec         `json:"repurchase_rate"`
	RefererTransactionBonusRate sdk.Dec         `json:"referer_transaction_bonus_rate"`
	RepurchaseDuration          int64           `json:"repurchase_duration"`
	LPMiningWeights             []*MiningWeight `json:"lp_mining_weights"`
	TradeMiningWeights          []*MiningWeight `json:"trade_mining_weights"`
	LPMiningRewardRate          sdk.Dec         `json:"lp_mining_reward_rate"`
	RepurchaseToken             string          `json:"repurchase_token"`
	RepurchaseRoutingToken      string          `json:"repurchase_routing_token"`
}

// NewParams creates a new Params instance
func NewParams(minLiquidity sdk.Int, limitSwapMatchingGas sdk.Uint, maxFeeRate, lpRewardRate, repurchaseRate, refererTransactionBonusRate sdk.Dec,
	repurchaseDuration int64, lpMiningWeights, tradeMiningWeights []*MiningWeight,
	lpMiningRewardRate sdk.Dec, repurchaseToken, repurchaseRoutingToken string) Params {
	return Params{
		MinimumLiquidity:            minLiquidity,
		LimitSwapMatchingGas:        limitSwapMatchingGas,
		MaxFeeRate:                  maxFeeRate,
		LpRewardRate:                lpRewardRate,
		RepurchaseRate:              repurchaseRate,
		RefererTransactionBonusRate: refererTransactionBonusRate,
		RepurchaseDuration:          repurchaseDuration,
		LPMiningWeights:             lpMiningWeights,
		TradeMiningWeights:          tradeMiningWeights,
		LPMiningRewardRate:          lpMiningRewardRate,
		RepurchaseToken:             repurchaseToken,
		RepurchaseRoutingToken:      repurchaseRoutingToken,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMinimumLiquidity, &p.MinimumLiquidity},
		{KeyLimitSwapMatchingGas, &p.LimitSwapMatchingGas},
		{KeyMaxFeeRate, &p.MaxFeeRate},
		{KeyLpRewardRate, &p.LpRewardRate},
		{KeyRepurchaseRate, &p.RepurchaseRate},
		{KeyRefererTransactionBonusRate, &p.RefererTransactionBonusRate},
		{KeyRepurchaseDuration, &p.RepurchaseDuration},
		{KeyLPMiningWeights, &p.LPMiningWeights},
		{KeyTradeMiningWeights, &p.TradeMiningWeights},
		{KeyLPMiningRewardRate, &p.LPMiningRewardRate},
		{KeyRepurchaseToken, &p.RepurchaseToken},
		{KeyRepurchaseRoutingToken, &p.RepurchaseRoutingToken},
	}
}

func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultMinimumLiquidity, DefaultLimitSwapMatchingGas, DefaultMaxFeeRate, DefaultLpRewardRate,
		DefaultRepurchaseRate, DefaultRefererTransactionBonusRate,
		DefaultRepurchaseDuration, DefaultLPMiningWeights,
		DefaultTradeMiningWeights, DefaultLPMiningRewardRate,
		DefaultRepurchaseToken, DefaultRepurchaseRoutingToken)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  MinimumLiquidity: %s
  LimitSwapMatchingGas: %s
  MaxFeeRate: %s
  LpRewardRate: %s
  RepurchaseRate: %s
  RefererTransactionBonusRate: %s
  RepurchaseDuration: %d
  RepurchaseToken: %s
  RepurchaseRoutingToken: %s
  LPMiningWeights: %v
  TradeMiningWeights: %v
  LPMiningRewardRate: %v`,
		p.MinimumLiquidity.String(), p.LimitSwapMatchingGas.String(), p.MaxFeeRate.String(),
		p.LpRewardRate.String(), p.RepurchaseRate.String(), p.RefererTransactionBonusRate.String(),
		p.RepurchaseDuration, p.RepurchaseToken, p.RepurchaseRoutingToken, p.LPMiningWeights,
		p.TradeMiningWeights, p.LPMiningRewardRate)
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.Codec, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		panic(err)
	}
	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.Codec, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &params)
	if err != nil {
		return
	}
	return
}

// validate a set of params
func (p Params) Validate() error {
	if !p.MinimumLiquidity.IsPositive() {
		return errors.New("minimun liquidity should be positive")
	}
	if p.MaxFeeRate.IsNegative() || p.MaxFeeRate.GT(sdk.OneDec()) {
		return errors.New("max fee rate must be between 0 to 1")
	}
	if p.LpRewardRate.IsNegative() {
		return errors.New("fee rate cannot be negative")
	}
	if p.RepurchaseRate.IsNegative() {
		return errors.New("repurchase rate cannot be negative")
	}
	if p.RefererTransactionBonusRate.IsNegative() {
		return errors.New("referer transaction bonus rate cannot be negative")
	}
	if p.LpRewardRate.Add(p.RepurchaseRate).Add(p.RefererTransactionBonusRate).GT(p.MaxFeeRate) {
		return errors.New("sum of fee rate must be less than max fee rate")
	}
	if p.RepurchaseDuration <= 0 {
		return errors.New("repurchase duration should be positive")
	}

	if err := validateMiningWeight(p.LPMiningWeights); err != nil {
		return err
	}
	if err := validateMiningWeight(p.TradeMiningWeights); err != nil {
		return err
	}
	if p.LPMiningRewardRate.IsNegative() {
		return errors.New("lp mining reward rate cannot be negative")
	}

	if p.LPMiningRewardRate.GT(sdk.OneDec()) {
		return errors.New("lp mining reward rates cannot be larger than 1.0")
	}

	return nil
}

func validateMiningWeight(miningWeights []*MiningWeight) error {
	exists := make(map[string]bool)
	for _, w := range miningWeights {
		if w.TokenA == w.TokenB {
			return errors.New("tokenA and tokenB can be the same")
		}
		tokenA, tokenB := w.TokenA, w.TokenB
		if tokenA > tokenB {
			tokenA, tokenB = tokenB, tokenA
		}
		pair := fmt.Sprintf("%s-%s", tokenA, tokenB)
		if exists[pair] {
			return fmt.Errorf("%s is duplicated", pair)
		}
		exists[pair] = true
		if !w.Weight.IsPositive() {
			return errors.New("weight should be positive")
		}
	}
	return nil
}
