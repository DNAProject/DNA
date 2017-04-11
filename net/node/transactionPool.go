package node

import (
	"DNA/common"
	"DNA/common/log"
	"DNA/core/transaction"
	msg "DNA/net/message"
	. "DNA/net/protocol"
	"sync"
	"fmt"
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
	return DeepCopy(list)
}

func DeepCopy(mapIn map[common.Uint256]*transaction.Transaction) map[common.Uint256]*transaction.Transaction {
	reply := make( map[common.Uint256]*transaction.Transaction)
	for k, v := range mapIn {
		reply[k] =v
	}
	return reply
}

// Attention: clean the trasaction Pool with committed transactions.
func (txnPool *TXNPool) CleanTxnPool(txHashes []*common.Uint256) error{
	txnPool.Lock()
	defer txnPool.Unlock()

	txsNum := len(txHashes)
	txInPoolNum := len(txnPool.list)
	cleaned :=0
	for _, txHash := range txHashes {
		if _,ok:= txnPool.list[*txHash]; ok{
			delete(txnPool.list,*txHash)
			cleaned ++
		}else{
			log.Fatal("Delete failed of transaction hash =",txHash)
		}
	}

	log.Fatal(fmt.Sprintf("[CleanTxnPool], Requested %d clean, %d transactions cleaned from localNode.TransPool and remains %d still in TxPool",txsNum,cleaned,txInPoolNum-cleaned))
	return nil
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
