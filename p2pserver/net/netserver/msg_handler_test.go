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

package netserver

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/DNAProject/DNA/account"
	"github.com/DNAProject/DNA/consensus/vbft/config"
	"github.com/ontio/ontology-crypto/keypair"
	"net"
	"os"
	"testing"
	"time"

	"github.com/DNAProject/DNA/common"
	"github.com/DNAProject/DNA/common/config"
	"github.com/DNAProject/DNA/common/log"
	"github.com/DNAProject/DNA/core/genesis"
	"github.com/DNAProject/DNA/core/ledger"
	"github.com/DNAProject/DNA/events"
	msgCommon "github.com/DNAProject/DNA/p2pserver/common"
	"github.com/DNAProject/DNA/p2pserver/dht/kbucket"
	msgpack "github.com/DNAProject/DNA/p2pserver/message/msg_pack"
	"github.com/DNAProject/DNA/p2pserver/message/types"
	p2p "github.com/DNAProject/DNA/p2pserver/net/protocol"
	"github.com/DNAProject/DNA/p2pserver/peer"
	"github.com/DNAProject/DNA/p2pserver/protocols"
	"github.com/stretchr/testify/assert"
)

var (
	network *MockP2P
)

type MockP2P struct {
	p2p.P2P
	SentMsgs []types.Message // stores all mock msgs
}

func (mock *MockP2P) Send(p *peer.Peer, msg types.Message) error {
	mock.SentMsgs = append(mock.SentMsgs, msg)
	return nil
}

func NewMockP2p() *MockP2P {
	return &MockP2P{NewNetServer(), make([]types.Message, 0)}
}

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
	// Start local network server and create message router
	network = NewMockP2p()

	events.Init()
	// Initial a ledger
	var err error
	ledger.DefLedger, err = ledger.NewLedger(config.DEFAULT_DATA_DIR, 0)
	if err != nil {
		log.Fatalf("NewLedger error %s", err)
	}

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
		log.Fatalf("DefLedger.Init error %s", err)
	}

	m.Run()

	_ = ledger.DefLedger.Close()
	_ = os.RemoveAll(config.DEFAULT_DATA_DIR)
}

// TestAddrReqHandle tests Function AddrReqHandle handling an address req
// testcase: no-mask neighbor
func TestAddrReqHandle(t *testing.T) {
	network = NewMockP2p()

	// Simulate a remote peer to be added to the neighbor peers
	testID := kbucket.PseudoKadIdFromUint64(0x7533345)

	remotePeer := peer.NewPeer()

	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:1234")

	network.AddNbrNode(remotePeer)
	remotePeer.SetState(msgCommon.ESTABLISH)

	// Construct an address request packet
	buf := msgpack.NewAddrReq()

	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "test",
		Payload: buf,
	}

	// Invoke AddrReqHandle to handle the msg
	ctx := newContext(t, msg, network)
	AddrReqHandle(ctx)

	// all neighbor peers should be in rsp msg
	for _, msg := range network.SentMsgs {
		addrMsg, ok := msg.(*types.Addr)
		if !ok {
			t.Fatalf("invalid addr msg %s", msg.CmdType())
		}
		if len(addrMsg.NodeAddrs) != 1 {
			t.Fatalf("invalid addr count: %v", addrMsg.NodeAddrs)
		}
		var ip net.IP
		ip = addrMsg.NodeAddrs[0].IpAddr[:]
		addr := fmt.Sprintf("%v:%d", ip, addrMsg.NodeAddrs[0].Port)
		if addr != remotePeer.Link.GetAddr() {
			t.Fatalf("invalid addr: %s vs %s", addr, remotePeer.Link.GetAddr())
		}
	}

	network.DelNbrNode(testID.ToUint64())
}

