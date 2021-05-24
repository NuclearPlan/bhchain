package types

import (
	"fmt"

	sdk "github.com/bluehelix-chain/bhchain/types"
	"github.com/bluehelix-chain/bhchain/x/params"
)

// Parameter store keys
var (
	KeyMintDenom                         = []byte("MintDenom")
	KeyMintForStakingPerBlock            = []byte("MintForStakingPerBlock")
	KeyMintForStakingHalveDuration       = []byte("MintForStakingHalveDuration")
	KeyMintForOpenswapBonusStartBlock    = []byte("MintForOpenswapBonusStartBlock")
	KeyMintForOpenswapBonusPerBlock      = []byte("MintForOpenswapBonusPerBlock")
	KeyMintForOpenswapBonusHalveDuration = []byte("MintForOpenswapBonusHalveDuration")
)

// mint parameters
type Params struct {
	MintDenom                         string  `json:"mint_denom" yaml:"mint_denom"` // type of coin to mint
	MintForStakingPerBlock            sdk.Int `json:"mint_for_staking_per_block" yaml:"mint_for_staking_per_block"`
	MintForStakingHalveDuration       sdk.Int `json:"mint_for_staking_halve_duration" yaml:"mint_for_staking_halve_duration"`
	MintForOpenswapBonusStartBlock    sdk.Int `json:"mint_for_openswap_bonus_start_block" yaml:"mint_for_openswap_bonus_start_block"`
	MintForOpenswapBonusPerBlock      sdk.Int `json:"mint_for_openswap_bonus_per_block" yaml:"mint_for_openswap_bonus_per_block"`
	MintForOpenswapBonusHalveDuration sdk.Int `json:"mint_for_openswap_bonus_halve_duration" yaml:"mint_for_openswap_bonus_halve_duration"`
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, mintForStakingPerBlock, mintForStakingHalveDuration, mintForOpenswapBonusStartBlock,
	mintForOpenswapBonusPerBlock, mintForOpenswapBonusHalveDuration sdk.Int) Params {
	return Params{
		MintDenom:                         mintDenom,
		MintForStakingPerBlock:            mintForStakingPerBlock,
		MintForStakingHalveDuration:       mintForStakingHalveDuration,
		MintForOpenswapBonusStartBlock:    mintForOpenswapBonusStartBlock,
		MintForOpenswapBonusPerBlock:      mintForOpenswapBonusPerBlock,
		MintForOpenswapBonusHalveDuration: mintForOpenswapBonusHalveDuration,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:                         sdk.DefaultBondDenom,
		MintForStakingPerBlock:            sdk.NewIntWithDecimal(2, 16), // 0.02 hbc
		MintForStakingHalveDuration:       sdk.NewInt(5256000),
		MintForOpenswapBonusStartBlock:    sdk.ZeroInt(),
		MintForOpenswapBonusPerBlock:      sdk.ZeroInt(),
		MintForOpenswapBonusHalveDuration: sdk.ZeroInt(),
	}
}

// validate params
func ValidateParams(params Params) error {
	if params.MintDenom == "" {
		return fmt.Errorf("mint parameter MintDenom can't be an empty string")
	}

	if params.MintForStakingPerBlock.IsNegative() {
		return fmt.Errorf("mint parameter MintForStakingPerBlock can't be negative")
	}

	if params.MintForStakingHalveDuration.IsNegative() {
		return fmt.Errorf("mint parameter MintForStakingHalveDuration can't be negative")
	}

	if params.MintForOpenswapBonusStartBlock.IsNegative() {
		return fmt.Errorf("mint parameter MintForOpenswapBonusStartBlock can't be negative")
	}

	if params.MintForOpenswapBonusPerBlock.IsNegative() {
		return fmt.Errorf("mint parameter MintForOpenswapBonusPerBlock can't be negative")
	}

	if params.MintForOpenswapBonusHalveDuration.IsNegative() {
		return fmt.Errorf("mint parameter MintForOpenswapBonusHalveDuration can't be negative")
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:                          %s
  MintForStakingPerBlock:              %s
  MintForStakingHalveDuration:         %s
  MintForOpenswapBonusStartBlock:      %s
  MintForOpenswapBonusPerBlock:        %s
  MintForOpenswapBonusHalveDuration:   %s
`,
		p.MintDenom, p.MintForStakingPerBlock.String(), p.MintForStakingHalveDuration.String(),
		p.MintForOpenswapBonusStartBlock.String(), p.MintForOpenswapBonusPerBlock.String(), p.MintForOpenswapBonusHalveDuration.String())
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMintDenom, &p.MintDenom},
		{KeyMintForStakingPerBlock, &p.MintForStakingPerBlock},
		{KeyMintForStakingHalveDuration, &p.MintForStakingHalveDuration},
		{KeyMintForOpenswapBonusStartBlock, &p.MintForOpenswapBonusStartBlock},
		{KeyMintForOpenswapBonusPerBlock, &p.MintForOpenswapBonusPerBlock},
		{KeyMintForOpenswapBonusHalveDuration, &p.MintForOpenswapBonusHalveDuration},
	}
}
