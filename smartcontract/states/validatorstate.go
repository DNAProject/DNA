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
	"DNA/crypto"
	"io"
	"bytes"
)

type ValidatorState struct {
	StateBase
	PublicKey *crypto.PubKey
}

func(v *ValidatorState) Serialize(w io.Writer) error {
	v.StateBase.Serialize(w)
	v.PublicKey.Serialize(w)
	return nil
}


func(v *ValidatorState)Deserialize(r io.Reader) error {
	stateBase := new(StateBase)
	err := stateBase.Deserialize(r)
	if err != nil {
		return err
	}
	v.StateBase = *stateBase
	p := new(crypto.PubKey)
	err = p.DeSerialize(r)
	if err != nil {
		return err
	}
	v.PublicKey = p
	return nil
}

func(v *ValidatorState) ToArray() []byte {
	b := new(bytes.Buffer)
	v.Serialize(b)
	return b.Bytes()
}

