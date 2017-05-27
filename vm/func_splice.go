package vm

import (
	. "DNA/vm/errors"
)

func opCat(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, nil
	}
	b2, err :=PopByteArray(e)
	if err != nil {
		return FAULT, nil
	}
	b1, err :=PopByteArray(e)
	if err != nil {
		return FAULT, nil
	}
	if len(b1) != len(b2) {
		return FAULT, nil
	}
	r := ByteArrZip(b1, b2, CAT)
	PushData(e, r)
	return NONE, nil
}

func opSubStr(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 3 {
		return FAULT, ErrLittleLen
	}
	count, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if count < 0 {
		return FAULT, nil
	}
	index, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	s, err :=PopByteArray(e)
	if err != nil {
		return FAULT, nil
	}
	l1 := index + count
	l2 := len(s)
	if l1 > l2 {
		return FAULT, nil
	}
	b := s[index : l2-l1+1]
	err = PushData(e, b)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opLeft(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, nil
	}
	count, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if count < 0 {
		return FAULT, nil
	}
	s, err :=PopByteArray(e)
	if err != nil {
		return FAULT, nil
	}
	if count > len(s) {
		return FAULT, nil
	}
	b := s[:count]
	err = PushData(e, b)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opRight(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, nil
	}
	count, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if count < 0 {
		return FAULT, nil
	}
	s, err := PopByteArray(e)
	if err != nil {
		return FAULT, nil
	}

	l := len(s)
	if count > l {
		return FAULT, nil
	}
	b := s[l-count:]
	err = PushData(e, b)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opSize(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, nil
	}
	x := Peek(e).GetStackItem()
	if x != nil {
		return FAULT, ErrBadType
	}
	err := PushData(e, len(x.GetByteArray()))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}
