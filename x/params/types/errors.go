package types

import (
	"fmt"

	sdk "github.com/bluehelix-chain/bhchain/types"
)

// Param module codespace constants
const (
	DefaultCodespace sdk.CodespaceType = "params"

	CodeUnknownSubspace     sdk.CodeType = 1
	CodeSettingParameter    sdk.CodeType = 2
	CodeEmptyData           sdk.CodeType = 3
	CodeMixedSubspace       sdk.CodeType = 4
	CodeInvalidChangeParams sdk.CodeType = 5
)

// ErrUnknownSubspace returns an unknown subspace error.
func ErrUnknownSubspace(codespace sdk.CodespaceType, space string) sdk.Error {
	return sdk.NewError(codespace, CodeUnknownSubspace, fmt.Sprintf("unknown subspace %s", space))
}

// ErrSettingParameter returns an error for failing to set a parameter.
func ErrSettingParameter(codespace sdk.CodespaceType, key, subkey, value, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeSettingParameter, fmt.Sprintf("error setting parameter %s on %s (%s): %s", value, key, subkey, msg))
}

// ErrEmptyChanges returns an error for empty parameter changes.
func ErrEmptyChanges(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyData, "submitted parameter changes are empty")
}

// ErrEmptySubspace returns an error for an empty subspace.
func ErrEmptySubspace(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyData, "parameter subspace is empty")
}

// ErrEmptyKey returns an error for when an empty key is given.
func ErrEmptyKey(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyData, "parameter key is empty")
}

// ErrEmptyValue returns an error for when an empty key is given.
func ErrEmptyValue(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyData, "parameter value is empty")
}

func ErrMixedSubspace(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMixedSubspace, "subspace is mixed")
}

// ErrInvalidChange returns an error for invalid parameter changes.
func ErrInvalidChange(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidChangeParams, msg)
}
