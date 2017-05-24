package LevelDBStore

import (
	. "DNA/core/store"
	"github.com/syndtr/goleveldb/leveldb"
)

type LevelDBCache struct {
	db    *leveldb.DB // LevelDB instance
	batch *leveldb.Batch
	prefix DataEntryPrefix
}

func NewLevelDBCache(db *leveldb.DB, batch *leveldb.Batch, prefix DataEntryPrefix) *LevelDBCache {
	return &LevelDBCache{
		db:    db,
		batch:  batch,
		prefix: prefix,
	}
}

func (self *LevelDBCache) Get(key []byte) ([]byte, error) {
	//data, err := self.db.Get(append(self.prefix, key...), nil)
	//return data, err
	return nil,nil
}

func (self *LevelDBCache) Commit(){

}