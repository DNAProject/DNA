package validation

import (
	"DNA/common"
	"DNA/core/ledger"
	tx "DNA/core/transaction"
	"DNA/core/transaction/payload"
	"errors"
	"math"
)

// VerifyTransaction verifys received single transaction
func VerifyTransaction(Tx *tx.Transaction) error {

	if err := CheckDuplicateInput(Tx); err != nil {
		return err
	}

	if err := CheckAssetPrecision(Tx); err != nil {
		return err
	}

	if err := CheckTransactionBalance(Tx); err != nil {
		return err
	}

	if err := CheckAttributeProgram(Tx); err != nil {
		return err
	}

	if err := CheckTransactionContracts(Tx); err != nil {
		return err
	}

	if err := CheckTransactionPayload(Tx); err != nil {
		return err
	}

	return nil
}

// VerifyTransactionWithBlock verifys a transaction with current transaction pool in memory
func VerifyTransactionWithBlock(TxPool []*tx.Transaction) error {
	//initial
	txnlist := make(map[common.Uint256]*tx.Transaction, 0)
	var txPoolInputs []string
	//sum all inputs in TxPool
	for _, Tx := range TxPool {
		for _, UTXOinput := range Tx.UTXOInputs {
			txPoolInputs = append(txPoolInputs, UTXOinput.ToString())
		}
	}
	//start check
	for _, txn := range TxPool {
		//1.check weather have duplicate transaction.
		if _, exist := txnlist[txn.Hash()]; exist {
			return errors.New("[VerifyTransactionWithBlock], duplicate transaction exist in block.")
		} else {
			txnlist[txn.Hash()] = txn
		}
		//2.check Duplicate Utxo input
		if err := CheckDuplicateUtxoInBlock(txn, txPoolInputs); err != nil {
			return err
		}
		//3.check issue amount
		switch txn.TxType {
		case tx.IssueAsset:
			//TODO: use delta mode to improve performance
			results := txn.GetMergedAssetIDValueFromOutputs()
			for k, _ := range results {
				//Get the Asset amount when RegisterAsseted.
				trx, err := tx.TxStore.GetTransaction(k)
				if trx.TxType != tx.RegisterAsset {
					return errors.New("[VerifyTransaction], TxType is illegal.")
				}
				AssetReg := trx.Payload.(*payload.RegisterAsset)

				//Get the amount has been issued of this assetID
				var quantity_issued common.Fixed64
				if AssetReg.Amount < common.Fixed64(0) {
					continue
				} else {
					quantity_issued, err = tx.TxStore.GetQuantityIssued(k)
					if err != nil {
						return errors.New("[VerifyTransaction], GetQuantityIssued failed.")
					}
				}

				//calc the amounts in txPool which are also IssueAsset
				var txPoolAmounts common.Fixed64
				for _, t := range TxPool {
					if t.TxType == tx.IssueAsset {
						outputResult := t.GetMergedAssetIDValueFromOutputs()
						for txidInPool, txValueInPool := range outputResult {
							if txidInPool == k {
								txPoolAmounts = txPoolAmounts + txValueInPool
							}
						}
					}
				}

				//calc weather out off the amount when Registed.
				//AssetReg.Amount : amount when RegisterAsset of this assedID
				//quantity_issued : amount has been issued of this assedID
				//txPoolAmounts   : amount in transactionPool of this assedID of issue transaction.
				if AssetReg.Amount-quantity_issued < txPoolAmounts {
					return errors.New("[VerifyTransaction], Amount check error.")
				}
			}
		}

	}

	return nil
}

// VerifyTransactionWithLedger verifys a transaction with history transaction in ledger
func VerifyTransactionWithLedger(Tx *tx.Transaction, ledger *ledger.Ledger) error {
	if IsDoubleSpend(Tx, ledger) {
		return errors.New("[IsDoubleSpend] faild.")
	}
	return nil
}

//validate the transaction of duplicate UTXO input
func CheckDuplicateInput(tx *tx.Transaction) error {
	if len(tx.UTXOInputs) == 0 {
		return nil
	}
	for i, utxoin := range tx.UTXOInputs {
		for j := 0; j < i; j++ {
			if utxoin.ReferTxID == tx.UTXOInputs[j].ReferTxID && utxoin.ReferTxOutputIndex == tx.UTXOInputs[j].ReferTxOutputIndex {
				return errors.New("invalid transaction")
			}
		}
	}
	return nil
}

func CheckDuplicateUtxoInBlock(tx *tx.Transaction, txPoolInputs []string) error {
	var txInputs []string
	for _, t := range tx.UTXOInputs {
		txInputs = append(txInputs, t.ToString())
	}
	for _, i := range txInputs {
		for _, j := range txPoolInputs {
			if i == j {
				return errors.New("Duplicated UTXO inputs found in tx pool")
			}
		}
	}
	return nil
}

func IsDoubleSpend(tx *tx.Transaction, ledger *ledger.Ledger) bool {
	return ledger.IsDoubleSpend(tx)
}

func CheckAssetPrecision(Tx *tx.Transaction) error {
	for k, outputs := range Tx.AssetOutputs {
		asset, err := ledger.DefaultLedger.GetAsset(k)
		if err != nil {
			return errors.New("The asset not exist in local blockchain.")
		}
		precision := asset.Precision
		for _, output := range outputs {
			if output.Value.GetData()%int64(math.Pow(10, 8-float64(precision))) != 0 {
				return errors.New("The precision of asset is incorrect.")
			}
		}
	}
	return nil
}

func CheckTransactionBalance(Tx *tx.Transaction) error {
	if len(Tx.AssetInputAmount) != len(Tx.AssetOutputAmount) {
		return errors.New("The number of asset is not same between inputs and outputs.")
	}

	for k, v := range Tx.AssetInputAmount {
		if v != Tx.AssetOutputAmount[k] {
			return errors.New("The amount of asset is not same between inputs and outputs.")
		}
	}
	return nil
}

func CheckAttributeProgram(Tx *tx.Transaction) error {
	//TODO: implement CheckAttributeProgram
	return nil
}

func CheckTransactionContracts(Tx *tx.Transaction) error {
	flag, err := VerifySignableData(Tx)
	if flag && err == nil {
		return nil
	} else {
		return err
	}
}

func CheckTransactionPayload(Tx *tx.Transaction) error {

	switch pld := Tx.Payload.(type) {
	case *payload.BookKeeper:
		//Todo: validate bookKeeper Cert
		_ = pld.Cert
		return nil
	case *payload.RegisterAsset:
	case *payload.IssueAsset:
	case *payload.TransferAsset:
	case *payload.BookKeeping:
	case *payload.PrivacyPayload:
	case *payload.Record:
	case *payload.DeployCode:
	case *payload.DataFile:
	default:
		return errors.New("[txValidator],invalidate transaction payload type.")
	}
	return nil
}
