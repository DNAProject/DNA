package vm

import (
	. "DNA/vm/errors"
)

func opBigInt(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	x, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, BigIntOp(x, e.opCode))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opNot(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	x, err := PopBoolean(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, !x)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opNz(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	x, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, BigIntComp(x, e.opCode))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opBigIntZip(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, ErrLittleLen
	}
	x2, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	x1, err := PopBigInt(e)
	if err != nil {
		return FAULT, ErrBadType
	}
	b := BigIntZip(x1, x2, e.opCode)
	err = PushData(e, b)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opBoolZip(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, ErrLittleLen
	}
	x2, err := PopBoolean(e)
	if err != nil {
		return FAULT, err
	}
	x1, err := PopBoolean(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, BoolZip(x1, x2, e.opCode))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opBigIntComp(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, ErrLittleLen
	}
	x2, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	x1, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, BigIntMultiComp(x1, x2, e.opCode))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opWithIn(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 3 {
		return FAULT, ErrLittleLen
	}
	b, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	a, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	c, err := PopBigInt(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, WithInOp(c, a, b))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}
