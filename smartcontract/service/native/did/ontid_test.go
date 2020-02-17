// SPDX-License-Identifier: LGPL-3.0-or-later
// Copyright 2019 DNA Dev team
//
/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */
package did

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/DNAProject/DNA/account"
	"github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/smartcontract/service/native"
	common2 "github.com/DNAProject/DNA/smartcontract/service/native/common"
	"github.com/DNAProject/DNA/smartcontract/service/native/testsuite"
	"github.com/DNAProject/DNA/smartcontract/service/native/utils"
	"github.com/ontio/ontology-crypto/keypair"
)

func testcase(t *testing.T, f func(t *testing.T, n *native.NativeService)) {
	testsuite.InvokeNativeContract(t, common2.DIDContractAddress,
		func(n *native.NativeService) ([]byte, error) {
			f(t, n)
			return nil, nil
		},
	)
}

func TestReg(t *testing.T) {
	testcase(t, CaseRegID)
}

func TestOwner(t *testing.T) {
	testcase(t, CaseOwner)
}

func TestOwnerSize(t *testing.T) {
	testcase(t, CaseOwnerSize)
}

func TestDoubleInitFail(t *testing.T) {
	testcase(t, CaseDoubleInit)
}

// Register id with account acc
func regID(n *native.NativeService, id string, a *account.Account) error {
	// make arguments
	sink := common.NewZeroCopySink(nil)
	sink.WriteVarBytes([]byte(id))
	pk := keypair.SerializePublicKey(a.PubKey())
	sink.WriteVarBytes(pk)
	n.Input = sink.Bytes()
	// set signing address
	n.Tx.SignedAddr = []common.Address{a.Address}
	// call
	_, err := regIdWithPublicKey(n)
	return err
}

func CaseRegID(t *testing.T, n *native.NativeService) {
	initSink := common.NewZeroCopySink(nil)
	initSink.WriteVarBytes([]byte("dna"))
	n.Input = initSink.Bytes()
	if _, err := didInit(n); err != nil {
		t.Errorf("failed to init did: %s", err)
	}

	id, err := account.GenerateID()
	if err != nil {
		t.Fatal(err)
	}
	a := account.NewAccount("")

	// 1. register invalid id, should fail
	if err := regID(n, "did:dna:abcd1234", a); err == nil {
		t.Error("invalid id registered")
	}

	// 2. register without valid signature, should fail
	sink := common.NewZeroCopySink(nil)
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{}
	if _, err := regIdWithPublicKey(n); err == nil {
		t.Error("id registered without signature")
	}

	// 3. register with invalid key, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes([]byte("invalid public key"))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a.Address}
	if _, err := regIdWithPublicKey(n); err == nil {
		t.Error("id registered with invalid key")
	}

	// 4. register id
	if err := regID(n, id, a); err != nil {
		t.Fatal(err)
	}

	// 5. get DDO
	sink.Reset()
	sink.WriteString(id)
	n.Input = sink.Bytes()
	_, err = GetDDO(n)
	if err != nil {
		t.Error(err)
	}

	// 6. register again, should fail
	if err := regID(n, id, a); err == nil {
		t.Error("id registered twice")
	}

	// 7. revoke with invalid key, should fail
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 2)
	n.Input = sink.Bytes()
	if _, err := revokeID(n); err == nil {
		t.Error("revoked by invalid key")
	}

	// 8. revoke without valid signature, should fail
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 1)
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{common.ADDRESS_EMPTY}
	if _, err := revokeID(n); err == nil {
		t.Error("revoked without valid signature")
	}

	// 9. revoke id
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 1)
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a.Address}
	if _, err := revokeID(n); err != nil {
		t.Fatal(err)
	}

	// 10. register again, should fail
	if err := regID(n, id, a); err == nil {
		t.Error("revoked id should not be registered again")
	}

	// 11. get DDO of the revoked id
	sink.Reset()
	sink.WriteString(id)
	n.Input = sink.Bytes()
	_, err = GetDDO(n)
	if err == nil {
		t.Error("get DDO of the revoked id should fail")
	}

}

