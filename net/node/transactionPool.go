package node

import (
	"DNA/common"
	"DNA/core/transaction"
	msg "DNA/net/message"
	. "DNA/net/protocol"
	"sync"
)

type TXNPool struct {
	sync.RWMutex
	txnCnt	uint64
	list map[common.Uint256]*transaction.Transaction
}

func (txnPool *TXNPool) GetTransaction(hash common.Uint256) *transaction.Transaction {
	txnPool.RLock()
	defer txnPool.RUnlock()
	txn := txnPool.list[hash]
	// Fixme need lock
	return txn
}

func (txnPool *TXNPool) AppendTxnPool(txn *transaction.Transaction) bool {
	txnPool.Lock()
	defer txnPool.Unlock()

	hash := txn.Hash()
	if _, ret := txnPool.list[hash]; ret {
		return false
	} else {
		txnPool.list[hash] = txn
		txnPool.txnCnt++
	}

	return true
}

// Attention: clean the trasaction Pool after the consensus confirmed all of the transcation
func (txnPool *TXNPool) GetTxnPool(cleanPool bool) map[common.Uint256]*transaction.Transaction {
	txnPool.Lock()
	defer txnPool.Unlock()

	list := txnPool.list
	if cleanPool == true {
		txnPool.init()
	}
	return list
}

func (txnPool *TXNPool) init() {
	txnPool.list = make(map[common.Uint256]*transaction.Transaction)
	txnPool.txnCnt = 0
}

func (node *node) SynchronizeTxnPool() {
	node.nbrNodes.RLock()
	defer node.nbrNodes.RUnlock()

	for _, n := range node.nbrNodes.List {
		if n.state == ESTABLISH {
			msg.ReqTxnPool(n)
		}
	}
}
