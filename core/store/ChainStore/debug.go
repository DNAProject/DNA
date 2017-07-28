package ChainStore

import (
	// . "DNA/common"
	"expvar"
	"fmt"
	// "time"
)

type storeStatus struct {
	CurrHeaderHeight uint32
	CurrBlockHeight  uint32
	HeaderCache      map[uint32]string
	BlockCache       map[uint32]string
	HeaderIndex      map[uint32]string
}

func expvarStore(store *ChainStore) func() interface{} {

	return func() interface{} {
		headerHeight := store.GetHeaderHeight()
		blockHeight := store.GetHeight()
		store.mu.RLock()
		defer store.mu.RUnlock()

		ss := storeStatus{
			HeaderIndex:      make(map[uint32]string, len(store.headerIndex)),
			CurrHeaderHeight: headerHeight,
			CurrBlockHeight:  blockHeight,
			HeaderCache:      make(map[uint32]string, len(store.headerCache)),
			BlockCache:       make(map[uint32]string, len(store.blockCache)),
		}

		for k, hash := range store.headerIndex {
			ss.HeaderIndex[k] = fmt.Sprintf("%x", hash)
		}

		for k, header := range store.headerCache {
			ss.HeaderCache[header.Blockdata.Height] = fmt.Sprintf("%x", k)
		}
		for k, block := range store.blockCache {
			ss.BlockCache[block.Blockdata.Height] = fmt.Sprintf("%x", k)
		}

		return ss
	}

}

func ExportStoreStatus(store *ChainStore) {

	expvar.Publish("dna_store", expvar.Func(expvarStore(store)))
}
