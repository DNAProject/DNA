package vm

import (
	"testing"
	"math/big"
	"DNA/vm/types"
)


func TestOpArraySize(t *testing.T) {
	engine.opCode = ARRAYSIZE

	bs := []byte{0x51, 0x52}
	i := big.NewInt(1)

	is := []types.StackItemInterface{types.NewByteArray(bs), types.NewInteger(i)}
	if err := PushData(engine, is); err != nil {
		t.Fatal(err)
	}
	_, err := opArraySize(engine)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("op array size result 2, execute result:", engine.GetEvaluationStack().Peek(0).GetStackItem().GetBigInteger())
}

func TestOpPack(t *testing.T) {
	engine.opCode = PACK

	bs := []byte{0x51, 0x52}
	i := big.NewInt(1)
	n := 2

	if err := PushData(engine, bs); err != nil {
		t.Fatal(err)
	}

	if err := PushData(engine, i); err != nil {
		t.Fatal(err)
	}

	if err := PushData(engine, n); err != nil {
		t.Fatal(err)
	}

	if _, err := opPack(engine); err != nil {
		t.Fatal(err)
	}
	array := engine.GetEvaluationStack().Peek(0).GetStackItem().GetArray()

	for _, v := range array {
		t.Log("value:", v.GetByteArray())
	}
}

func TestOpUnPack(t *testing.T) {
	engine.opCode = UNPACK

	if _, err := opUnpack(engine); err != nil {
		t.Fatal(err)
	}
	t.Log(engine.GetEvaluationStack().Pop().GetStackItem().GetBigInteger())
	t.Log(engine.GetEvaluationStack().Pop().GetStackItem().GetBigInteger())
	t.Log(engine.GetEvaluationStack().Pop().GetStackItem().GetByteArray())

}

func TestOpPickItem(t *testing.T) {
	engine.opCode = PICKITEM

	bs := []byte{0x51, 0x52}
	i := big.NewInt(1)

	is := []types.StackItemInterface{types.NewByteArray(bs), types.NewInteger(i)}
	if err := PushData(engine, is); err != nil {
		t.Fatal(err)
	}

	if err := PushData(engine, 0); err != nil {
		t.Fatal(err)

	}

	if _, err := opPickItem(engine); err != nil {
		t.Fatal(err)
	}
	t.Log(engine.GetEvaluationStack().Pop().GetStackItem().GetByteArray())

}


