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

// Package p2p provides an network interface
package p2p

import (
	"github.com/DNAProject/DNA/p2pserver/common"
	"github.com/DNAProject/DNA/p2pserver/dht/kbucket"
	"github.com/DNAProject/DNA/p2pserver/message/types"
	"github.com/DNAProject/DNA/p2pserver/peer"
	"github.com/ontio/ontology-eventbus/actor"
)

//P2P represent the net interface of p2p package
type P2P interface {
	Start()
	Halt()
	Connect(addr string) error
	GetID() uint64
	GetKId() kbucket.KadId
	GetVersion() uint32
	GetPort() uint16
	GetRelay() bool
	GetHeight() uint64
	GetServices() uint64
	GetNeighbors() []*peer.Peer
	GetNeighborAddrs() []common.PeerAddr
	GetConnectionCnt() uint32
	GetMaxPeerBlockHeight() uint64
	GetNp() *peer.NbrPeers
	GetPeer(id uint64) *peer.Peer
	SetHeight(uint64)
	IsPeerEstablished(p *peer.Peer) bool
	Send(p *peer.Peer, msg types.Message) error
	GetPeerFromAddr(addr string) *peer.Peer
	GetOutConnRecordLen() uint
	AddPeerAddress(addr string, p *peer.Peer)
	RemovePeerAddress(addr string)
	AddNbrNode(*peer.Peer)
	DelNbrNode(id uint64) (*peer.Peer, bool)
	NodeEstablished(id uint64) bool
	Xmit(msg types.Message)
	IsOwnAddress(addr string) bool

	UpdateDHT(id kbucket.KadId) bool
	RemoveDHT(id kbucket.KadId) bool
	BetterPeers(id kbucket.KadId, count int) []kbucket.KadId
	GetKadKeyId() *kbucket.KadKeyId

	GetPeerStringAddr() map[uint64]string
	SetPID(pid *actor.PID)
}
