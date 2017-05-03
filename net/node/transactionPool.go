package node

import (
	"DNA/common"
	"DNA/common/log"
	"DNA/core/ledger"
	"DNA/core/transaction"
	"DNA/core/transaction/payload"
	va "DNA/core/validation"
	"DNA/errors"
	msg "DNA/net/message"
	. "DNA/net/protocol"
	"fmt"
	"sync"
)

type TXNPool struct {
	sync.RWMutex
	txProcessList     // transaction which have been verifyed will put into this map
	tXNPendingList    // transaction which didn't pass the verify will put into this map
	assetIssueSummary // transaction which pass the verify will summary the amout to this map
	inputUTXOList     // transaction which pass the verify will add the UTXO to this map
}

func (txnPool *TXNPool) init() {
	txnPool.initAssetIssueSummary()
	txnPool.initTxProcessList()
	txnPool.initTXNPendingList()
	txnPool.initInputUTXOList()
}

func (txnPool *TXNPool) AppendTxnPool(txn *transaction.Transaction) bool {
	//verify transaction with Concurrency
	if err, errCode := va.VerifyTransactionCanConcurrency(txn); err != nil {
		if errCode != errors.ErrDuplicatedTx {
			txnPool.appendTXNPendingList(txn, errCode)
		}
		return false
	}

	//verify transaction by pool with lock
	txnPool.Lock()
	defer txnPool.Unlock()
	if ok := txnPool.verifyTransactionWithTxPool(txn); !ok {
		return false
	}
	txnPool.appendToProcessList(txn)
	return true
}

// Attention: clean the trasaction Pool after the consensus confirmed all of the transcation
func (txnPool *TXNPool) GetTxnPool(cleanPool bool) map[common.Uint256]*transaction.Transaction {
	txnPool.Lock()
	defer txnPool.Unlock()
	txList := txnPool.getProcessTxnList(cleanPool)
	return txList
}

//Attention: clean the trasaction Pool with committed transactions.
func (txnPool *TXNPool) CleanSubmittedTransactions(block *ledger.Block) error {
	txnPool.Lock()
	defer txnPool.Unlock()
	log.Trace()

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		txnPool.cleanProcessTxnPool(block.Transactions)
		wg.Done()
	}()
	go func() {
		txnPool.cleanTxnPoolUtxoMap(block.Transactions)
		wg.Done()
	}()
	go func() {
		txnPool.cleanSubmittedTransactionsOfTXNPendingList(block.Transactions)
		wg.Done()
	}()
	wg.Wait()
	txnPool.clearAndReverifyAssetIssueAmoutSummaryMap()
	return nil
}

func (txnPool *TXNPool) GetTransaction(hash common.Uint256) *transaction.Transaction {
	txnPool.RLock()
	defer txnPool.RUnlock()
	return txnPool.getTransaction(hash)
}

func (txnPool *TXNPool) ClearTransactionsPendingPool() {
	txnPool.Lock()
	defer txnPool.Unlock()
	txnPool.clearTransactionsPendingPool()
	return
}

