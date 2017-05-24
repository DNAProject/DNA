package vm

import (
	. "DNA/vm/errors"
)

func opArraySize(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x :=  e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	arr := x.GetArray()
	err := PushData(e, len(arr))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opPack(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	size := int(x.GetBigInteger().Int64())
	if size < 0 || size > e.evaluationStack.Count() {
		return FAULT, ErrBadValue
	}
	items := NewStackItems()
	for {
		if size == 0 {
			break
		}
		items = append(items, e.evaluationStack.Pop().GetStackItem())
		size--
	}
	err := PushData(e, items)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opUnpack(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadValue
	}
	arr := x.GetArray()
	l := len(arr)
	for i := l - 1; i >= 0; i-- {
		e.evaluationStack.Push(NewStackItem(arr[i]))
	}
	err := PushData(e, l)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opPickItem(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadValue
	}
	index := int(x.GetBigInteger().Int64())
	if index < 0 {
		return FAULT, ErrBadValue
	}
	x = e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadValue
	}
	items := x.GetArray()
	if index >= len(items) {
		return FAULT, ErrOverLen
	}
	err := PushData(e, items[index])
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}