// create two neighbors, one masked, one un-masked
// send addr-req from un-mask peer, get itself in addr-rsp
func TestAddrReqHandle_maskok(t *testing.T) {
	network = NewMockP2p()

	// Simulate a remote peer to be added to the neighbor peers
	testID := kbucket.PseudoKadIdFromUint64(123456)
	info := peer.NewPeerInfo(testID, 1, 12345678, true, 0, 20336, 12345, "1.5.2", "1.2.3.4:5001")
	remotePeer := peer.NewPeer()
	remotePeer.SetInfo(info)
	remotePeer.Link.SetAddr("1.2.3.4:5001")
	network.AddNbrNode(remotePeer)
	remotePeer.SetState(msgCommon.ESTABLISH)

	testID2 := kbucket.PseudoKadIdFromUint64(1234567)
	info2 := peer.NewPeerInfo(testID2, 1, 12345678, true, 0,
		20336, 12345, "1.5.2", "1.2.3.5:5002")
	remotePeer2 := peer.NewPeer()
	remotePeer2.SetInfo(info2)
	remotePeer2.Link.SetAddr("1.2.3.5:5002")
	network.AddNbrNode(remotePeer2)
	remotePeer2.SetState(msgCommon.ESTABLISH)

	// Construct an address request packet
	buf := msgpack.NewAddrReq()

	msg := &types.MsgPayload{
		Id:      testID2.ToUint64(),
		Addr:    "test",
		Payload: buf,
	}

	config.DefConfig.P2PNode.ReservedPeersOnly = true
	config.DefConfig.P2PNode.ReservedCfg.MaskPeers = []string{"1.2.3.4"}
	// Invoke AddrReqHandle to handle the msg
	ctx := newContext(t, msg, network)
	AddrReqHandle(ctx)

	// verify 1.2.3.4 is masked
	for _, msg := range network.SentMsgs {
		addrMsg, ok := msg.(*types.Addr)
		if !ok {
			t.Fatalf("invalid addr msg %s", msg.CmdType())
		}
		if len(addrMsg.NodeAddrs) != 1 {
			t.Fatalf("invalid addr count: %v", addrMsg.NodeAddrs)
		}
		var ip net.IP
		ip = addrMsg.NodeAddrs[0].IpAddr[:]
		addr := fmt.Sprintf("%v:%d", ip, addrMsg.NodeAddrs[0].Port)
		if addr != remotePeer2.Link.GetAddr() {
			t.Fatalf("invalid addr: %s vs %s", addr, remotePeer2.Link.GetAddr())
		}
	}

	network.DelNbrNode(testID.ToUint64())
}

func newContext(t *testing.T, msg *types.MsgPayload, n p2p.P2P) *protocols.Context {
	sender := n.GetPeer(msg.Id)
	assert.NotNil(t, sender)

	return protocols.NewContext(sender, n, nil, msg.PayloadSize)
}

// create one masked neighbor
// send addr-req, get itself in addr-rsp
func TestAddrReqHandle_unmaskok(t *testing.T) {
	network = NewMockP2p()

	// Simulate a remote peer to be added to the neighbor peers
	testID := kbucket.PseudoKadIdFromUint64(123456)

	remotePeer := peer.NewPeer()

	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("1.2.3.4:5001")

	network.AddNbrNode(remotePeer)
	remotePeer.SetState(msgCommon.ESTABLISH)

	// Construct an address request packet
	buf := msgpack.NewAddrReq()

	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "test",
		Payload: buf,
	}

	config.DefConfig.P2PNode.ReservedPeersOnly = true
	config.DefConfig.P2PNode.ReservedCfg.MaskPeers = []string{"1.2.3.4"}

	// Invoke AddrReqHandle to handle the msg
	ctx := newContext(t, msg, network)
	AddrReqHandle(ctx)

	for _, msg := range network.SentMsgs {
		addrMsg, ok := msg.(*types.Addr)
		if !ok {
			t.Fatalf("invalid addr msg %s", msg.CmdType())
		}
		if len(addrMsg.NodeAddrs) != 1 {
			t.Fatalf("invalid addr count: %v", addrMsg.NodeAddrs)
		}
		var ip net.IP
		ip = addrMsg.NodeAddrs[0].IpAddr[:]
		addr := fmt.Sprintf("%v:%d", ip, addrMsg.NodeAddrs[0].Port)
		if addr != remotePeer.Link.GetAddr() {
			t.Fatalf("invalid addr: %s vs %s", addr, remotePeer.Link.GetAddr())
		}
	}

	network.DelNbrNode(testID.ToUint64())
}

