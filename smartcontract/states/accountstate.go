package states

import (
	"DNA/common"
	"io"
	"DNA/common/serialization"
	. "DNA/errors"
	"bytes"
)

type AccountState struct {
	CodeHash common.Uint160
	IsFrozen bool
	Balances map[common.Uint256]common.Fixed64
	*StateBase
}

func(accountState *AccountState)Serialize(w io.Writer) error {
	accountState.CodeHash.Serialize(w)
	serialization.WriteBool(w, accountState.IsFrozen)
	serialization.WriteUint64(w, uint64(len(accountState.Balances)))
	for k, v := range accountState.Balances {
		k.Serialize(w)
		v.Serialize(w)
	}
	return nil
}

func(accountState *AccountState)Deserialize(r io.Reader) error {
	accountState.CodeHash.Deserialize(r)
	isFrozen, err := serialization.ReadBool(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "AccountState IsFrozen Deserialize fail.")
	}
	accountState.IsFrozen = isFrozen
	l, err := serialization.ReadUint64(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "AccountState Balances Len Deserialize fail.")
	}
	balances := make(map[common.Uint256]common.Fixed64, 0)
	u := new(common.Uint256)
	f := new(common.Fixed64)
	for i:=0; i<int(l); i++ {
		err = u.Deserialize(r)
		if err != nil {
			return NewDetailErr(err, ErrNoCode, "AccountState Balances key Deserialize fail.")
		}
		err = f.Deserialize(r)
		if err != nil {
			return NewDetailErr(err, ErrNoCode, "AccountState Balances value Deserialize fail.")
		}
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