func (txnPool *TXNPool) GetTransactionPendingReason(hash common.Uint256) errors.ErrCode {
	return txnPool.getTransactionPendingReason(hash)
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

func (txnPool *TXNPool) clearAndReverifyAssetIssueAmoutSummaryMap() {
	txnPool.clearAssetIssueAmoutSummaryMap()
	txnPool.reverifyIssueAmoutInPool(txnPool.getTxList())

}

func (txnPool *TXNPool) verifyTransactionWithTxPool(txn *transaction.Transaction) bool {
	//check weather have duplicate UTXO input
	ok, duplicateTxn := txnPool.checkWithUTXOinPool(txn)
	if !ok {
		//delete from associated map
		txnPool.deleteFromProcessTransactionList(duplicateTxn)
		txnPool.deleteFromSummary(duplicateTxn)
		txnPool.deleteFromUTXOList(duplicateTxn)
		//append erro transaction to pending list
		txnPool.appendTXNPendingList(duplicateTxn, errors.ErrCheckDuplicatedUTXOInput)
		txnPool.appendTXNPendingList(txn, errors.ErrCheckDuplicatedUTXOInput)
		return false
	}
	//check issue transaction weather occur exceed issue range.
	if ok := txnPool.verifyIssueAmoutInPool(txn); !ok {
		txnPool.appendTXNPendingList(txn, errors.ErrCheckExceedTheRegAmount)
		return false
	}
	return true
}

type inputUTXOList struct {
	inputUTXOList map[*transaction.TxOutput]*transaction.Transaction
}

func (list *inputUTXOList) initInputUTXOList() {
	list.inputUTXOList = make(map[*transaction.TxOutput]*transaction.Transaction)
}

func (list *inputUTXOList) checkWithUTXOinPool(txn *transaction.Transaction) (bool, *transaction.Transaction) {
	reference, err := txn.GetReference()
	if err != nil {
		return false, nil
	}
	for _, v := range reference {
		if transaction, ok := list.inputUTXOList[v]; ok {
			return false, transaction
		} else {
			list.inputUTXOList[v] = txn
		}
	}
	return true, nil
}

func (list *inputUTXOList) deleteFromUTXOList(txn *transaction.Transaction) {
	result, err := txn.GetReference()
	if err != nil {
		log.Warn(fmt.Sprintf("Transaction =%x not Exist in Pool when delete.", txn.Hash()))
	}
	for _, v := range result {
		if _, ok := list.inputUTXOList[v]; ok {
			delete(list.inputUTXOList, v)
		}
	}
}

func (list *inputUTXOList) cleanTxnPoolUtxoMap(txs []*transaction.Transaction) {
	for _, tx := range txs {
		inputUtxos, _ := tx.GetReference()
		for _, v := range inputUtxos {
			delete(list.inputUTXOList, v)
		}
	}
}

type tXNPendingList struct {
	txnCnt uint64
	txList map[common.Uint256]errors.ErrCode // map[transaction hash]errorCode
}

func (tpl *tXNPendingList) initTXNPendingList() {
	tpl.txList = make(map[common.Uint256]errors.ErrCode)
	tpl.txnCnt = 0
}

func (tpl *tXNPendingList) getTransactionPendingReason(hash common.Uint256) errors.ErrCode {
	reason := tpl.txList[hash]
	return reason
}

func (tpl *tXNPendingList) appendTXNPendingList(txn *transaction.Transaction, errcode errors.ErrCode) bool {
	hash := txn.Hash()
	if _, ok := tpl.txList[hash]; ok {
		return true
	}
	tpl.txList[hash] = errcode
	tpl.txnCnt++
	return true
}

// Attention: clean the trasaction Pool with committed transactions.
func (tpl *tXNPendingList) cleanSubmittedTransactionsOfTXNPendingList(txns []*transaction.Transaction) error {
	cleaned := 0
	for _, tx := range txns {
		delete(tpl.txList, tx.Hash())
		cleaned++
	}
	if len(tpl.txList) > 0 {
		log.Debug(fmt.Sprintf("%d Transactions pending at pendingPool.\n", len(tpl.txList)))
	}
	return nil
}

func (tpl *tXNPendingList) clearTransactionsPendingPool() {
	for k, _ := range tpl.txList {
		delete(tpl.txList, k)
	}
	log.Info("Transactions pendingPool clear completed.\n")
}

type assetIssueSummary struct {
	assetIssueAmoutSummary map[common.Uint256]common.Fixed64
}

func (ais *assetIssueSummary) initAssetIssueSummary() {
	ais.assetIssueAmoutSummary = make(map[common.Uint256]common.Fixed64)
}

// verifyTransactionWithTxPool verifys a transaction with current transaction pool in memory which CAN NOT Concurrency.
func (asm *assetIssueSummary) verifyIssueAmoutInPool(txn *transaction.Transaction) bool {
	if txn.TxType != transaction.IssueAsset {
		return true
	}
	transactionResult := txn.GetMergedAssetIDValueFromOutputs()
	for k, delta := range transactionResult {
		//update the amount in txPool
		if amout, ok := asm.assetIssueAmoutSummary[k]; ok {
			asm.assetIssueAmoutSummary[k] = amout + delta
		} else {
			asm.assetIssueAmoutSummary[k] = delta
		}

		//Check weather occur exceed the amount when RegisterAsseted
		//1. Get the Asset amount when RegisterAsseted.
		trx, err := transaction.TxStore.GetTransaction(k)
		if err != nil {
			return false
		}
		if trx.TxType != transaction.RegisterAsset {
			return false
		}
		AssetReg := trx.Payload.(*payload.RegisterAsset)

		//2. Get the amount has been issued of this assetID
		var quantity_issued common.Fixed64
		if AssetReg.Amount < common.Fixed64(0) {
			continue
		} else {
			quantity_issued, err = transaction.TxStore.GetQuantityIssued(k)
			if err != nil {
				return false
			}
		}

		//3. calc weather out off the amount when Registed.
		//AssetReg.Amount : amount when RegisterAsset of this assedID
		//quantity_issued : amount has been issued of this assedID
		//txnPool.assetIssueAmoutSummary[k] : amount in transactionPool of this assedID
		if AssetReg.Amount-quantity_issued < asm.assetIssueAmoutSummary[k] {
			return false
		}
	}
	return true
}

func (asm *assetIssueSummary) deleteFromSummary(txn *transaction.Transaction) {
	if txn.TxType != transaction.IssueAsset {
		return
	}
	transactionResult := txn.GetMergedAssetIDValueFromOutputs()
	for k, delta := range transactionResult {
		if amout, ok := asm.assetIssueAmoutSummary[k]; ok {
			asm.assetIssueAmoutSummary[k] = amout - delta
		}
	}
}

func (asm *assetIssueSummary) clearAssetIssueAmoutSummaryMap() {
	for assit, _ := range asm.assetIssueAmoutSummary {
		delete(asm.assetIssueAmoutSummary, assit)
	}
}

func (asm *assetIssueSummary) reverifyIssueAmoutInPool(txns map[common.Uint256]*transaction.Transaction) {
	for _, v := range txns {
		asm.verifyIssueAmoutInPool(v)
	}
}

type txProcessList struct {
	txnCnt uint64 // count
	txList map[common.Uint256]*transaction.Transaction
}

func (txp *txProcessList) initTxProcessList() {
	txp.txList = make(map[common.Uint256]*transaction.Transaction)
	txp.txnCnt = 0
}

func (txp *txProcessList) getTxList() map[common.Uint256]*transaction.Transaction {
	return txp.txList
}

// Attention: clean the trasaction Pool with committed transactions.
func (txp *txProcessList) deleteFromProcessTransactionList(txs *transaction.Transaction) error {
	delete(txp.txList, txs.Hash())
	txp.txnCnt--
	return nil
}

// Attention: clean the trasaction Pool with committed transactions.
func (txp *txProcessList) appendToProcessList(txs *transaction.Transaction) {
	txp.txList[txs.Hash()] = txs
	txp.txnCnt++
}

// Attention: clean the trasaction Pool with committed transactions.
func (txp *txProcessList) duplicateCheck(txs *transaction.Transaction) bool {
	if _, ok := txp.txList[txs.Hash()]; ok {
		return true
	}
	return false
}

// Attention: clean the trasaction Pool after the consensus confirmed all of the transcation
func (txp *txProcessList) getProcessTxnList(cleanPool bool) map[common.Uint256]*transaction.Transaction {
	txList := txp.txList
	if cleanPool == true {
		txp.initTxProcessList()
	}
	return deepCopy(txList)
}

func deepCopy(mapIn map[common.Uint256]*transaction.Transaction) map[common.Uint256]*transaction.Transaction {
	reply := make(map[common.Uint256]*transaction.Transaction)
	for k, v := range mapIn {
		reply[k] = v
	}
	return reply
}

// Attention: clean the trasaction Pool with committed transactions.
func (txp *txProcessList) cleanProcessTxnPool(txs []*transaction.Transaction) error {
	txsNum := len(txs) - 1
	txInPoolNum := len(txp.txList)
	cleaned := 0
	for _, tx := range txs {
		if tx.TxType != transaction.BookKeeping {
			delete(txp.txList, tx.Hash())
			cleaned++
		}
	}
	if txsNum-cleaned != 0 {
		log.Info(fmt.Sprintf("The Transactions num Unmatched. Expect %d, got %d .\n", txsNum, cleaned))
	}
	log.Debug(fmt.Sprintf("[cleanProcessTxnPool], Requested %d clean, %d transactions cleaned from localNode.TransPool and remains %d still in TxPool", txsNum, cleaned, txInPoolNum-cleaned))
	return nil
}

func (txp *txProcessList) getTransaction(hash common.Uint256) *transaction.Transaction {
	txn := txp.txList[hash]
	return txn
}
