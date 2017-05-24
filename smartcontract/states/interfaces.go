package states

import (
	"io"
	"DNA/vm/interfaces"
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
