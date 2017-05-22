package vm

import (
	. "DNA/vm/errors"
)

func opBigInt(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	err := PushData(e, BigIntOp(x.GetBigInteger(), e.opCode))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opNot(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	err := PushData(e, !x.GetBoolean())
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opNz(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	err := PushData(e, BigIntComp(x.GetBigInteger(), e.opCode))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opBigIntZip(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	x2 := x.GetBigInteger()
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	x1 := x.GetBigInteger()
	b := BigIntZip(x1, x2, e.opCode)
	err := PushData(e, b)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opBoolZip(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	x2 := x.GetBoolean()
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	x1 := x.GetBoolean()
	err := PushData(e, BoolZip(x1, x2, e.opCode))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opBigIntComp(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	x2 := x.GetBigInteger()
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	x1 := x.GetBigInteger()
	err := PushData(e, BigIntMultiComp(x1, x2, e.opCode))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opWithIn(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 3 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	b := x.GetBigInteger()
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	a := x.GetBigInteger()
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	c := x.GetBigInteger()
	err := PushData(e, WithInOp(c, a, b))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}
