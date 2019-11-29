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
package testsuite

import (
	"testing"

	"github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/smartcontract/service/native"
	"github.com/DNAProject/DNA/smartcontract/service/native/gas"
	_ "github.com/DNAProject/DNA/smartcontract/service/native/init"
	"github.com/DNAProject/DNA/smartcontract/service/native/utils"
	"github.com/DNAProject/DNA/smartcontract/storage"
	"github.com/stretchr/testify/assert"
)

func setBalance(db *storage.CacheDB, addr common.Address, value uint64) {
	balanceKey := gas.GenBalanceKey(utils.GasContractAddress, addr)
	item := utils.GenUInt64StorageItem(value)
	db.Put(balanceKey, item.ToArray())
}

func getBalanceOf(native *native.NativeService, addr common.Address) int {
	sink := common.NewZeroCopySink(nil)
	utils.EncodeAddress(sink, addr)
	native.Input = sink.Bytes()
	buf, _ := gas.GasBalanceOf(native)
	val := common.BigIntFromNeoBytes(buf)
	return int(val.Uint64())
}

func makeTransfer(native *native.NativeService, from, to common.Address, value uint64) error {
	native.Tx.SignedAddr = append(native.Tx.SignedAddr, from)

	state := gas.State{from, to, value}
	native.Input = common.SerializeToBytes(&gas.Transfers{States: []gas.State{state}})

	_, err := gas.GasTransfer(native)
	return err
}

func TestTransfer(t *testing.T) {
	InvokeNativeContract(t, utils.GasContractAddress, func(native *native.NativeService) ([]byte, error) {
		a := RandomAddress()
		b := RandomAddress()
		c := RandomAddress()
		setBalance(native.CacheDB, a, 10000)

		assert.Equal(t, getBalanceOf(native, a), 10000)
		assert.Equal(t, getBalanceOf(native, b), 0)
		assert.Equal(t, getBalanceOf(native, c), 0)

		assert.Nil(t, makeTransfer(native, a, b, 10))
		assert.Equal(t, getBalanceOf(native, a), 9990)
		assert.Equal(t, getBalanceOf(native, b), 10)

		assert.Nil(t, makeTransfer(native, b, c, 10))
		assert.Equal(t, getBalanceOf(native, b), 0)
		assert.Equal(t, getBalanceOf(native, c), 10)

		return nil, nil
	})
}
