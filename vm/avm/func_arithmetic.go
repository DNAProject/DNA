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

func opBigInt(e *ExecutionEngine) (VMState, error) {
	x := PopBigInt(e)
	PushData(e, BigIntOp(x, e.opCode))
	return NONE, nil
}

func opNot(e *ExecutionEngine) (VMState, error) {
	x := PopBoolean(e)
	PushData(e, !x)
	return NONE, nil
}

func opNz(e *ExecutionEngine) (VMState, error) {
	x := PopBigInt(e)
	PushData(e, BigIntComp(x, e.opCode))
	return NONE, nil
}

func opBigIntZip(e *ExecutionEngine) (VMState, error) {
	x2 := PopBigInt(e)
	x1 := PopBigInt(e)
	b := BigIntZip(x1, x2, e.opCode)
	PushData(e, b)
	return NONE, nil
}

func opBoolZip(e *ExecutionEngine) (VMState, error) {
	x2 := PopBoolean(e)
	x1 := PopBoolean(e)
	PushData(e, BoolZip(x1, x2, e.opCode))
	return NONE, nil
}

func opBigIntComp(e *ExecutionEngine) (VMState, error) {
	x2 := PopBigInt(e)
	x1 := PopBigInt(e)
	PushData(e, BigIntMultiComp(x1, x2, e.opCode))
	return NONE, nil
}

func opWithIn(e *ExecutionEngine) (VMState, error) {
	b := PopBigInt(e)
	a := PopBigInt(e)
	c := PopBigInt(e)
	PushData(e, WithInOp(c, a, b))
	return NONE, nil
}
