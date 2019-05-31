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

package states

import (
	"DNA/common"
	"io"
	"DNA/common/serialization"
	. "DNA/errors"
)

type StorageKey struct {
	CodeHash *common.Uint160
	Key []byte
}

func NewStorageKey(codeHash *common.Uint160, key []byte) *StorageKey {
	var storageKey StorageKey
	storageKey.CodeHash = codeHash
	storageKey.Key = key
	return &storageKey
}

func (storageKey *StorageKey) Serialize(w io.Writer) (int, error) {
	storageKey.CodeHash.Serialize(w)
	serialization.WriteVarBytes(w, storageKey.Key)
	return 0, nil
}

func (storageKey *StorageKey) Deserialize(r io.Reader) error {
	u := new(common.Uint160)
	err := u.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "StorageKey CodeHash Deserialize fail.")
	}
	storageKey.CodeHash = u
	key, err := serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "StorageKey Key Deserialize fail.")
	}
	storageKey.Key = key
	return nil
}

