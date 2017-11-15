package payload

import (
	"DNA/common"
	"DNA/common/serialization"
	"io"
)

type InvokeCode struct {
	CodeHash    common.Uint160
	Code        []byte
	ProgramHash common.Uint160
}

func (ic *InvokeCode) Data(version byte) []byte {
	return []byte{0}
}

func (ic *InvokeCode) Serialize(w io.Writer, version byte) error {
	ic.CodeHash.Serialize(w)
	err := serialization.WriteVarBytes(w, ic.Code)
	if err != nil {
		return err
	}
	_, err = ic.ProgramHash.Serialize(w)
	if err != nil {
		return err
	}
	return nil
}

func (ic *InvokeCode) Deserialize(r io.Reader, version byte) error {
	u := new(common.Uint160)
	if err := u.Deserialize(r); err != nil {
		return err
	}
	ic.CodeHash = *u
	code, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	ic.Code = code

	p := new(common.Uint160)
	err = p.Deserialize(r)
	if err != nil {
		return err
	}
	ic.ProgramHash = *p
	return nil
}
