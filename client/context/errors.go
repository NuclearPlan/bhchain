package context

import (
	"fmt"

	sdk "github.com/bluehelix-chain/bhchain/types"
)

// ErrInvalidAccount returns a standardized error reflecting that a given
// CU address does not exist.
func ErrInvalidAccount(addr sdk.CUAddress) error {
	return fmt.Errorf(`No CU with address %s was found in the state.
Are you sure there has been a transaction involving it?`, addr)
}

// ErrVerifyCommit returns a common error reflecting that the blockchain commit at a given
// height can't be verified. The reason is that the base checkpoint of the certifier is
// newer than the given height
func ErrVerifyCommit(height int64) error {
	return fmt.Errorf(`The height of base truststore in the light client is higher than height %d. 
Can't verify blockchain proof at this height. Please set --trust-node to true and try again`, height)
}
