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

package dht

import (
	"time"

	"github.com/DNAProject/DNA/common/log"
	kb "github.com/DNAProject/DNA/p2pserver/dht/kbucket"
)

// Pool size is the number of nodes used for group find/set RPC calls
var PoolSize = 6

// K is the maximum number of requests to perform before returning failure.
var KValue = 20

// Alpha is the concurrency factor for asynchronous requests.
var AlphaValue = 3

type DHT struct {
	localKeyId *kb.KadKeyId
	birth      time.Time // When this peer started up

	bucketSize int
	routeTable *kb.RouteTable // Array of routing tables for differently distanced nodes

	AutoRefresh           bool
	RtRefreshQueryTimeout time.Duration
	RtRefreshPeriod       time.Duration
}

// RouteTable return dht's routeTable
func (dht *DHT) RouteTable() *kb.RouteTable {
	return dht.routeTable
}

// NewDHT creates a new DHT with the specified host and options.
func NewDHT() *DHT {
	bucketSize := KValue
	keyId := kb.RandKadKeyId()
	rt := kb.NewRoutingTable(bucketSize, keyId.Id)

	rt.PeerAdded = func(p kb.KadId) {
		log.Debugf("dht: peer: %d added to dht", p)
	}

	rt.PeerRemoved = func(p kb.KadId) {
		log.Debugf("dht: peer: %d removed from dht", p)
	}

	return &DHT{
		localKeyId:            keyId,
		birth:                 time.Now(),
		routeTable:            rt,
		bucketSize:            bucketSize,
		AutoRefresh:           true,
		RtRefreshPeriod:       10 * time.Second,
		RtRefreshQueryTimeout: 10 * time.Second,
	}
}

// Update signals the routeTable to Update its last-seen status
// on the given peer.
func (dht *DHT) Update(peer kb.KadId) bool {
	err := dht.routeTable.Update(peer)
	return err == nil
}

func (dht *DHT) Remove(peer kb.KadId) {
	dht.routeTable.Remove(peer)
}

func (dht *DHT) GetKadKeyId() *kb.KadKeyId {
	return dht.localKeyId
}

func (dht *DHT) BetterPeers(id kb.KadId, count int) []kb.KadId {
	closer := dht.routeTable.NearestPeers(id, count)
	filtered := make([]kb.KadId, 0, len(closer))
	// don't include self and target id
	for _, curID := range closer {
		if curID == dht.localKeyId.Id {
			continue
		}
		if curID == id {
			continue
		}
		filtered = append(filtered, curID)
	}

	return filtered
}
