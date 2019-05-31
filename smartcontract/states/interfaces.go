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
	"DNA/vm/avm/interfaces"
	"DNA/core/store"
	"bytes"
)

type IStateValueInterface interface {
	Serialize(w io.Writer) error
	Deserialize(r io.Reader) error
	interfaces.IInteropInterface
}

type IStateKeyInterface interface {
	Serialize(w io.Writer) (int, error)
	Deserialize(r io.Reader) error
}

var (
	StatesMap = map[store.DataEntryPrefix]IStateValueInterface{
		store.ST_Contract: new(ContractState),
		store.ST_Storage: new(StorageItem),
		store.ST_ACCOUNT: new(AccountState),
		store.ST_AssetState: new(AssetState),
		store.ST_Validator: new(ValidatorState),
	}
)

func GetStateValue(prefix store.DataEntryPrefix, data []byte) (IStateValueInterface, error){
	r := bytes.NewBuffer(data)
	state := StatesMap[prefix]
	if err := state.Deserialize(r); err != nil {
		return nil, err
	}
	return state, nil
}
