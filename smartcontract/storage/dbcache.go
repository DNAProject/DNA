// Copyright 2016 DNA Dev team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"DNA/core/store"
	"DNA/smartcontract/states"
	"DNA/common"
	"math/big"
)

type DBCache interface {
	GetOrAdd(prefix store.DataEntryPrefix, key string, value states.IStateValueInterface) (states.IStateValueInterface, error)
	TryGet(prefix store.DataEntryPrefix, key string) (states.IStateValueInterface, error)
	GetWriteSet() *RWSet
	GetState(codeHash common.Uint160, loc common.Hash) (common.Hash, error)
	SetState(codeHash common.Uint160, loc, value common.Hash)
	GetCode(codeHash common.Uint160) ([]byte, error)
	SetCode(codeHash common.Uint160, code []byte)
	GetBalance(common.Uint160) *big.Int
	GetCodeSize(common.Uint160) int
	AddBalance(common.Uint160, *big.Int)
	Suicide(codeHash common.Uint160) bool
}

type CloneCache struct {
	innerCache DBCache
	dbCache DBCache
}

func NewCloneDBCache(innerCache DBCache, dbCache DBCache) *CloneCache {
	return &CloneCache{
		innerCache: innerCache,
		dbCache: dbCache,
	}
}

func (cloneCache *CloneCache) GetInnerCache() DBCache {
	return cloneCache.innerCache
}

func (cloneCache *CloneCache) Commit() {
	for _, v := range cloneCache.innerCache.GetWriteSet().WriteSet {
		if v.IsDeleted {
			cloneCache.dbCache.GetWriteSet().Delete(v.Key)
		}else {
			cloneCache.dbCache.GetWriteSet().Add(v.Prefix, v.Key, v.Item)
		}
	}
}

func (cloneCache *CloneCache) TryGet(prefix store.DataEntryPrefix, key string) (states.IStateValueInterface, error) {
	if v, ok := cloneCache.innerCache.GetWriteSet().WriteSet[key]; ok {
		return v.Item, nil
	}else {
		return cloneCache.dbCache.TryGet(prefix, key)
	}
}