func CaseOwner(t *testing.T, n *native.NativeService) {
	initSink := common.NewZeroCopySink(nil)
	initSink.WriteVarBytes([]byte("dna"))
	n.Input = initSink.Bytes()
	if _, err := didInit(n); err != nil {
		t.Errorf("failed to init did: %s", err)
	}

	// 1. register ID
	id, err := account.GenerateID()
	if err != nil {
		t.Fatal("generate ID error")
	}
	a0 := account.NewAccount("")
	if err := regID(n, id, a0); err != nil {
		t.Fatal("register ID error", err)
	}

	// 2. add key without valid signature, should fail
	a1 := account.NewAccount("")
	sink := common.NewZeroCopySink(nil)
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{common.ADDRESS_EMPTY}
	if _, err = addKey(n); err == nil {
		t.Error("key added without valid signature")
	}

	// 3. add key by invalid owner, should fail
	a2 := account.NewAccount("")
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a2.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a2.Address}
	if _, err = addKey(n); err == nil {
		t.Error("key added by invalid owner")
	}

	// 4. add invalid key, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes([]byte("test invalid key"))
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a0.Address}
	if _, err = addKey(n); err == nil {
		t.Error("invalid key added")
	}

	// 5. add key
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a0.Address}
	if _, err = addKey(n); err != nil {
		t.Fatal(err)
	}

	// 6. verify new key
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 2)
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a1.Address}
	res, err := verifySignature(n)
	if err != nil || !bytes.Equal(res, utils.BYTE_TRUE) {
		t.Fatal("verify the added key failed")
	}

	// 7. add key again, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a0.Address}
	if _, err = addKey(n); err == nil {
		t.Fatal("should not add the same key twice")
	}

	// 8. remove key without valid signature, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a2.Address}
	if _, err = removeKey(n); err == nil {
		t.Error("key removed without valid signature")
	}

	// 9. remove key by invalid owner, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a2.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a2.Address}
	if _, err = removeKey(n); err == nil {
		t.Error("key removed by invalid owner")
	}

	// 10. remove invalid key, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a2.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a1.Address}
	if _, err = removeKey(n); err == nil {
		t.Error("invalid key removed")
	}

	// 11. remove key
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	n.Input = sink.Bytes()
	if _, err = removeKey(n); err != nil {
		t.Fatal(err)
	}

	// 12. check removed key
	sink.Reset()
	sink.WriteString(id)
	utils.EncodeVarUint(sink, 1)
	n.Input = sink.Bytes()
	n.Tx.SignedAddr = []common.Address{a0.Address}
	res, err = verifySignature(n)
	if err == nil && bytes.Equal(res, utils.BYTE_TRUE) {
		t.Fatal("removed key passed verification")
	}

	// 13. add removed key again, should fail
	sink.Reset()
	sink.WriteString(id)
	sink.WriteVarBytes(keypair.SerializePublicKey(a0.PubKey()))
	sink.WriteVarBytes(keypair.SerializePublicKey(a1.PubKey()))
	n.Input = sink.Bytes()
	res, err = verifySignature(n)
	if err == nil && bytes.Equal(res, utils.BYTE_TRUE) {
		t.Error("the removed key should not be added again")
	}

	// 14. query removed key
	sink.Reset()
	sink.WriteString(id)
	sink.WriteInt32(1)
	n.Input = sink.Bytes()
	_, err = GetPublicKeyByID(n)
	if err == nil {
		t.Error("query removed key should fail")
	}
}

func CaseOwnerSize(t *testing.T, n *native.NativeService) {
	sink := common.NewZeroCopySink(nil)
	sink.WriteVarBytes([]byte("dna"))
	n.Input = sink.Bytes()
	if _, err := didInit(n); err != nil {
		t.Errorf("failed to init did: %s", err)
	}

	id, _ := account.GenerateID()
	a := account.NewAccount("")
	err := regID(n, id, a)
	if err != nil {
		t.Fatal(err)
	}

	enc, err := encodeID(common2.DIDContractAddress, []byte(id))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, OWNER_TOTAL_SIZE)
	_, err = insertPk(n, enc, buf)
	if err == nil {
		t.Fatal("total size of the owner's key should be limited")
	}
}

func CaseDoubleInit(t *testing.T, n *native.NativeService) {
	initSink := common.NewZeroCopySink(nil)
	initSink.WriteVarBytes([]byte("dna"))
	n.Input = initSink.Bytes()
	if _, err := didInit(n); err != nil {
		t.Errorf("failed to init did: %s", err)
	}

	initSink = common.NewZeroCopySink(nil)
	initSink.WriteVarBytes([]byte("abc"))
	n.Input = initSink.Bytes()
	if _, err := didInit(n); err == nil {
		t.Error("failed to double init did")
	}
}

func GetPublicKeyByID(srvc *native.NativeService) ([]byte, error) {
	args := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(args)
	if err != nil {
		return nil, errors.New("get public key failed: argument 0 error")
	}
	// arg1: key ID
	arg1, err := utils.DecodeUint32(args)
	if err != nil {
		return nil, errors.New("get public key failed: argument 1 error")
	}

	key, err := encodeID(common2.DIDContractAddress, arg0)
	if err != nil {
		return nil, fmt.Errorf("get public key failed: %s", err)
	}

	pk, err := getPk(srvc, key, arg1)
	if err != nil {
		return nil, fmt.Errorf("get public key failed: %s", err)
	} else if pk == nil {
		return nil, errors.New("get public key failed: not found")
	} else if pk.revoked {
		return nil, errors.New("get public key failed: revoked")
	}

	return pk.key, nil
}
