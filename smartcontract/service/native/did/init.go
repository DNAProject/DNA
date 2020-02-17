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
	"fmt"

	common2 "github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/core/states"
	"github.com/DNAProject/DNA/smartcontract/service/native"
	"github.com/DNAProject/DNA/smartcontract/service/native/common"
	"github.com/DNAProject/DNA/smartcontract/service/native/utils"
)

func Init() {
	native.Contracts[common.DIDContractAddress] = RegisterIDContract
}

func RegisterIDContract(srvc *native.NativeService) {
	srvc.Register("initDID", didInit)
	srvc.Register("getDIDMethod", getDIDMethod)
	srvc.Register("regIDWithPublicKey", regIdWithPublicKey)
	srvc.Register("regIDWithController", regIdWithController)
	srvc.Register("revokeID", revokeID)
	srvc.Register("revokeIDByController", revokeIDByController)
	srvc.Register("removeController", removeController)
	srvc.Register("addRecovery", addRecovery)
	srvc.Register("changeRecovery", changeRecovery)
	srvc.Register("setRecovery", setRecovery)
	srvc.Register("updateRecovery", updateRecovery)
	srvc.Register("addKey", addKey)
	srvc.Register("removeKey", removeKey)
	srvc.Register("addKeyByController", addKeyByController)
	srvc.Register("removeKeyByController", removeKeyByController)
	srvc.Register("addKeyByRecovery", addKeyByRecovery)
	srvc.Register("removeKeyByRecovery", removeKeyByRecovery)
	srvc.Register("regIDWithAttributes", regIdWithAttributes)
	srvc.Register("addAttributes", addAttributes)
	srvc.Register("removeAttribute", removeAttribute)
	srvc.Register("addAttributesByController", addAttributesByController)
	srvc.Register("removeAttributeByController", removeAttributeByController)
	srvc.Register("verifySignature", verifySignature)
	srvc.Register("verifyController", verifyController)
	srvc.Register("getPublicKeys", GetPublicKeys)
	srvc.Register("getKeyState", GetKeyState)
	srvc.Register("getAttributes", GetAttributes)
	srvc.Register("getDDO", GetDDO)
	return
}

func didInit(srvc *native.NativeService) ([]byte, error) {
	didMethod, err := utils.DecodeVarBytes(common2.NewZeroCopySource(srvc.Input))
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("init did, contract param deserialize err: %s", err)
	}
	if len(didMethod) != 3 {
		return utils.BYTE_FALSE, fmt.Errorf("init did, invalid length of did-method: %s", string(didMethod))
	}

	// check if has initialized
	contract := srvc.ContextRef.CurrentContext().ContractAddress
	didMethodBytes, err := srvc.CacheDB.Get(utils.ConcatKey(contract, []byte{FIELD_DID_METHOD}))
	if didMethodBytes != nil || err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("init did, already inited")
	}

	// save did-method
	srvc.CacheDB.Put(utils.ConcatKey(contract, []byte{FIELD_DID_METHOD}), states.GenRawStorageItem(didMethod))
	return utils.BYTE_TRUE, nil
}

func getDIDMethod(srvc *native.NativeService) ([]byte, error) {
	contract := srvc.ContextRef.CurrentContext().ContractAddress
	if contract == common.DIDContractAddress {
		// default did contract
		return []byte("dna"), nil
	}

	didMethodBytes, err := srvc.CacheDB.Get(utils.ConcatKey(contract, []byte{FIELD_DID_METHOD}))
	if err != nil {
		return nil, err
	}
	return states.GetValueFromRawStorageItem(didMethodBytes)
}
