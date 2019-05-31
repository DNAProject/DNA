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

func opPushData(e *ExecutionEngine) (VMState, error) {
	data := getPushData(e)
	PushData(e, data)
	return NONE, nil
}

func getPushData(e *ExecutionEngine) interface{} {
	var data interface{}
	if e.opCode >= PUSHBYTES1 && e.opCode <= PUSHBYTES75 {
		data = e.context.OpReader.ReadBytes(int(e.opCode))
	}
	switch e.opCode {
	case PUSH0:
		data = []byte{}
	case PUSHDATA1:
		d, _ := e.context.OpReader.ReadByte()
		data = e.context.OpReader.ReadBytes(int(d))
	case PUSHDATA2:
		data = e.context.OpReader.ReadBytes(int(e.context.OpReader.ReadUint16()))
	case PUSHDATA4:
		i := int(e.context.OpReader.ReadInt32())
		data = e.context.OpReader.ReadBytes(i)
	case PUSHM1, PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8, PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16:
		data = int8(e.opCode - PUSH1 + 1)
	}

	return data
}
