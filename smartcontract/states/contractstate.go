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
	. "DNA/errors"
	"DNA/core/code"
	"bytes"
	"DNA/common/serialization"
	"DNA/smartcontract/types"
	"DNA/common"
)

type ContractState struct {
	StateBase
	Code *code.FunctionCode
	Name string
	Version string
	Author string
	Email string
	Description string
	Language types.LangType
	ProgramHash common.Uint160
}

func(contractState *ContractState) Serialize(w io.Writer) error {
	contractState.StateBase.Serialize(w)
	err := contractState.Code.Serialize(w)
	if err != nil {
		return err
	}
	err = serialization.WriteVarString(w, contractState.Name)
	if err != nil {
		return err
	}
	err = serialization.WriteVarString(w, contractState.Version)
	if err != nil {
		return err
	}
	err = serialization.WriteVarString(w, contractState.Author)
	if err != nil {
		return err
	}
	err = serialization.WriteVarString(w, contractState.Email)
	if err != nil {
		return err
	}
	err = serialization.WriteVarString(w, contractState.Description)
	if err != nil {
		return err
	}
	err = serialization.WriteByte(w, byte(contractState.Language))
	if err != nil {
		return err
	}
	_, err = contractState.ProgramHash.Serialize(w)
	if err != nil {
		return err
	}
	return nil
}

func(contractState *ContractState) Deserialize(r io.Reader) error {
	stateBase := new(StateBase)
	err := stateBase.Deserialize(r)
	if err != nil {
		return err
	}
	contractState.StateBase = *stateBase
	f := new(code.FunctionCode)
	err = f.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Code Deserialize fail.")
	}
	contractState.Code = f
	contractState.Name, err = serialization.ReadVarString(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Name Deserialize fail.")

	}
	contractState.Version, err = serialization.ReadVarString(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Version Deserialize fail.")

	}
	contractState.Author, err = serialization.ReadVarString(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Author Deserialize fail.")

	}
	contractState.Email, err = serialization.ReadVarString(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Email Deserialize fail.")

	}
	contractState.Description, err = serialization.ReadVarString(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Description Deserialize fail.")

	}
	language, err := serialization.ReadByte(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Description Deserialize fail.")

	}
	contractState.Language = types.LangType(language)
	u := new(common.Uint160)
	err = u.Deserialize(r)
	if err != nil {
		return err
	}
	contractState.ProgramHash = *u
	return nil
}

func(contractState *ContractState) ToArray() []byte {
	b := new(bytes.Buffer)
	contractState.Serialize(b)
	return b.Bytes()
}


