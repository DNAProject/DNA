package vm

import (
	"DNA/vm/interfaces"
	"io"
	_ "math/big"
	_ "sort"
	. "DNA/vm/errors"
	"DNA/common/log"
)

func NewExecutionEngine(container interfaces.ICodeContainer, crypto interfaces.ICrypto, table interfaces.ICodeTable, service IInteropService) *ExecutionEngine {
	var engine ExecutionEngine

	engine.crypto = crypto
	engine.table = table

	engine.codeContainer = container
	engine.invocationStack = NewRandAccessStack()
	engine.opCount = 0

	engine.evaluationStack = NewRandAccessStack()
	engine.altStack = NewRandAccessStack()
	engine.state = BREAK

	engine.context = nil
	engine.opCode = 0

	engine.service = NewInteropService()

	if service != nil {
		engine.service.MergeMap(service.GetServiceMap())
	}
	return &engine
}

type ExecutionEngine struct {
	crypto  interfaces.ICrypto
	table   interfaces.ICodeTable
	service *InteropService

	codeContainer interfaces.ICodeContainer
	invocationStack *RandomAccessStack
	opCount         int

	evaluationStack *RandomAccessStack
	altStack        *RandomAccessStack
	state           VMState

	context *ExecutionContext

	//current opcode
	opCode OpCode
}

func (e *ExecutionEngine) GetState() VMState {
	return e.state
}

func (e *ExecutionEngine) GetEvaluationStack() *RandomAccessStack {
	return e.evaluationStack
}

func (e *ExecutionEngine) GetExecuteResult() bool {
	return e.evaluationStack.Pop().GetStackItem().GetBoolean()
}

func (e *ExecutionEngine) ExecutingCode() []byte {
	context := e.invocationStack.Peek(0).GetExecutionContext()
	if context != nil {
		return context.Code
	}
	return nil
}


func (e *ExecutionEngine) CurrentContext() *ExecutionContext {
	context := e.invocationStack.Peek(0).GetExecutionContext()
	if context != nil {
		return context
	}
	return nil
}

func (e *ExecutionEngine) CallingContext() *ExecutionContext {
	context := e.invocationStack.Peek(1).GetExecutionContext()
	if context !=  nil {
		return context
	}
	return nil
}

func (e *ExecutionEngine) EntryContext() *ExecutionContext {
	context := e.invocationStack.Peek(e.invocationStack.Count() - 1).GetExecutionContext()
	if context != nil {
		return context
	}
	return nil
}

func (e *ExecutionEngine) LoadCode(script []byte, pushOnly bool) {
	e.invocationStack.Push(NewExecutionContext(e, script, pushOnly, nil))
}

func (e *ExecutionEngine) Execute() {
	e.state = e.state & (^BREAK)
	for {
		if e.state == FAULT || e.state == HALT || e.state == BREAK {
			break
		}
		e.StepInto()
	}
}

func (e *ExecutionEngine) StepInto() {
	if e.invocationStack.Count() == 0 {
		e.state = HALT
		return
	}
	context := e.CurrentContext()
	if context == nil {
		e.state = FAULT
		return
	}
	var opCode OpCode

	if context.GetInstructionPointer() >= len(context.Code) {
		opCode = RET
	} else {
		o, err := context.OpReader.ReadByte()
		if err == io.EOF {
			e.state = FAULT
			return
		}
		opCode = OpCode(o)
	}
	e.opCount++
	state, err := e.ExecuteOp(opCode, context)
	if state == HALT || state == FAULT {
		e.state = state
		log.Error(err)
		return
	}
}

func (e *ExecutionEngine) ExecuteOp(opCode OpCode, context *ExecutionContext) (VMState, error) {
	if opCode > PUSH16 && opCode != RET && context.PushOnly {
		return FAULT, ErrBadValue
	}

	if opCode >= PUSHBYTES1 && opCode <= PUSHBYTES75 {
		err := PushData(e, context.OpReader.ReadBytes(int(opCode)))
		if err != nil {
			return FAULT, err
		}
		return NONE, nil
	}
	e.opCode = opCode
	e.context = context
	opExec := OpExecList[opCode]
	if opExec.Exec == nil {
		return FAULT, ErrNotSupportOpCode
	}
	return opExec.Exec(e)
}

func (e *ExecutionEngine) StepOut() {
	e.state = e.state & (^BREAK)
	c := e.invocationStack.Count()
	for {
		if e.state == FAULT || e.state == HALT || e.state == BREAK || e.invocationStack.Count() >= c {
			break
		}
		e.StepInto()
	}
}

func (e *ExecutionEngine) StepOver() {
	if e.state == FAULT || e.state == HALT {
		return
	}
	e.state = e.state & (^BREAK)
	c := e.invocationStack.Count()
	for {
		if e.state == FAULT || e.state == HALT || e.state == BREAK || e.invocationStack.Count() > c {
			break
		}
		e.StepInto()
	}
}

func (e *ExecutionEngine) AddBreakPoint(position uint) {
	//b := e.context.BreakPoints
	//b = append(b, position)
}

func (e *ExecutionEngine) RemoveBreakPoint(position uint) bool {
	//if e.invocationStack.Count() == 0 { return false }
	//b := e.context.BreakPoints
	return true
}
