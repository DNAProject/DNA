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

package req

import (
	"github.com/DNAProject/DNA/common/log"
	"github.com/DNAProject/DNA/core/types"
	tc "github.com/DNAProject/DNA/txnpool/common"
	"github.com/ontio/ontology-eventbus/actor"
)

var txnPoolPid *actor.PID

func SetTxnPoolPid(txnPid *actor.PID) {
	txnPoolPid = txnPid
}

//add txn to txnpool
func AddTransaction(transaction *types.Transaction) {
	if txnPoolPid == nil {
		log.Error("[p2p]net_server AddTransaction(): txnpool pid is nil")
		return
	}
	txReq := &tc.TxReq{
		Tx:         transaction,
		Sender:     tc.NetSender,
		TxResultCh: nil,
	}
	txnPoolPid.Tell(txReq)
}
