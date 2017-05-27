package vm

import (
	. "DNA/vm/errors"
)

func opInvert(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	i, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, i.Not(i))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opEqual(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, ErrLittleLen
	}
	x := PopStackItem(e)
	if x == nil {
		return FAULT, ErrBadType
	}
	b1 := x
	x = PopStackItem(e)
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
