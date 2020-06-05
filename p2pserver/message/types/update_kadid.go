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

package types

import (
	common2 "github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/p2pserver/common"
)

type UpdatePeerKeyId struct {
	//TODO remove this legecy field when upgrade network layer protocal
	KadKeyId *common.PeerKeyId
}

//Serialize message payload
func (this *UpdatePeerKeyId) Serialization(sink *common2.ZeroCopySink) {
	this.KadKeyId.Serialization(sink)
}

func (this *UpdatePeerKeyId) Deserialization(source *common2.ZeroCopySource) error {
	this.KadKeyId = &common.PeerKeyId{}
	return this.KadKeyId.Deserialization(source)
}

func (this *UpdatePeerKeyId) CmdType() string {
	return common.UPDATE_KADID_TYPE
}
