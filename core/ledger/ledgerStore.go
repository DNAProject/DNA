package ledger

import (
	. "DNA/common"
	. "DNA/core/asset"
	tx "DNA/core/transaction"
)

// ILedgerStore provides func with store package.
type ILedgerStore interface {
	//TODO: define the state store func
	SaveBlock(b *Block,ledger *Ledger) error
	GetBlock(hash Uint256) (*Block, error)
	GetBlockHash(height uint32) (Uint256, error)
	InitLedgerStore(ledger *Ledger) error
	IsDoubleSpend(tx *tx.Transaction) bool

	//SaveHeader(header *Header,ledger *Ledger) error
	AddHeaders(headers []Header, ledger *Ledger) error
	GetHeader(hash Uint256) (*Header, error)

	GetTransaction(hash Uint256) (*tx.Transaction,error)

	SaveAsset(assetid Uint256,asset *Asset) error
	GetAsset(hash Uint256) (*Asset, error)

	GetCurrentBlockHash() Uint256
	GetCurrentHeaderHash() Uint256
	GetHeaderHeight() uint32
	GetHeight() uint32

	Close() error

	InitLevelDBStoreWithGenesisBlock( genesisblock * Block  ) (uint32, error)

	GetQuantityIssued(assetid Uint256) (Fixed64, error)

	GetUnspent(txid Uint256, index uint16) (*tx.TxOutput, error)
	ContainsUnspent(txid Uint256, index uint16) (bool,error)
}
