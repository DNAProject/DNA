package vm

import (
	. "DNA/vm/errors"
)

func opArraySize(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	arr, err := PopArray(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, len(arr))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opPack(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	size, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if size < 0 || size > e.evaluationStack.Count() {
		return FAULT, ErrBadValue
	}
	items := NewStackItems()
	for {
		if size == 0 {
			break
		}
		items = append(items, PopStackItem(e))
		size--
	}
	err = PushData(e, items)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opUnpack(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	arr, err := PopArray(e)
	if err != nil {
		return FAULT, err
	}
	l := len(arr)
	for i := l - 1; i >= 0; i-- {
		Push(e, NewStackItem(arr[i]))
	}
	err = PushData(e, l)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opPickItem(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	index, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if index < 0 {
		return FAULT, ErrBadValue
	}
	items, err := PopArray(e)
	if err != nil {
		return FAULT, err
	}
	if index >= len(items) {
		return FAULT, ErrOverLen
	}
	err = PushData(e, items[index])
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}
