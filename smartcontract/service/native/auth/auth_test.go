// SPDX-License-Identifier: LGPL-3.0-or-later
// Copyright 2020 DNA Dev team
//
package auth

import (
	"testing"

	"github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/smartcontract/service/native"
	common2 "github.com/DNAProject/DNA/smartcontract/service/native/common"
	"github.com/DNAProject/DNA/smartcontract/service/native/testsuite"
)

func testcase(t *testing.T, f func(t *testing.T, n *native.NativeService)) {
	testsuite.InvokeNativeContract(t, common2.AuthContractAddress,
		func(n *native.NativeService) ([]byte, error) {
			f(t, n)
			return nil, nil
		},
	)
}

func TestAuthInit(t *testing.T) {
	testcase(t, CaseAuthInit)
}

func TestAuthDoubleInitFail(t *testing.T) {
	testcase(t, CaseAuthDoubleInitFail)
}

func CaseAuthInit(t *testing.T, n *native.NativeService) {
	initSink := common.NewZeroCopySink(nil)
	initSink.WriteVarBytes(common2.DIDContractAddress[:])
	n.Input = initSink.Bytes()
	if _, err := authInit(n); err != nil {
		t.Errorf("failed to init auth: %s", err)
	}
}

func CaseAuthDoubleInitFail(t *testing.T, n *native.NativeService) {
	initSink := common.NewZeroCopySink(nil)
	initSink.WriteVarBytes(common2.DIDContractAddress[:])
	n.Input = initSink.Bytes()
	if _, err := authInit(n); err != nil {
		t.Errorf("failed to init auth: %s", err)
	}

	initSink = common.NewZeroCopySink(nil)
	initSink.WriteVarBytes(common2.DIDContractAddress[:])
	n.Input = initSink.Bytes()
	if _, err := authInit(n); err == nil {
		t.Error("failed to double init auth")
	}
}
