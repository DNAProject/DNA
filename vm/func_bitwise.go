package vm

import (
	. "DNA/vm/errors"
)

func opInvert(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	i := x.GetBigInteger()
	err := PushData(e, i.Not(i))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opEqual(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	b1 := x
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	b2 := x
	err := PushData(e, b1.Equals(b2))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}
