package vm

import (
	. "DNA/vm/errors"
)

func opToAltStack(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	e.altStack.Push(e.evaluationStack.Pop())
	return NONE, nil
}

func opFromAltStack(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	e.evaluationStack.Push(e.altStack.Pop())
	return NONE, nil
}

func opXDrop(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	n := int(x.GetBigInteger().Int64())
	if n < 0 {
		return FAULT, nil
	}
	e.evaluationStack.Remove(n)
	return NONE, nil
}

func opXSwap(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	n := int(x.GetBigInteger().Int64())
	if n < 0 || n > e.evaluationStack.Count()-1 {
		return FAULT, ErrBadValue
	}
	if n == 0 {
		return NONE, ErrBadValue
	}
	e.evaluationStack.Swap(0, n)
	return NONE, nil
}

func opXTuck(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	n := int(x.GetBigInteger().Int64())

	if n < 0 || n > e.evaluationStack.Count()-1 {
		return FAULT, ErrBadValue
	}
	e.evaluationStack.Insert(n, e.evaluationStack.Peek(0))
	return NONE, nil
}

func opDepth(e *ExecutionEngine) (VMState, error) {
	PushData(e, e.evaluationStack.Count())
	return NONE, nil
}

func opDrop(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	e.evaluationStack.Pop()
	return NONE, nil
}

func opDup(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, ErrLittleLen
	}
	e.evaluationStack.Push(e.evaluationStack.Peek(0))
	return NONE, nil
}

func opNip(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, ErrLittleLen
	}
	x2 := e.evaluationStack.Pop()
	e.evaluationStack.Pop()
	e.evaluationStack.Push(x2)
	return NONE, nil
}

func opOver(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, ErrLittleLen
	}
	x2 := e.evaluationStack.Pop()
	x1 := e.evaluationStack.Peek(0)
	e.evaluationStack.Push(x2)
	e.evaluationStack.Push(x1)
	return NONE, nil
}

func opPick(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, nil
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	n := int(x.GetBigInteger().Int64())
	if n < 0 {
		return FAULT, nil
	}
	e.evaluationStack.Push(e.evaluationStack.Peek(n))
	return NONE, nil
}

func opRoll(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, ErrOverLen
	}
	x := e.evaluationStack.Pop().GetStackItem()
	if x == nil {
		return FAULT, ErrBadType
	}
	n := int(x.GetBigInteger().Int64())
	if n < 0 {
		return FAULT, ErrBadType
	}
	if n == 0 {
		return NONE, nil
	}
	if e.evaluationStack.Count() < n+1 {
		return FAULT, ErrBadType
	}
	e.evaluationStack.Push(e.evaluationStack.Remove(n))
	return NONE, nil
}

func opRot(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 3 {
		return FAULT, nil
	}
	x3 := e.evaluationStack.Pop()
	x2 := e.evaluationStack.Pop()
	x1 := e.evaluationStack.Pop()
	e.evaluationStack.Push(x2)
	e.evaluationStack.Push(x3)
	e.evaluationStack.Push(x1)
	return NONE, nil
}

func opSwap(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, nil
	}
	x2 := e.evaluationStack.Pop()
	x1 := e.evaluationStack.Pop()
	e.evaluationStack.Push(x2)
	e.evaluationStack.Push(x1)
	return NONE, nil
}

func opTuck(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, nil
	}
	x2 := e.evaluationStack.Pop()
	x1 := e.evaluationStack.Pop()
	e.evaluationStack.Push(x2)
	e.evaluationStack.Push(x1)
	e.evaluationStack.Push(x2)
	return NONE, nil
}

func PushData(e *ExecutionEngine, data interface{}) error {
	d, err := NewStackItemInterface(data)
	if err != nil {
		return err
	}
	e.evaluationStack.Push(NewStackItem(d))
	return nil
}
