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
	"DNA/vm/avm/utils"
	"io"
	"DNA/vm/avm/types"
)

type ExecutionContext struct {
	Code               []byte
	OpReader           *utils.VmReader
	PushOnly           bool
	BreakPoints        []uint
	InstructionPointer int
	CodeHash           []byte
	engine             *ExecutionEngine
}

func NewExecutionContext(engine *ExecutionEngine, code []byte, pushOnly bool, breakPoints []uint) *ExecutionContext {
	var executionContext ExecutionContext
	executionContext.Code = code
	executionContext.OpReader = utils.NewVmReader(code)
	executionContext.PushOnly = pushOnly
	executionContext.BreakPoints = breakPoints
	executionContext.InstructionPointer = 0
	executionContext.engine = engine
	return &executionContext
}

func (ec *ExecutionContext) GetInstructionPointer() int {
	return ec.OpReader.Position()
}

func (ec *ExecutionContext) SetInstructionPointer(offset int64) {
	ec.OpReader.Seek(offset, io.SeekStart)
}

func (ec *ExecutionContext) GetCodeHash() []byte {
	if ec.CodeHash == nil {
		ec.CodeHash = ec.engine.crypto.Hash160(ec.Code)
	}
	return ec.CodeHash
}

func (ec *ExecutionContext) NextInstruction() OpCode {
	return OpCode(ec.Code[ec.OpReader.Position()])
}

func (ec *ExecutionContext) Clone() *ExecutionContext {
	executionContext := NewExecutionContext(ec.engine, ec.Code, ec.PushOnly, ec.BreakPoints)
	executionContext.InstructionPointer = ec.InstructionPointer
	executionContext.SetInstructionPointer(int64(ec.GetInstructionPointer()))
	return executionContext
}

func (ec *ExecutionContext) GetStackItem() types.StackItemInterface {
	return nil
}

func (ec *ExecutionContext) GetExecutionContext() *ExecutionContext {
	return ec
}