// TestHeadersReqHandle tests Function HeadersReqHandle handling a header req
func TestHeadersReqHandle(t *testing.T) {
	// Simulate a remote peer to be added to the neighbor peers
	testID := kbucket.PseudoKadIdFromUint64(123456)

	remotePeer := peer.NewPeer()
	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:50010")

	network.AddNbrNode(remotePeer)

	// Construct a headers request of packet
	headerHash := ledger.DefLedger.GetCurrentHeaderHash()
	buf := msgpack.NewHeadersReq(headerHash)

	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: buf,
	}

	// Invoke HeadersReqhandle to handle the msg
	ctx := newContext(t, msg, network)
	HeadersReqHandle(ctx, buf.(*types.HeadersReq))
	network.DelNbrNode(testID.ToUint64())
}

// TestPingHandle tests Function PingHandle handling a ping message
func TestPingHandle(t *testing.T) {
	// Simulate a remote peer to be added to the neighbor peers
	testID := kbucket.PseudoKadIdFromUint64(123456)

	remotePeer := peer.NewPeer()
	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:50010")

	network.AddNbrNode(remotePeer)

	// Construct a ping packet
	height := ledger.DefLedger.GetCurrentBlockHeight()

	buf := msgpack.NewPingMsg(uint64(height))

	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: buf,
	}

	// Invoke PingHandle to handle the msg
	ctx := newContext(t, msg, network)
	PingHandle(ctx, buf)

	network.DelNbrNode(testID.ToUint64())
}

// TestPingHandle tests Function PingHandle handling a pong message
func TestPongHandle(t *testing.T) {
	// Simulate a remote peer to be added to the neighbor peers
	testID := kbucket.PseudoKadIdFromUint64(123456)

	remotePeer := peer.NewPeer()
	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:50010")

	network.AddNbrNode(remotePeer)

	// Construct a pong packet
	height := ledger.DefLedger.GetCurrentBlockHeight()

	buf := msgpack.NewPongMsg(uint64(height))

	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: buf,
	}

	// Invoke PingHandle to handle the msg
	ctx := newContext(t, msg, network)
	PongHandle(ctx, buf)

	network.DelNbrNode(testID.ToUint64())
}

// TestBlkHeaderHandle tests Function BlkHeaderHandle handling a sync header msg
func TestBlkHeaderHandle(t *testing.T) {
	// Simulate a remote peer to be added to the neighbor peers
	testID := kbucket.PseudoKadIdFromUint64(123456)

	remotePeer := peer.NewPeer()
	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336, testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:50010")

	network.AddNbrNode(remotePeer)

	// Construct a sync header packet
	hash := ledger.DefLedger.GetBlockHash(0)
	assert.NotEqual(t, hash, common.UINT256_EMPTY)

	headers, err := GetHeadersFromHash(hash, hash)
	assert.Nil(t, err)

	buf := msgpack.NewHeaders(headers)
	sink := common.NewZeroCopySink(nil)
	types.WriteMessage(sink, buf)
	realHeaderMsg, _, err := types.ReadMessage(bytes.NewBuffer(sink.Bytes()))
	assert.Nil(t, err)
	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: realHeaderMsg,
	}

	// Invoke BlkHeaderHandle to handle the msg
	ctx := newContext(t, msg, network)
	BlkHeaderHandle(ctx, realHeaderMsg.(*types.BlkHeader))

	network.DelNbrNode(testID.ToUint64())
}

