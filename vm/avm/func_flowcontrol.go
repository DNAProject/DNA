// Copyright 2016 DNA Dev team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package avm

import (
	. "DNA/vm/avm/errors"
	"DNA/common/log"
)

func opNop(e *ExecutionEngine) (VMState, error) {
	return NONE, nil
}

func opJmp(e *ExecutionEngine) (VMState, error) {
	offset := int(e.context.OpReader.ReadInt16())

	offset = e.context.GetInstructionPointer() + offset - 3

	if offset > len(e.context.Code) {
		return FAULT, ErrFault
	}
	var (
		fValue = true
	)
	if e.opCode > JMP {
		fValue = PopBoolean(e)

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
	return NONE, nil
}

func opAppCall(e *ExecutionEngine) (VMState, error) {
	codeHash := e.context.OpReader.ReadBytes(20)
	code, err := e.table.GetCode(codeHash)
	if code == nil {
		return FAULT, err
	}
	if e.opCode == TAILCALL {
		e.invocationStack.Pop()
	}
	e.LoadCode(code, false)
	return NONE, nil
}

func opSysCall(e *ExecutionEngine) (VMState, error) {
	s := e.context.OpReader.ReadVarString()

	log.Error("[opSysCall] service name:", s)

	success, err := e.service.Invoke(s, e)
	if success {
		return NONE, nil
	} else {
		return FAULT, err
	}
}

