package smartcontract

import (
	"DNA/core/ledger"
	"DNA/common/log"
)

type CachedCodeTable struct {
}

func NewCachedCodeTable() *CachedCodeTable {
	var cachedCodeTable CachedCodeTable
	return &cachedCodeTable
}

func(c *CachedCodeTable) GetCode(codeHash []byte) ([]byte) {
	code, err := ledger.DefaultLedger.Store.GetContract(codeHash)
	if err != nil {
		log.Error("Get Contract Error:", err)
		return []byte{}
	}
	return code
}