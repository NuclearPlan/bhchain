// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	types "github.com/bluehelix-chain/bhchain/types"
	mock "github.com/stretchr/testify/mock"
)

// StakingKeeper is an autogenerated mock type for the StakingKeeper type
type StakingKeeper struct {
	mock.Mock
}

// GetCurrentEpoch provides a mock function with given fields: ctx
func (_m *StakingKeeper) GetCurrentEpoch(ctx types.Context) types.Epoch {
	ret := _m.Called(ctx)

	var r0 types.Epoch
	if rf, ok := ret.Get(0).(func(types.Context) types.Epoch); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(types.Epoch)
	}

	return r0
}

// GetEpochByHeight provides a mock function with given fields: ctx, height
func (_m *StakingKeeper) GetEpochByHeight(ctx types.Context, height uint64) types.Epoch {
	ret := _m.Called(ctx, height)

	var r0 types.Epoch
	if rf, ok := ret.Get(0).(func(types.Context, uint64) types.Epoch); ok {
		r0 = rf(ctx, height)
	} else {
		r0 = ret.Get(0).(types.Epoch)
	}

	return r0
}

// JailByOperator provides a mock function with given fields: ctx, operator
func (_m *StakingKeeper) JailByOperator(ctx types.Context, operator types.ValAddress) {
	_m.Called(ctx, operator)
}

// SlashByOperator provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *StakingKeeper) SlashByOperator(_a0 types.Context, _a1 types.ValAddress, _a2 int64, _a3 types.Dec) {
	_m.Called(_a0, _a1, _a2, _a3)
}
