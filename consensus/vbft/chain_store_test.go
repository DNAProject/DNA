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

package vbft

import (
	"os"
	"testing"

	"github.com/DNAProject/DNA/account"
	"github.com/DNAProject/DNA/common/config"
	"github.com/DNAProject/DNA/common/log"
	"github.com/DNAProject/DNA/consensus/vbft/config"
	"github.com/DNAProject/DNA/core/genesis"
	"github.com/DNAProject/DNA/core/ledger"
	"github.com/ontio/ontology-crypto/keypair"
)

var testBookkeeperAccounts []*account.Account

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
		AdminOntID:           "did:ont:AdjfcJgwru2FD8kotCPvLDXYzRjqFjc9Tb",
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

func newTestChainStore(t *testing.T) *ChainStore {
	log.InitLog(log.InfoLog, log.Stdout)
	var err error
	acct := account.NewAccount("SHA256withECDSA")
	if acct == nil {
		t.Fatalf("GetDefaultAccount error: acc is nil")
	}
	os.RemoveAll(config.DEFAULT_DATA_DIR)
	db, err := ledger.NewLedger(config.DEFAULT_DATA_DIR, 0)
	if err != nil {
		t.Fatalf("NewLedger error %s", err)
	}

	var bookkeepers []keypair.PublicKey
	if len(testBookkeeperAccounts) == 0 {
		for i := 0; i < 7; i++ {
			acc := account.NewAccount("")
			testBookkeeperAccounts = append(testBookkeeperAccounts, acc)
			bookkeepers = append(bookkeepers, acc.PublicKey)
		}
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
		t.Fatalf("BuildGenesisBlock error %s", err)
	}

	err = db.Init(bookkeepers, block)
	if err != nil {
		t.Fatalf("InitLedgerStoreWithGenesisBlock error %s", err)
	}
	chainstore, err := OpenBlockStore(db, nil)
	if err != nil {
		t.Fatalf("openblockstore failed: %v\n", err)
	}
	return chainstore
}

func cleanTestChainStore() {
	os.RemoveAll(config.DEFAULT_DATA_DIR)
	testBookkeeperAccounts = make([]*account.Account, 0)
}

func TestGetChainedBlockNum(t *testing.T) {
	chainstore := newTestChainStore(t)
	if chainstore == nil {
		t.Error("newChainStrore error")
		return
	}
	defer cleanTestChainStore()

	blocknum := chainstore.GetChainedBlockNum()
	t.Logf("TestGetChainedBlockNum :%d", blocknum)
}
