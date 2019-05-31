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

package account

import (
	"DNA/common"
	"io"
	"bytes"
	"DNA/common/serialization"
)

type AccountState struct {
	ProgramHash common.Uint160
	IsFrozen bool
	Balances map[common.Uint256]common.Fixed64
}

func NewAccountState(programHash common.Uint160, balances map[common.Uint256]common.Fixed64) *AccountState {
	var accountState AccountState
	accountState.ProgramHash = programHash
	accountState.Balances = balances
	accountState.IsFrozen = false
	return &accountState
}

func(accountState *AccountState)Serialize(w io.Writer) error {
	accountState.ProgramHash.Serialize(w)
	serialization.WriteBool(w, accountState.IsFrozen)
	serialization.WriteUint64(w, uint64(len(accountState.Balances)))
	for k, v := range accountState.Balances {
		k.Serialize(w)
		v.Serialize(w)
	}
	return nil
}

func(accountState *AccountState)Deserialize(r io.Reader) error {
	accountState.ProgramHash.Deserialize(r)
	isFrozen, err := serialization.ReadBool(r)
	if err != nil { return err }
	accountState.IsFrozen = isFrozen
	l, err := serialization.ReadUint64(r)
	if err != nil { return err }
	balances := make(map[common.Uint256]common.Fixed64, 0)
	u := new(common.Uint256)
	f := new(common.Fixed64)
	for i:=0; i<int(l); i++ {
		err = u.Deserialize(r)
		if err != nil { return err }
		err = f.Deserialize(r)
		if err != nil { return err }
		balances[*u] = *f
	}
	accountState.Balances = balances
	return nil
}

func(accountState *AccountState) ToArray() []byte {
	b := new(bytes.Buffer)
	accountState.Serialize(b)
	return b.Bytes()
}


