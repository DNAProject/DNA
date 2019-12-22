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

package txnpool

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DNAProject/DNA/account"
	"github.com/DNAProject/DNA/common/config"
	"github.com/DNAProject/DNA/common/log"
	"github.com/DNAProject/DNA/consensus/vbft/config"
	"github.com/DNAProject/DNA/core/genesis"
	"github.com/DNAProject/DNA/core/ledger"
	"github.com/DNAProject/DNA/core/payload"
	"github.com/DNAProject/DNA/core/types"
	tc "github.com/DNAProject/DNA/txnpool/common"
	tp "github.com/DNAProject/DNA/txnpool/proc"
	"github.com/DNAProject/DNA/validator/stateful"
	"github.com/DNAProject/DNA/validator/stateless"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology-eventbus/actor"
)

var testGenesisConfig = &config.GenesisConfig{
	SeedList: []string{
		"localhost:20338",
		"localhost:20438",
		"localhost:20538",
		"localhost:20638",
		"localhost:20738"},
	ConsensusType: config.CONSENSUS_TYPE_VBFT,
	VBFT: &config.VBFTConfig{
		N:                    7,
		C:                    2,
		K:                    7,
		L:                    112,
		BlockMsgDelay:        10000,
		HashMsgDelay:         10000,
		PeerHandshakeTimeout: 10,
		MaxBlockChangeView:   120000,
		AdminOntID:           "did:dna:AdjfcJgwru2FD8kotCPvLDXYzRjqFjc9Tb",
		MinInitStake:         100000,
		VrfValue:             "",
		VrfProof:             "",
		Peers: []*config.VBFTPeerStakeInfo{
			{Index: 1},
			{Index: 2},
			{Index: 3},
			{Index: 4},
			{Index: 5},
			{Index: 6},
			{Index: 7},
		},
	},
	DBFT: &config.DBFTConfig{},
	SOLO: &config.SOLOConfig{},
}

func TestMain(m *testing.M) {
	log.InitLog(log.InfoLog, log.Stdout)

	var err error
	ledger.DefLedger, err = ledger.NewLedger(config.DEFAULT_DATA_DIR, 0)
	if err != nil {
		log.Errorf("failed  to new ledger")
		return
	}

	m.Run()

	// tear down
	ledger.DefLedger.Close()
	os.RemoveAll(config.DEFAULT_DATA_DIR)
}

func initTestTx() *types.Transaction {
	log.InitLog(log.InfoLog, log.Stdout)
	//topic := "TXN"

	mutable := &types.MutableTransaction{
		TxType:  types.InvokeNeo,
		Nonce:   uint32(time.Now().Unix()),
		Payload: &payload.InvokeCode{Code: []byte{}},
	}

	tx, _ := mutable.IntoImmutable()
	return tx
}

func startActor(obj interface{}) *actor.PID {
	props := actor.FromProducer(func() actor.Actor {
		return obj.(actor.Actor)
	})

	pid := actor.Spawn(props)
	return pid
}

func Test_RCV(t *testing.T) {
	var s *tp.TXPoolServer
	var wg sync.WaitGroup

	var bookkeepers []keypair.PublicKey
	testBookkeeperAccounts := make([]*account.Account, 0)
	for i := 0; i < 7; i++ {
		acc := account.NewAccount("")
		testBookkeeperAccounts = append(testBookkeeperAccounts, acc)
		bookkeepers = append(bookkeepers, acc.PublicKey)
	}

	config.DefConfig.Genesis = testGenesisConfig
	genesisConfig := config.DefConfig.Genesis

	// update peers in genesis
	for i, p := range genesisConfig.VBFT.Peers {
		if i < len(testBookkeeperAccounts) {
			p.PeerPubkey = vconfig.PubkeyID(testBookkeeperAccounts[i].PublicKey)
			p.Address = testBookkeeperAccounts[i].Address.ToBase58()
		}
	}

	block, err := genesis.BuildGenesisBlock(bookkeepers, genesisConfig)
	if err != nil {
		log.Fatalf("failed to build genesis block: %s", err)
	}
	err = ledger.DefLedger.Init(bookkeepers, block)
	if err != nil {
		t.Error("failed to initialize default ledger", err)
		return
	}

	// Start txnpool server to receive msgs from p2p, consensus and valdiators
	s = tp.NewTxPoolServer(tc.MAX_WORKER_NUM, true, false)

	// Initialize an actor to handle the msgs from valdiators
	rspActor := tp.NewVerifyRspActor(s)
	rspPid := startActor(rspActor)
	if rspPid == nil {
		t.Error("Fail to start verify rsp actor")
		return
	}
	s.RegisterActor(tc.VerifyRspActor, rspPid)

	// Initialize an actor to handle the msgs from consensus
	tpa := tp.NewTxPoolActor(s)
	txPoolPid := startActor(tpa)
	if txPoolPid == nil {
		t.Error("Fail to start txnpool actor")
		return
	}
	s.RegisterActor(tc.TxPoolActor, txPoolPid)

	// Initialize an actor to handle the msgs from p2p and api
	ta := tp.NewTxActor(s)
	txPid := startActor(ta)
	if txPid == nil {
		t.Error("Fail to start txn actor")
		return
	}
	s.RegisterActor(tc.TxActor, txPid)

	// Start stateless validator
	statelessV, err := stateless.NewValidator("stateless")
	if err != nil {
		t.Errorf("failed to new stateless valdiator: %s", err)
		return
	}
	statelessV.Register(rspPid)

	statelessV2, err := stateless.NewValidator("stateless2")
	if err != nil {
		t.Errorf("failed to new stateless valdiator: %s", err)
		return
	}
	statelessV2.Register(rspPid)

	statelessV3, err := stateless.NewValidator("stateless3")
	if err != nil {
		t.Errorf("failed to new stateless valdiator: %s", err)
		return
	}
	statelessV3.Register(rspPid)

	statefulV, err := stateful.NewValidator("stateful")
	if err != nil {
		t.Errorf("failed to new stateful valdiator: %s", err)
		return
	}
	statefulV.Register(rspPid)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			var j int
			defer wg.Done()

			tx := initTestTx()
			for {
				j++
				txReq := &tc.TxReq{
					Tx:     tx,
					Sender: tc.NilSender,
				}
				txPid.Tell(txReq)

				if j >= 4 {
					return
				}
			}
		}()
	}

	wg.Wait()
	time.Sleep(1 * time.Second)
	txPoolPid.Tell(&tc.GetTxnPoolReq{ByCount: true})
	txPoolPid.Tell(&tc.GetPendingTxnReq{ByCount: true})
	time.Sleep(2 * time.Second)

	statelessV.UnRegister(rspPid)
	statelessV2.UnRegister(rspPid)
	statelessV3.UnRegister(rspPid)
	statefulV.UnRegister(rspPid)
	s.Stop()
}
