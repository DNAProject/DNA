package main

import (
	"time"
	"runtime"
	"GoOnchain/node"
	"GoOnchain/node/httpjsonrpc"
)

const (
	// The number of the CPU cores for parallel optimization,TODO set from config file
	NCPU	 = 4
)

func init() {
	runtime.GOMAXPROCS(NCPU)
	go httpjsonrpc.StartServer()
}


func main() {
	time.Sleep(2 * time.Second)

	node.InitNodes()
	node.StartProtocol()
	httpjsonrpc.StartClient()

	// Modules start sample
	//ledger.Start(net.NetToLedgerCh <-chan *Msg, net.LedgerToNetCh chan<- *Msg)
	//consensus.Start(net.NetToConsensusCh <-chan *Msg, net.ConsensusToNetCh chan<- *Msg)

	for {

	}
}
