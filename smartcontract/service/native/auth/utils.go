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

package auth

import (
	"fmt"
	"sort"

	"github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/core/states"
	"github.com/DNAProject/DNA/smartcontract/event"
	"github.com/DNAProject/DNA/smartcontract/service/native"
	common2 "github.com/DNAProject/DNA/smartcontract/service/native/common"
	"github.com/DNAProject/DNA/smartcontract/service/native/utils"
)

var (
	PreAdmin           = []byte{0x01}
	PreRoleFunc        = []byte{0x02}
	PreRoleToken       = []byte{0x03}
	PreDelegateStatus  = []byte{0x04}
	PreDIDContractAddr = []byte{0x05}
)

//type(this.contractAddr.Admin) = []byte
func concatContractAdminKey(native *native.NativeService, contractAddr common.Address) []byte {
	this := native.ContextRef.CurrentContext().ContractAddress
	adminKey := append(this[:], contractAddr[:]...)
	adminKey = append(adminKey, PreAdmin...)

	return adminKey
}

func getContractAdmin(native *native.NativeService, contractAddr common.Address) ([]byte, error) {
	key := concatContractAdminKey(native, contractAddr)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil { //is not set
		return nil, nil
	}
	return item.Value, nil
}

func putContractAdmin(native *native.NativeService, contractAddr common.Address, adminDID []byte) error {
	key := concatContractAdminKey(native, contractAddr)
	utils.PutBytes(native, key, adminDID)
	return nil
}

//type(this.contractAddr.RoleFunc.role) = roleFuncs
func concatRoleFuncKey(native *native.NativeService, contractAddr common.Address, role []byte) []byte {
	this := native.ContextRef.CurrentContext().ContractAddress
	roleFuncKey := append(this[:], contractAddr[:]...)
	roleFuncKey = append(roleFuncKey, PreRoleFunc...)
	roleFuncKey = append(roleFuncKey, role...)

	return roleFuncKey
}

func getRoleFunc(native *native.NativeService, contractAddr common.Address, role []byte) (*roleFuncs, error) {
	key := concatRoleFuncKey(native, contractAddr, role)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil { //is not set
		return nil, nil
	}
	source := common.NewZeroCopySource(item.Value)
	rF := new(roleFuncs)
	err = rF.Deserialization(source)
	if err != nil {
		return nil, fmt.Errorf("deserialize roleFuncs object failed. data: %x", item.Value)
	}
	return rF, nil
}

func putRoleFunc(native *native.NativeService, contractAddr common.Address, role []byte, funcs *roleFuncs) error {
	key := concatRoleFuncKey(native, contractAddr, role)
	utils.PutBytes(native, key, common.SerializeToBytes(funcs))
	return nil
}

//type(this.contractAddr.RoleP.ontID) = roleTokens
func concatDIDTokenKey(native *native.NativeService, contractAddr common.Address, DID []byte) []byte {
	this := native.ContextRef.CurrentContext().ContractAddress
	tokenKey := append(this[:], contractAddr[:]...)
	tokenKey = append(tokenKey, PreRoleToken...)
	tokenKey = append(tokenKey, DID...)

	return tokenKey
}

func getDIDToken(native *native.NativeService, contractAddr common.Address, DID []byte) (*roleTokens, error) {
	key := concatDIDTokenKey(native, contractAddr, DID)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil { //is not set
		return nil, nil
	}
	source := common.NewZeroCopySource(item.Value)
	rT := new(roleTokens)
	err = rT.Deserialization(source)
	if err != nil {
		return nil, fmt.Errorf("deserialize roleTokens object failed. data: %x", item.Value)
	}
	return rT, nil
}

func putDIDToken(native *native.NativeService, contractAddr common.Address, DID []byte, tokens *roleTokens) error {
	key := concatDIDTokenKey(native, contractAddr, DID)
	utils.PutBytes(native, key, common.SerializeToBytes(tokens))
	return nil
}

//type(this.contractAddr.DelegateStatus.ontID)
func concatDelegateStatusKey(native *native.NativeService, contractAddr common.Address, DID []byte) []byte {
	this := native.ContextRef.CurrentContext().ContractAddress
	key := append(this[:], contractAddr[:]...)
	key = append(key, PreDelegateStatus...)
	key = append(key, DID...)

	return key
}

func getDelegateStatus(native *native.NativeService, contractAddr common.Address, DID []byte) (*Status, error) {
	key := concatDelegateStatusKey(native, contractAddr, DID)
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	status := new(Status)
	source := common.NewZeroCopySource(item.Value)
	err = status.Deserialization(source)
	if err != nil {
		return nil, fmt.Errorf("deserialize Status object failed. data: %x", item.Value)
	}
	return status, nil
}

func putDelegateStatus(native *native.NativeService, contractAddr common.Address, DID []byte, status *Status) error {
	key := concatDelegateStatusKey(native, contractAddr, DID)
	utils.PutBytes(native, key, common.SerializeToBytes(status))
	return nil
}

//remove duplicates in the slice of string and sorts the slice in increasing order.
func StringsDedupAndSort(s []string) []string {
	smap := make(map[string]int)
	for i, str := range s {
		if str == "" {
			continue
		}
		smap[str] = i
	}
	ret := make([]string, len(smap))
	i := 0
	for str := range smap {
		ret[i] = str
		i++
	}
	sort.Strings(ret)
	return ret
}

func pushEvent(native *native.NativeService, s interface{}) {
	event := new(event.NotifyEventInfo)
	event.ContractAddress = native.ContextRef.CurrentContext().ContractAddress
	event.States = s
	native.Notifications = append(native.Notifications, event)
}

func serializeAddress(sink *common.ZeroCopySink, addr common.Address) {
	sink.WriteVarBytes(addr[:])
}

func getDIDContractAddr(native *native.NativeService) (common.Address, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	if contract == common2.AuthContractAddress {
		// default auth contract
		return common2.DIDContractAddress, nil
	}
	didContractAddrBytes, err := native.CacheDB.Get(utils.ConcatKey(contract, PreDIDContractAddr))
	if err != nil {
		return common.ADDRESS_EMPTY, fmt.Errorf("get did contract: %s", err)
	}
	value, err := states.GetValueFromRawStorageItem(didContractAddrBytes)
	if err != nil {
		return common.ADDRESS_EMPTY, fmt.Errorf("parse did contract: %s", err)
	}
	didContractAddr, err := common.AddressParseFromBytes(value)
	if err != nil {
		return common.ADDRESS_EMPTY, fmt.Errorf("invalid did contract addr: %s", err)
	}
	return didContractAddr, nil
}

func getDIDMethod(native *native.NativeService) ([]byte, error) {
	didContractAddr, err := getDIDContractAddr(native)
	if err != nil {
		return nil, fmt.Errorf("invalid did contract addr: %s", err)
	}
	didMethod, err := native.NativeCall(didContractAddr, "getDIDMethod", []byte{})
	if err != nil {
		return nil, fmt.Errorf("get did method from contract: %s", err)
	}
	if result, ok := didMethod.([]byte); ok {
		return result, nil
	}
	return nil, fmt.Errorf("get did method from contract failed")
}
