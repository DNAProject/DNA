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

package peer

import (
	"fmt"
	"testing"
	"time"

	"github.com/DNAProject/DNA/p2pserver/common"
)

var id45 common.PeerId
var id46 common.PeerId
var id47 common.PeerId

func init() {
	id45 = common.PseudoPeerIdFromUint64(uint64(0x7533345))
	id46 = common.PseudoPeerIdFromUint64(uint64(0x7533346))
	id47 = common.PseudoPeerIdFromUint64(uint64(0x7533347))
}

func createPeers(cnt uint16) []*Peer {
	np := []*Peer{}
	var syncport uint16
	var height uint64
	for i := uint16(0); i < cnt; i++ {
		syncport = 20224 + i
		id := common.PseudoPeerIdFromUint64(0x7533345 + uint64(i))
		height = 434923 + uint64(i)
		p := NewPeer()
		p.UpdateInfo(time.Now(), 2, 3, syncport, id, 0, height, "1.5.2")
		p.SetState(3)
		p.SetHttpInfoState(true)
		p.Link.SetAddr("127.0.0.1:10338")
		np = append(np, p)
	}
	return np

}

func initTestNbrPeers() *NbrPeers {
	nm := NewNbrPeers()
	np := createPeers(5)
	for _, v := range np {
		nm.List[v.GetID()] = v
	}
	return nm
}

func TestNodeExisted(t *testing.T) {
	nm := initTestNbrPeers()

	if !nm.NodeExisted(id45) {
		t.Fatal("0x7533345 should in nbr peers")
	}
}

func TestGetPeer(t *testing.T) {
	nm := initTestNbrPeers()
	p := nm.GetPeer(id45)
	if p == nil {
		t.Fatal("TestGetPeer error")
	}
}

func TestDelNbrNode(t *testing.T) {
	nm := initTestNbrPeers()

	cnt := len(nm.List)
	p, delOK := nm.DelNbrNode(id45)
	if p == nil || !delOK {
		t.Fatal("TestDelNbrNode err")
	}
	if len(nm.List) != cnt-1 {
		t.Fatal("TestDelNbrNode not work")
	}
	p.DumpInfo()
}

func TestGetNeighborAddrs(t *testing.T) {
	nm := initTestNbrPeers()
	p := nm.GetPeer(id46)
	if p == nil {
		t.Fatal("TestGetNeighborAddrs:get peer error")
	}
	p.SetState(common.ESTABLISH)
	p = nm.GetPeer(id47)
	if p == nil {
		t.Fatal("TestGetNeighborAddrs:get peer error")
	}
	p.SetState(common.ESTABLISH)

	pList := nm.GetNeighborAddrs()
	for i := 0; i < len(pList); i++ {
		fmt.Printf("peer id = %s \n", pList[i].ID.ToHexString())
	}
	if len(pList) != 2 {
		t.Fatal("TestGetNeighborAddrs error")
	}
}

func TestGetNeighborHeights(t *testing.T) {
	nm := initTestNbrPeers()

	p := nm.GetPeer(id46)
	if p == nil {
		t.Fatal("TestGetNeighborHeights:get peer error")
	}
	p.SetState(common.ESTABLISH)
	p.SetHeight(110)

	p = nm.GetPeer(id47)
	if p == nil {
		t.Fatal("TestGetNeighborHeights:get peer error")
	}
	p.SetState(common.ESTABLISH)
	p.SetHeight(110)

	pMap := nm.GetNeighborHeights()
	if len(pMap) != 2 {
		t.Fatalf("expect pmap size to 2, got %d", len(pMap))
	}

	for k, v := range pMap {
		fmt.Printf("peer id = %s height = %d \n", k.ToHexString(), v)
		if v != 110 {
			t.Fatalf("expect height is 110, got %d", v)
		}
	}
}

func TestGetNeighbors(t *testing.T) {
	nm := initTestNbrPeers()
	p := nm.GetPeer(id46)
	if p == nil {
		t.Fatal("TestGetNeighbors:get peer error")
	}
	p.SetState(common.ESTABLISH)

	p = nm.GetPeer(id47)
	if p == nil {
		t.Fatal("TestGetNeighbors:get peer error")
	}
	p.SetState(common.ESTABLISH)

	pList := nm.GetNeighbors()
	if len(pList) != 2 {
		t.Fatalf("expect neigbor size is 2, got %d", len(pList))
	}

	for _, v := range pList {
		v.DumpInfo()
	}
}

func TestGetNbrNodeCnt(t *testing.T) {
	nm := initTestNbrPeers()

	p := nm.GetPeer(id46)
	if p == nil {
		t.Fatal("TestGetNbrNodeCnt:get peer error")
	}
	p.SetState(common.ESTABLISH)

	p = nm.GetPeer(id47)
	if p == nil {
		t.Fatal("TestGetNbrNodeCnt:get peer error")
	}
	p.SetState(common.ESTABLISH)

	if nm.GetNbrNodeCnt() != 2 {
		t.Fatalf("expect 2 neigbors got: %d", nm.GetNbrNodeCnt())
	}
}
