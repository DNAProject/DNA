package vm

import (
	. "DNA/vm/errors"
)

func opNop(e *ExecutionEngine) (VMState, error) {
	return NONE, nil
}

func opJmp(e *ExecutionEngine) (VMState, error) {
	offset := int(e.context.OpReader.ReadInt16())

	offset = e.context.GetInstructionPointer() + offset - 3

	if offset < 0 || offset > len(e.context.Code) {
		return FAULT, ErrFault
	}
	fValue := true
	if e.opCode > JMP {
		x := e.evaluationStack.Pop().GetStackItem()
		if x == nil {
			return FAULT, ErrBadType
		}
		fValue = x.GetBoolean()
		if e.opCode == JMPIFNOT {
			fValue = !fValue
		}
	}
	if fValue {
		e.context.SetInstructionPointer(int64(offset))
	}
	return NONE, nil
}

func opCall(e *ExecutionEngine) (VMState, error) {
	e.invocationStack.Push(e.context.Clone())
	e.context.SetInstructionPointer(int64(e.context.GetInstructionPointer() + 2))
	e.opCode = JMP
	e.context = e.CurrentContext()
	opJmp(e)
	return NONE, nil
}

func opRet(e *ExecutionEngine) (VMState, error) {
	e.invocationStack.Pop()
	if e.invocationStack.Count() == 0 {
		return HALT, ErrLittleLen
	}
	return NONE, nil
}

func opAppCall(e *ExecutionEngine) (VMState, error) {
	if e.table == nil {
		return FAULT, nil
	}
	codeHash := e.context.OpReader.ReadBytes(20)
	code := e.table.GetCode(codeHash)
	if code == nil {
		return FAULT, nil
	}
	if e.opCode == TAILCALL {
		e.invocationStack.Pop()
	}
	e.LoadCode(code, false)
	return NONE, nil
}

func opSysCall(e *ExecutionEngine) (VMState, error) {
	if e.service == nil {
		return FAULT, nil
	}

	s := e.context.OpReader.ReadVarString()

	success, err := e.service.Invoke(s, e)
	if success {
		return NONE, nil
	} else {
		return FAULT, err
	}
}

