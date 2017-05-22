package vm

import (
	. "DNA/vm/errors"
)

func opCat(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, nil
	}
	x2 := e.evaluationStack.Pop()
	x1 := e.evaluationStack.Pop()
	b1 := x1.GetStackItem().GetByteArray()
	b2 := x2.GetStackItem().GetByteArray()
	if len(b1) != len(b2) {
		return FAULT, nil
	}
	r := ByteArrZip(b1, b2, CAT)
	PushData(e, r)
	return NONE, nil
}

func opSubStr(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 3 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	count := int(x.GetBigInteger().Int64())
	if count < 0 {
		return FAULT, nil
	}
	index := int(e.evaluationStack.Pop().GetStackItem().GetBigInteger().Int64())
	if index < 0 {
		return FAULT, nil
	}
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	s := x.GetByteArray()
	l1 := index + count
	l2 := len(s)
	if l1 > l2 {
		return FAULT, nil
	}
	b := s[index : l2-l1+1]
	err := PushData(e, b)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opLeft(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, nil
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	count := int(x.GetBigInteger().Int64())
	if count < 0 {
		return FAULT, nil
	}
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	s := x.GetByteArray()
	if count > len(s) {
		return FAULT, nil
	}
	b := s[:count]
	err := PushData(e, b)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opRight(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, nil
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	count := int(x.GetBigInteger().Int64())
	if count < 0 {
		return FAULT, nil
	}
	if x == nil {
		return FAULT, ErrBadType
	}
	s := x.GetByteArray()

	l := len(s)
	if count > l {
		return FAULT, nil
	}
	b := s[l-count:]
	err := PushData(e, b)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opSize(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, nil
	}
	x := e.evaluationStack.Peek(0)
	s := x.GetStackItem().GetByteArray()
	err := PushData(e, len(s))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}
