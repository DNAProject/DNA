package vm

import (
	"crypto/sha1"
	"crypto/sha256"
	"hash"
)

func opHash(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 1 {
		return FAULT, nil
	}
	x := e.evaluationStack.Pop().GetStackItem().GetByteArray()
	err := PushData(e, Hash(x, e))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opCheckSig(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 2 {
		return FAULT, nil
	}
	pubkey := e.evaluationStack.Pop().GetStackItem().GetByteArray()
	signature := e.evaluationStack.Pop().GetStackItem().GetByteArray()

	ver, err := e.crypto.VerifySignature(e.codeContainer.GetMessage(), signature, pubkey)
	err = PushData(e, ver)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opCheckMultiSig(e *ExecutionEngine) (VMState, error) {
	if e.evaluationStack.Count() < 4 {
		return FAULT, nil
	}
	n := int(e.evaluationStack.Pop().GetStackItem().GetBigInteger().Int64())
	if n < 1 {
		return FAULT, nil
	}
	if e.evaluationStack.Count() < n+2 {
		return FAULT, nil
	}
	e.opCount += n

	pubkeys := make([][]byte, n)
	for i := 0; i < n; i++ {
		pubkeys[i] = e.evaluationStack.Pop().GetStackItem().GetByteArray()
	}

	m := int(e.evaluationStack.Pop().GetStackItem().GetBigInteger().Int64())
	if m < 1 || m > n {
		return FAULT, nil
	}
	if e.evaluationStack.Count() < m {
		return FAULT, nil
	}

	signatures := make([][]byte, m)
	for i := 0; i < m; i++ {
		signatures[i] = e.evaluationStack.Pop().GetStackItem().GetByteArray()
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
	err := PushData(e, fSuccess)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func Hash(b []byte, e *ExecutionEngine) []byte {
	var sh hash.Hash
	var bt []byte
	switch e.opCode {
	case SHA1:
		sh = sha1.New()
		sh.Write(b)
		bt = sh.Sum(nil)
	case SHA256:
		sh = sha256.New()
		sh.Write(b)
		bt = sh.Sum(nil)
	case HASH160:
		bt = e.crypto.Hash160(b)
	case HASH256:
		bt = e.crypto.Hash256(b)
	}
	return bt
}
