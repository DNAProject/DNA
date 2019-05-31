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
	"DNA/smartcontract/states"
	"bytes"
	"DNA/core/store"
)

type RWSet struct {
	ReadSet map[string]*Read
	WriteSet map[string]*Write
}

type Write struct {
	Prefix store.DataEntryPrefix
	Key string
	Item states.IStateValueInterface
	IsDeleted bool
}

type Read struct {
	Key states.IStateKeyInterface
	Version string
}

func NewRWSet() *RWSet {
	var rwSet RWSet
	rwSet.WriteSet = make(map[string]*Write, 0)
	rwSet.ReadSet = make(map[string]*Read, 0)
	return &rwSet
}

func(rw *RWSet) Add(prefix store.DataEntryPrefix, key string, value states.IStateValueInterface) {
	if _, ok := rw.WriteSet[key]; !ok {
		rw.WriteSet[key] = &Write{
			Prefix: prefix,
			Key: key,
			Item: value,
			IsDeleted: false,
		}
	}

}

func(rw *RWSet) Delete(key string){
	if _, ok := rw.WriteSet[key]; ok {
		rw.WriteSet[key].Item = nil
		rw.WriteSet[key].IsDeleted = true
	}else {
		rw.WriteSet[key] = &Write{
			Key: key,
			Item: nil,
			IsDeleted: true,
		}
	}
}

func KeyToStr(key states.IStateKeyInterface) string {
	k := new(bytes.Buffer)
	key.Serialize(k)
	return k.String()
}



