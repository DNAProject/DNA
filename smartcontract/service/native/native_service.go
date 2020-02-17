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

package native

import (
	"fmt"
	"github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/core/payload"
	"github.com/DNAProject/DNA/core/types"
	"github.com/DNAProject/DNA/errors"
	"github.com/DNAProject/DNA/smartcontract/context"
	"github.com/DNAProject/DNA/smartcontract/event"
	common2 "github.com/DNAProject/DNA/smartcontract/service/native/common"
	"github.com/DNAProject/DNA/smartcontract/states"
	sstates "github.com/DNAProject/DNA/smartcontract/states"
	"github.com/DNAProject/DNA/smartcontract/storage"
)

type (
	Handler         func(native *NativeService) ([]byte, error)
	RegisterService func(native *NativeService)
)

var (
	Contracts = make(map[common.Address]RegisterService)
)

// Native service struct
// Invoke a native smart contract, new a native service
type NativeService struct {
	CacheDB       *storage.CacheDB
	ServiceMap    map[string]Handler
	Notifications []*event.NotifyEventInfo
	InvokeParam   sstates.ContractInvokeParam
	Input         []byte
	Tx            *types.Transaction
	Height        uint32
	Time          uint32
	BlockHash     common.Uint256
	ContextRef    context.ContextRef
}

func (this *NativeService) Register(methodName string, handler Handler) {
	this.ServiceMap[methodName] = handler
}

func (this *NativeService) getBaseContractAddress(addr common.Address) (common.Address, error) {
	if common2.IsNativeContract(addr) {
		return addr, nil
	}

	dep, err := this.CacheDB.GetContract(addr)
	if err != nil {
		return common.ADDRESS_EMPTY, fmt.Errorf("[getNativeContract] get contract context error: %s", err)
	}
	if dep == nil {
		return common.ADDRESS_EMPTY, errors.NewErr("[NativeVmService] deploy code type error!")
	}

	ndc, err := payload.NewNativeDeployCode(dep.GetRawCode())
	if err != nil {
		return common.ADDRESS_EMPTY, fmt.Errorf("[NativeVmService] native deploy code type error: %s", err)
	}

	return ndc.BaseContractAddress, nil
}

func (this *NativeService) Invoke() ([]byte, error) {
	contract := this.InvokeParam

	// update baseContractAddr from deploy code
	baseContractAddr, err := this.getBaseContractAddress(contract.Address)
	if err != nil {
		return common2.BYTE_FALSE, fmt.Errorf("get native contract base address: %s", err)
	}

	services, ok := Contracts[baseContractAddr]
	if !ok {
		return common2.BYTE_FALSE, fmt.Errorf("Native contract address %x haven't been registered.", contract.Address)
	}
	services(this)
	service, ok := this.ServiceMap[contract.Method]
	if !ok {
		return common2.BYTE_FALSE, fmt.Errorf("Native contract %x doesn't support this function %s.",
			contract.Address, contract.Method)
	}
	args := this.Input
	this.Input = contract.Args
	this.ContextRef.PushContext(&context.Context{ContractAddress: contract.Address})
	notifications := this.Notifications
	this.Notifications = []*event.NotifyEventInfo{}
	result, err := service(this)
	if err != nil {
		return result, errors.NewDetailErr(err, errors.ErrNoCode, "[Invoke] Native serivce function execute error!")
	}
	this.ContextRef.PopContext()
	this.ContextRef.PushNotifications(this.Notifications)
	this.Notifications = notifications
	this.Input = args
	return result, nil
}

func (this *NativeService) NativeCall(address common.Address, method string, args []byte) (interface{}, error) {
	c := states.ContractInvokeParam{
		Address: address,
		Method:  method,
		Args:    args,
	}
	this.InvokeParam = c
	return this.Invoke()
}
