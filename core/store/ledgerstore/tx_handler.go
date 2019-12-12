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

package ledgerstore

import (
	"fmt"
	"math"
	"strconv"

	"github.com/DNAProject/DNA/common"
	sysconfig "github.com/DNAProject/DNA/common/config"
	"github.com/DNAProject/DNA/common/log"
	"github.com/DNAProject/DNA/core/payload"
	"github.com/DNAProject/DNA/core/store"
	scommon "github.com/DNAProject/DNA/core/store/common"
	"github.com/DNAProject/DNA/core/store/overlaydb"
	"github.com/DNAProject/DNA/core/types"
	"github.com/DNAProject/DNA/smartcontract"
	"github.com/DNAProject/DNA/smartcontract/event"
	"github.com/DNAProject/DNA/smartcontract/service/native/global_params"
	"github.com/DNAProject/DNA/smartcontract/service/native/utils"
	"github.com/DNAProject/DNA/smartcontract/service/neovm"
	"github.com/DNAProject/DNA/smartcontract/service/wasmvm"
	"github.com/DNAProject/DNA/smartcontract/storage"
)

//HandleDeployTransaction deal with smart contract deploy transaction
func (self *StateStore) HandleDeployTransaction(store store.LedgerStore, overlay *overlaydb.OverlayDB, gasTable map[string]uint64, cache *storage.CacheDB,
	tx *types.Transaction, block *types.Block, notify *event.ExecuteNotify) error {
	deploy := tx.Payload.(*payload.DeployCode)
	var (
		notifies    []*event.NotifyEventInfo
		gasConsumed uint64
		err         error
	)

	if deploy.VmType() == payload.WASMVM_TYPE {
		_, err = wasmvm.ReadWasmModule(deploy.GetRawCode(), true)
		if err != nil {
			return err
		}
	}

	address := deploy.Address()
	log.Infof("deploy contract address:%s", address.ToHexString())
	// store contract message
	dep, err := cache.GetContract(address)
	if err != nil {
		return err
	}
	if dep == nil {
		cache.PutContract(deploy)
	}
	cache.Commit()

	notify.Notify = append(notify.Notify, notifies...)
	notify.GasConsumed = gasConsumed
	notify.State = event.CONTRACT_STATE_SUCCESS
	return nil
}

//HandleInvokeTransaction deal with smart contract invoke transaction
func (self *StateStore) HandleInvokeTransaction(store store.LedgerStore, overlay *overlaydb.OverlayDB, gasTable map[string]uint64, cache *storage.CacheDB,
	tx *types.Transaction, block *types.Block, notify *event.ExecuteNotify) error {
	invoke := tx.Payload.(*payload.InvokeCode)

	// init smart contract configuration info
	config := &smartcontract.Config{
		Time:      block.Header.Timestamp,
		Height:    block.Header.Height,
		Tx:        tx,
		BlockHash: block.Hash(),
	}

	var (
		costGasLimit      uint64
		availableGasLimit uint64
		err               error
	)

	availableGasLimit = tx.GasLimit

	//init smart contract info
	sc := smartcontract.SmartContract{
		Config:       config,
		CacheDB:      cache,
		Store:        store,
		GasTable:     gasTable,
		Gas:          availableGasLimit,
		WasmExecStep: sysconfig.DEFAULT_WASM_MAX_STEPCOUNT,
		PreExec:      false,
	}

	//start the smart contract executive function
	engine, _ := sc.NewExecuteEngine(invoke.Code, tx.TxType)

	_, err = engine.Invoke()
	if err != nil {
		return err
	}
	costGasLimit = availableGasLimit - sc.Gas

	var notifies []*event.NotifyEventInfo
	notify.Notify = append(notify.Notify, sc.Notifications...)
	notify.Notify = append(notify.Notify, notifies...)
	notify.GasConsumed = costGasLimit
	notify.State = event.CONTRACT_STATE_SUCCESS
	sc.CacheDB.Commit()
	return nil
}

func SaveNotify(eventStore scommon.EventStore, txHash common.Uint256, notify *event.ExecuteNotify) error {
	if !sysconfig.DefConfig.Common.EnableEventLog {
		return nil
	}
	if err := eventStore.SaveEventNotifyByTx(txHash, notify); err != nil {
		return fmt.Errorf("SaveEventNotifyByTx error %s", err)
	}
	event.PushSmartCodeEvent(txHash, 0, event.EVENT_NOTIFY, notify)
	return nil
}

func refreshGlobalParam(config *smartcontract.Config, cache *storage.CacheDB, store store.LedgerStore) error {
	sink := common.NewZeroCopySink(nil)
	utils.EncodeVarUint(sink, uint64(len(neovm.GAS_TABLE_KEYS)))
	for _, value := range neovm.GAS_TABLE_KEYS {
		sink.WriteString(value)
	}

	sc := smartcontract.SmartContract{
		Config:  config,
		CacheDB: cache,
		Store:   store,
		Gas:     math.MaxUint64,
	}

	service, _ := sc.NewNativeService()
	result, err := service.NativeCall(utils.ParamContractAddress, "getGlobalParam", sink.Bytes())
	if err != nil {
		return err
	}
	params := new(global_params.Params)
	if err := params.Deserialization(common.NewZeroCopySource(result.([]byte))); err != nil {
		return fmt.Errorf("deserialize global params error:%s", err)
	}
	neovm.GAS_TABLE.Range(func(key, value interface{}) bool {
		n, ps := params.GetParam(key.(string))
		if n != -1 && ps.Value != "" {
			pu, err := strconv.ParseUint(ps.Value, 10, 64)
			if err != nil {
				log.Errorf("[refreshGlobalParam] failed to parse uint %v\n", ps.Value)
			} else {
				neovm.GAS_TABLE.Store(key, pu)
			}
		}
		return true
	})
	return nil
}

func calcGasByCodeLen(codeLen int, codeGas uint64) uint64 {
	return uint64(codeLen/neovm.PER_UNIT_CODE_LEN) * codeGas
}
