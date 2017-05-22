package payload

import (
	"io"
	"DNA/common/serialization"
)

type InvokeCode struct {
	Code []byte
}

func (ic *InvokeCode) Data() []byte {

	return []byte{0}
}

func (ic *InvokeCode) Serialize(w io.Writer) error {
	err := serialization.WriteVarBytes(w, ic.Code)
	if err != nil {
		return err
	}

	return nil
}

func (ic *InvokeCode) Deserialize(r io.Reader) error {
	code, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	ic.Code = code
	return nil
}
