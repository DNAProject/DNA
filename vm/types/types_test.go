package types

import (
	"testing"
	"math/big"
	"fmt"
)

func TestTypes(t *testing.T) {
	i := NewInteger(big.NewInt(1))
	ba := NewByteArray([]byte{1})
	b := NewBoolean(false)
	a1 := NewArray([]StackItem{i})
	//a2 := NewArray([]StackItem{ba})
	fmt.Printf("%+v", i.GetByteArray())
	fmt.Printf("\n%+v", ba.GetBoolean())
	fmt.Printf("\n%v", b.Equals(NewBoolean(false)))
	fmt.Println("\nequal", a1.Equals(NewArray([]StackItem{NewInteger(big.NewInt(1))})))
}
