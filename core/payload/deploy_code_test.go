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
package payload

import (
	"bytes"
	"testing"

	"github.com/DNAProject/DNA/common"
	common2 "github.com/DNAProject/DNA/smartcontract/service/native/common"
	"github.com/stretchr/testify/assert"
)

func TestDeployCode_Serialize(t *testing.T) {
	ty, _ := VmTypeFromByte(1)
	deploy, err := NewDeployCode([]byte{1, 2, 3}, ty, "", "", "", "", "")
	assert.Nil(t, err)
	sink := common.NewZeroCopySink(nil)
	deploy.Serialization(sink)
	bs := sink.Bytes()
	var deploy2 DeployCode

	source := common.NewZeroCopySource(bs)
	deploy2.Deserialization(source)
	assert.Equal(t, &deploy2, deploy)

	source = common.NewZeroCopySource(bs[:len(bs)-1])
	err = deploy2.Deserialization(source)
	assert.NotNil(t, err)
}

func TestNativeDeployCode_Serialization(t *testing.T) {
	ndc := &NativeDeployCode{
		BaseContractAddress: common2.GasContractAddress,
		InitParam:           []byte("abcd"),
	}

	sink := common.NewZeroCopySink(nil)
	ndc.Serialization(sink)

	ndc2 := &NativeDeployCode{}
	if err := ndc2.Deserialization(common.NewZeroCopySource(sink.Bytes())); err != nil {
		t.Fatalf("deserialize NativeDeployCode failed: %s", err)
	}
	if ndc.BaseContractAddress != ndc2.BaseContractAddress ||
		bytes.Compare(ndc.InitParam, ndc2.InitParam) != 0 {
		t.Fatalf("deserialize NativeDeployCode unmatched")
	}

	ndc3 := &NativeDeployCode{}
	if err := ndc3.Deserialization(common.NewZeroCopySource(common2.GasContractAddress[:])); err != nil {
		t.Fatalf("deserialize NativeDeployCode from contract addr failed: %s", err)
	}
	if ndc3.BaseContractAddress != common2.GasContractAddress || ndc3.InitParam != nil {
		t.Fatalf("deserialize NativeDeployCode from contract addr not matched")
	}
}
