package states

import (
	"DNA/common"
	"io"
	. "DNA/errors"
	"DNA/core/code"
	"bytes"
)

type ContractState struct {
	Code *code.FunctionCode
	CodeHash common.Uint160
	*StateBase
}

func(contractState *ContractState)Serialize(w io.Writer) error {
	contractState.Code.Serialize(w)
	contractState.CodeHash.Serialize(w)
	return nil
}

func(contractState *ContractState)Deserialize(r io.Reader) error {
	u := new(common.Uint160)
	f := new(code.FunctionCode)
	err := f.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Code Deserialize fail.")
	}
	contractState.Code = f
	err = u.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "AssetState CodeHash Deserialize fail.")
	}
	contractState.CodeHash = *u
	return nil
}

func(contractState *ContractState) ToArray() []byte {
	b := new(bytes.Buffer)
	contractState.Serialize(b)
	return b.Bytes()
}


