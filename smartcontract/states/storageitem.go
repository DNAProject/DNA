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
	"io"
	"DNA/common/serialization"
	. "DNA/errors"
	"bytes"
)

type StorageItem struct {
	StateBase
	Value []byte
}

func NewStorageItem(value []byte) *StorageItem {
	var storageItem StorageItem
	storageItem.Value = value
	return &storageItem
}

func(storageItem *StorageItem)Serialize(w io.Writer) error {
	storageItem.StateBase.Serialize(w)
	serialization.WriteVarBytes(w, storageItem.Value)
	return nil
}

func(storageItem *StorageItem)Deserialize(r io.Reader) error {
	stateBase := new(StateBase)
	err := stateBase.Deserialize(r)
	if err != nil {
		return err
	}
	storageItem.StateBase = *stateBase
	value, err := serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Code Deserialize fail.")
	}
	storageItem.Value = value
	return nil
}

func(storageItem *StorageItem) ToArray() []byte {
	b := new(bytes.Buffer)
	storageItem.Serialize(b)
	return b.Bytes()
}
