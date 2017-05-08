package transaction

import (
	. "DNA/common"
)

// ILedgerStore provides func with store package.
type ITxStore interface {
	GetTransaction(hash Uint256) (*Transaction, error)
	GetQuantityIssued(AssetId Uint256) (Fixed64, error)
}
