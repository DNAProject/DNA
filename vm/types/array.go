package types

import (
	"math/big"
	"DNA/vm/interfaces"
)

type Array struct {
	_array []StackItemInterface
}

func NewArray(value []StackItemInterface) *Array{
	var a Array
	a._array = value
	return &a
}

func (a *Array) Equals(other StackItemInterface) bool{
	if _, ok := other.(*Array); !ok {
		return false
	}
	a1 := a._array
	a2 := other.GetArray()
	l1 := len(a1)
	l2 := len(a2)
	if l1 != l2 { return false }
	for i := 0; i<l1; i++ {
		if !a1[i].Equals(a2[i]) {
			return false
		}
	}
	return true

}

func (a *Array) GetBigInteger() *big.Int{
	if len(a._array) == 0 {  return big.NewInt(0) }
	return a._array[0].GetBigInteger()
}

func (a *Array) GetBoolean() bool{
	if len(a._array) == 0 { return false }
	return a._array[0].GetBoolean()
}

func (a *Array) GetByteArray() []byte{
	return []byte{}
}

func (a *Array) GetInterface() interfaces.IInteropInterface {
	return nil
}

func (a *Array) GetArray() []StackItemInterface{
	return a._array
}
