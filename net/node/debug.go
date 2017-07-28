package node

import (
	. "DNA/net/protocol"
	"expvar"
	"fmt"
	"time"
)

type PeerStatus struct {
	State         uint32
	FlightHeights []uint32
	LastContact   string
	TryTimes      uint32
	Addr          string // The address of the node
	Height        uint64 // The node latest block height
}

type TxPoolStatus struct {
	TxCount int
}

type nodeStatus struct {
	Id          uint64
	TxnCnt      uint64 // The transactions be transmit by this node
	RxTxnCnt    uint64 // The transaction received by this node
	PublicKey   string
	Peers       []PeerStatus
	Connectings []string
	TxPool      TxPoolStatus
}

func expvarNodeInfo(node *node) func() interface{} {

	return func() interface{} {
		pbkey, _ := node.publicKey.EncodePoint(true)

		ns := nodeStatus{
			Id:          node.id,
			TxnCnt:      node.txnCnt,
			RxTxnCnt:    node.rxTxnCnt,
			PublicKey:   fmt.Sprintf("%x", pbkey),
			TxPool:      TxPoolStatus{TxCount: node.TXNPool.Len()},
			Connectings: node.ConnectingAddrs,
		}

		node.nbrNodes.RLock()
		for _, n := range node.nbrNodes.List {
			peer := PeerStatus{
				State:         n.state,
				Height:        n.height,
				FlightHeights: n.flightHeights,
				LastContact:   fmt.Sprintf("%gs", float64(time.Now().Sub(n.link.time))/float64(time.Second)),
				TryTimes:      n.tryTimes,
				Addr:          fmt.Sprintf("%s:%d", n.link.addr, n.link.port),
			}

			ns.Peers = append(ns.Peers, peer)
		}
		node.nbrNodes.RUnlock()

		return ns
	}
}

func ExportNodeStatus(nd Noder) {
	expvar.Publish("dna_node", expvar.Func(expvarNodeInfo(nd.(*node))))
}
