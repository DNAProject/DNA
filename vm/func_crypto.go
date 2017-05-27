package vm


func opHash(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 1 {
		return FAULT, nil
	}
	x, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	err = PushData(e, Hash(x, e))
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opCheckSig(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 2 {
		return FAULT, nil
	}
	pubkey, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	signature, err := PopByteArray(e)
	if err != nil {
		return FAULT, err
	}
	ver, err := e.crypto.VerifySignature(e.codeContainer.GetMessage(), signature, pubkey)
	err = PushData(e, ver)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}

func opCheckMultiSig(e *ExecutionEngine) (VMState, error) {
	if Count(e) < 4 {
		return FAULT, nil
	}
	n, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}
	if n < 1 {
		return FAULT, nil
	}
	if Count(e) < n+2 {
		return FAULT, nil
	}
	e.opCount += n

	pubkeys := make([][]byte, n)
	for i := 0; i < n; i++ {
		pubkeys[i], err = PopByteArray(e)
		if err != nil {
			return FAULT, err
		}
	}

	m, err := PopInt(e)
	if err != nil {
		return FAULT, err
	}

	if m < 1 || m > n {
		return FAULT, nil
	}
	if Count(e) < m {
		return FAULT, nil
	}

	signatures := make([][]byte, m)
	for i := 0; i < m; i++ {
		signatures[i], err = PopByteArray(e)
		if err != nil {
			return FAULT, err
		}
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
	err = PushData(e, fSuccess)
	if err != nil {
		return FAULT, err
	}
	return NONE, nil
}
