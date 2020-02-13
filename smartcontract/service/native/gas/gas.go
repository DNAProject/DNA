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

package gas

import (
	"fmt"
	"math/big"

	"github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/common/constants"
	"github.com/DNAProject/DNA/errors"
	"github.com/DNAProject/DNA/smartcontract/service/native"
	common2 "github.com/DNAProject/DNA/smartcontract/service/native/common"
	"github.com/DNAProject/DNA/smartcontract/service/native/utils"
)

const (
	TRANSFER_FLAG byte = 1
	APPROVE_FLAG  byte = 2
)

func InitGas() {
	native.Contracts[common2.GasContractAddress] = RegisterGasContract
}

func RegisterGasContract(native *native.NativeService) {
	native.Register(INIT_NAME, GasInit)
	native.Register(TRANSFER_NAME, GasTransfer)
	native.Register(APPROVE_NAME, GasApprove)
	native.Register(TRANSFERFROM_NAME, GasTransferFrom)
	native.Register(NAME_NAME, GasName)
	native.Register(SYMBOL_NAME, GasSymbol)
	native.Register(DECIMALS_NAME, GasDecimals)
	native.Register(TOTALSUPPLY_NAME, GasTotalSupply)
	native.Register(BALANCEOF_NAME, GasBalanceOf)
	native.Register(ALLOWANCE_NAME, GasAllowance)
}

func GasInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init has been completed!")
	}

	distribute := make(map[common.Address]uint64)
	source := common.NewZeroCopySource(native.Input)
	buf, _, irregular, eof := source.NextVarBytes()
	if eof {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "serialization.ReadVarBytes, contract params deserialize error!")
	}
	if irregular {
		return utils.BYTE_FALSE, common.ErrIrregularData
	}
	input := common.NewZeroCopySource(buf)
	num, err := utils.DecodeVarUint(input)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("read number error:%v", err)
	}
	sum := uint64(0)
	overflow := false
	for i := uint64(0); i < num; i++ {
		addr, err := utils.DecodeAddress(input)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("read address error:%v", err)
		}
		value, err := utils.DecodeVarUint(input)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("read value error:%v", err)
		}
		sum, overflow = common.SafeAdd(sum, value)
		if overflow {
			return utils.BYTE_FALSE, errors.NewErr("wrong config. overflow detected")
		}
		distribute[addr] += value
	}
	if sum != constants.GAS_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("wrong config. total supply %d != %d", sum, constants.GAS_TOTAL_SUPPLY)
	}

	for addr, val := range distribute {
		balanceKey := GenBalanceKey(contract, addr)
		item := utils.GenUInt64StorageItem(val)
		native.CacheDB.Put(balanceKey, item.ToArray())
		AddNotifications(native, contract, &State{To: addr, Value: val})
	}
	native.CacheDB.Put(GenTotalSupplyKey(contract), utils.GenUInt64StorageItem(constants.GAS_TOTAL_SUPPLY).ToArray())

	return utils.BYTE_TRUE, nil
}

func GasTransfer(native *native.NativeService) ([]byte, error) {
	var transfers Transfers
	source := common.NewZeroCopySource(native.Input)
	if err := transfers.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[Transfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if v.Value == 0 {
			continue
		}
		if v.Value > constants.GAS_TOTAL_SUPPLY {
			return utils.BYTE_FALSE, fmt.Errorf("transfer ont amount:%d over totalSupply:%d", v.Value, constants.GAS_TOTAL_SUPPLY)
		}
		if _, _, err := Transfer(native, contract, &v); err != nil {
			return utils.BYTE_FALSE, err
		}

		AddNotifications(native, contract, &v)
	}
	return utils.BYTE_TRUE, nil
}

func GasTransferFrom(native *native.NativeService) ([]byte, error) {
	var state TransferFrom
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GasTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.GAS_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("transferFrom amount:%d over totalSupply:%d", state.Value, constants.GAS_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	_, _, err := TransferedFrom(native, contract, &state)
	if err != nil {
		return utils.BYTE_FALSE, err
	}
	AddNotifications(native, contract, &State{From: state.From, To: state.To, Value: state.Value})
	return utils.BYTE_TRUE, nil
}

func GasApprove(native *native.NativeService) ([]byte, error) {
	var state State
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OngApprove] state deserialize error!")
	}
	if state.Value > constants.GAS_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve amount:%d over totalSupply:%d", state.Value, constants.GAS_TOTAL_SUPPLY)
	}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("authentication failed!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CacheDB.Put(GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value).ToArray())
	return utils.BYTE_TRUE, nil
}

func GasName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.GAS_NAME), nil
}

func GasDecimals(native *native.NativeService) ([]byte, error) {
	return common.BigIntToNeoBytes(big.NewInt(int64(constants.GAS_DECIMALS))), nil
}

func GasSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.GAS_SYMBOL), nil
}

func GasTotalSupply(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[TotalSupply] get totalSupply error!")
	}
	return common.BigIntToNeoBytes(big.NewInt(int64(amount))), nil
}

func GasBalanceOf(native *native.NativeService) ([]byte, error) {
	return GetBalanceValue(native, TRANSFER_FLAG)
}

func GasAllowance(native *native.NativeService) ([]byte, error) {
	return GetBalanceValue(native, APPROVE_FLAG)
}

func GetBalanceValue(native *native.NativeService, flag byte) ([]byte, error) {
	source := common.NewZeroCopySource(native.Input)
	from, err := utils.DecodeAddress(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] get from address error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	var key []byte
	if flag == APPROVE_FLAG {
		to, err := utils.DecodeAddress(source)
		if err != nil {
			return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] get from address error!")
		}
		key = GenApproveKey(contract, from, to)
	} else if flag == TRANSFER_FLAG {
		key = GenBalanceKey(contract, from)
	}
	amount, err := utils.GetStorageUInt64(native, key)
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[GetBalanceValue] address parse error!")
	}
	return common.BigIntToNeoBytes(big.NewInt(int64(amount))), nil
}

func getApproveArgs(native *native.NativeService, contract, ongContract, address common.Address, value uint64) ([]byte, error) {
	bf := common.NewZeroCopySink(nil)
	approve := State{
		From:  contract,
		To:    address,
		Value: value,
	}

	stateValue, err := utils.GetStorageUInt64(native, GenApproveKey(ongContract, approve.From, approve.To))
	if err != nil {
		return nil, err
	}

	approve.Value += stateValue
	approve.Serialization(bf)
	return bf.Bytes(), nil
}
