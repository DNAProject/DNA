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


func opHash(e *ExecutionEngine) (VMState, error) {
	x := PopByteArray(e)
	PushData(e, Hash(x, e))
	return NONE, nil
}

func opCheckSig(e *ExecutionEngine) (VMState, error) {
	pubkey := PopByteArray(e)
	signature := PopByteArray(e)
	ver, err := e.crypto.VerifySignature(e.codeContainer.GetMessage(), signature, pubkey)
	if err != nil { return FAULT, err }
	PushData(e, ver)
	return NONE, nil
}

func opCheckMultiSig(e *ExecutionEngine) (VMState, error) {
	n := PopInt(e)
	if n < 1 { return FAULT, nil }
	if Count(e) < n+2 {
		return FAULT, nil
	}
	e.opCount += n

	pubkeys := make([][]byte, n)
	for i := 0; i < n; i++ {
		pubkeys[i] = PopByteArray(e)
	}

	m := PopInt(e)
	if m < 1 || m > n { return FAULT, nil }

	signatures := make([][]byte, m)
	for i := 0; i < m; i++ {
		signatures[i] = PopByteArray(e)
	}

	message := e.codeContainer.GetMessage()
	fSuccess := true

	for i, j := 0, 0; fSuccess && i < m && j < n; {
		ver, _ := e.crypto.VerifySignature(message, signatures[i], pubkeys[j])
		if ver {
			i++
		}
		j++
		if m-i > n-j {
			fSuccess = false
		}
	}
	PushData(e, fSuccess)
	return NONE, nil
}