// TestBlockHandle tests Function BlockHandle handling a block message
func TestBlockHandle(t *testing.T) {
	// Simulate a remote peer to be added to the neighbor peers
	testID := kbucket.PseudoKadIdFromUint64(123456)

	remotePeer := peer.NewPeer()
	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:50010")

	network.AddNbrNode(remotePeer)

	// Construct a block packet
	hash := ledger.DefLedger.GetBlockHash(0)
	assert.NotEqual(t, hash, common.UINT256_EMPTY)

	block, err := ledger.DefLedger.GetBlockByHash(hash)
	assert.Nil(t, err)

	mr, err := common.Uint256FromHexString("1b8fa7f242d0eeb4395f89cbb59e4c29634047e33245c4914306e78a88e14ce5")
	assert.Nil(t, err)
	buf := msgpack.NewBlock(block, mr)

	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: buf,
	}

	// Invoke BlockHandle to handle the msg
	ctx := newContext(t, msg, network)
	BlockHandle(ctx, buf.(*types.Block))

	network.DelNbrNode(testID.ToUint64())
}

// TestDataReqHandle tests Function DataReqHandle handling a data req(block/Transaction)
func TestDataReqHandle(t *testing.T) {
	testID := kbucket.PseudoKadIdFromUint64(0x7533345)

	remotePeer := peer.NewPeer()
	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:50010")

	network.AddNbrNode(remotePeer)

	hash := ledger.DefLedger.GetBlockHash(0)
	assert.NotEqual(t, hash, common.UINT256_EMPTY)
	buf := msgpack.NewBlkDataReq(hash)

	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: buf,
	}

	ctx := newContext(t, msg, network)
	DataReqHandle(ctx, buf.(*types.DataReq))

	tempStr := "3369930accc1ddd067245e8edadcd9bea207ba5e1753ac18a51df77a343bfe92"
	hex, _ := hex.DecodeString(tempStr)
	var txHash common.Uint256
	txHash.Deserialize(bytes.NewReader(hex))
	buf = msgpack.NewTxnDataReq(txHash)
	msg = &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: buf,
	}

	ctx = newContext(t, msg, network)
	DataReqHandle(ctx, buf.(*types.DataReq))

	network.DelNbrNode(testID.ToUint64())
}

// TestInvHandle tests Function InvHandle handling an inventory message
func TestInvHandle(t *testing.T) {
	testID := kbucket.PseudoKadIdFromUint64(0x7533345)

	remotePeer := peer.NewPeer()
	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:50010")

	network.AddNbrNode(remotePeer)

	hash := ledger.DefLedger.GetBlockHash(0)
	assert.NotEqual(t, hash, common.UINT256_EMPTY)

	buf := bytes.NewBuffer([]byte{})
	hash.Serialize(buf)
	invPayload := msgpack.NewInvPayload(common.BLOCK, []common.Uint256{hash})
	buffer := msgpack.NewInv(invPayload)
	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: buffer,
	}

	ctx := newContext(t, msg, network)
	InvHandle(ctx, buffer.(*types.Inv))

	network.DelNbrNode(testID.ToUint64())
}

// TestDisconnectHandle tests Function DisconnectHandle handling a disconnect event
func TestDisconnectHandle(t *testing.T) {
	testID := kbucket.PseudoKadIdFromUint64(0x7533345)

	remotePeer := peer.NewPeer()
	remotePeer.UpdateInfo(time.Now(), 1, 12345678, 20336,
		testID, 0, 12345, "1.5.2")
	remotePeer.Link.SetAddr("127.0.0.1:50010")

	network.AddNbrNode(remotePeer)

	msg := &types.MsgPayload{
		Id:      testID.ToUint64(),
		Addr:    "127.0.0.1:50010",
		Payload: &types.Disconnected{},
	}

	ctx := newContext(t, msg, network)
	DisconnectHandle(ctx)

	network.DelNbrNode(testID.ToUint64())
}
