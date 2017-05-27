package vm

import (
	. "DNA/vm/errors"
)

func opToAltStack(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	e.altStack.Push(Pop(e))
	return NONE, nil
}

func opFromAltStack(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	Push(e, e.altStack.Pop())
	return NONE, nil
}

func opXDrop(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if n < 0 {
		return FAULT, nil
	}
	e.evaluationStack.Remove(n)
	return NONE, nil
}

func opXSwap(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if n < 0 || n > Count(e)-1 {
		return FAULT, ErrBadValue
	}
	if n == 0 {
		return NONE, nil
	}
	e.evaluationStack.Swap(0, n)
	return NONE, nil
}

func opXTuck(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}

	if n < 0 || n > Count(e)-1 {
		return FAULT, ErrBadValue
	}
	e.evaluationStack.Insert(n, Peek(e))
	return NONE, nil
}

func opDepth(e *ExecutionEngine) (VMState, error) {
	PushData(e, Count(e))
	return NONE, nil
}

func opDrop(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	Pop(e)
	return NONE, nil
}

func opDup(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, ErrLittleLen
	}
	Push(e, Peek(e))
	return NONE, nil
}

func opNip(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, ErrLittleLen
	}
	x2 := Pop(e)
	Pop(e)
	Push(e, x2)
	return NONE, nil
}

func opOver(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, ErrLittleLen
	}
	x2 := Pop(e)
	x1 := Peek(e)
	Push(e, x2)
	Push(e, x1)
	return NONE, nil
}

func opPick(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, nil
	}
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if n < 0 {
		return FAULT, nil
	}
	Push(e, e.evaluationStack.Peek(n))
	return NONE, nil
}

func opRoll(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, ErrOverLen
	}
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if n < 0 {
		return FAULT, ErrBadType
	}
	if n == 0 {
		return NONE, nil
	}
	if Count(e) < n+1 {
		return FAULT, ErrBadType
	}
	Push(e, e.evaluationStack.Remove(n))
	return NONE, nil
}

func opRot(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 3 {
		return FAULT, nil
	}
	x3 := Pop(e)
	x2 := Pop(e)
	x1 := Pop(e)
	Push(e, x2)
	Push(e, x3)
	Push(e, x1)
	return NONE, nil
}

func opSwap(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, nil
	}
	x2 := Pop(e)
	x1 := Pop(e)
	Push(e, x2)
	Push(e, x1)
	return NONE, nil
}

func opTuck(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, nil
	}
	x2 := Pop(e)
	x1 := Pop(e)
	Push(e, x2)
	Push(e, x1)
	Push(e, x2)
	return NONE, nil
}

