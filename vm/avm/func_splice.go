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

func opCat(e *ExecutionEngine) (VMState, error) {
	b2 := PopByteArray(e)
	b1 := PopByteArray(e)
	r := ByteArrZip(b1, b2, CAT)
	PushData(e, r)
	return NONE, nil
}

func opSubStr(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	index := PopInt(e)
	arr  := PopByteArray(e)
	b := arr[index: index + count]
	PushData(e, b)
	return NONE, nil
}

func opLeft(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	s := PopByteArray(e)
	b := s[:count]
	PushData(e, b)
	return NONE, nil
}

func opRight(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	arr := PopByteArray(e)
	b := arr[len(arr)-count:]
	PushData(e, b)
	return NONE, nil
}

func opSize(e *ExecutionEngine) (VMState, error) {
	x := Peek(e).GetStackItem()
	PushData(e, len(x.GetByteArray()))
	return NONE, nil
}
