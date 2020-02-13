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

// Package actor privides communication with other actor
package actor

import (
	"github.com/DNAProject/DNA/common"
	common2 "github.com/DNAProject/DNA/smartcontract/service/native/common"
)

func updateNativeSCAddr(hash common.Address) common.Address {
	if hash == common2.GasContractAddress {
		hash = common.AddressFromVmCode(common2.GasContractAddress[:])
	} else if hash == common2.DIDContractAddress {
		hash = common.AddressFromVmCode(common2.DIDContractAddress[:])
	} else if hash == common2.ParamContractAddress {
		hash = common.AddressFromVmCode(common2.ParamContractAddress[:])
	} else if hash == common2.AuthContractAddress {
		hash = common.AddressFromVmCode(common2.AuthContractAddress[:])
	} else if hash == common2.GovernanceContractAddress {
		hash = common.AddressFromVmCode(common2.GovernanceContractAddress[:])
	}
	return hash
}
