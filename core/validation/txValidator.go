package validation

import (
	"GoOnchain/core/ledger"
	tx "GoOnchain/core/transaction"
	"errors"
	"math"
)

//Verfiy the transcation for following points
//- Well format
//- No duplicated inputs
//- inputs/outputs balance
//- Transcation contracts pass
func VerifyTransaction(Tx *tx.Transaction, ledger *ledger.Ledger, TxPool []*tx.Transaction) error {

	err := CheckDuplicateInput(Tx)
	if(err != nil){return err}

	err = IsDoubleSpend(Tx,ledger)
	if(err != nil){return err}

	if TxPool != nil{
		err = CheckMemPool(Tx,TxPool)
		if(err != nil){return err}
	}


	err = CheckAssetPrecision(Tx)
	if(err != nil){return err}

	err = CheckTransactionBalance(Tx)
	if(err != nil){return err}

	err = CheckAttributeProgram(Tx)
	if(err != nil){return err}

	err = CheckTransactionContracts(Tx)
	if(err != nil){return err}

	return nil
}

func CheckMemPool(tx *tx.Transaction, TxPool []*tx.Transaction) error {

	for _, poolTx := range TxPool {
		for _, poolInput := range poolTx.UTXOInputs {
			for _, txInput := range tx.UTXOInputs {
				if poolInput.Equals(txInput) {
					return errors.New("There is duplicated Tx Input with Tx Pool.")
				}
			}
		}
	}
	return nil
}

func CheckDuplicateInput(tx *tx.Transaction) error {
	for i, utxoin := range tx.UTXOInputs {
		for j := 0; j < i; j++ {
			if utxoin.ReferTxID == tx.UTXOInputs[j].ReferTxID && utxoin.ReferTxOutputIndex == tx.UTXOInputs[j].ReferTxOutputIndex {
				return errors.New("invalid transaction")
			}
		}
	}
	return nil
}

func IsDoubleSpend(tx *tx.Transaction, ledger *ledger.Ledger) error {
	return ledger.IsDoubleSpend(tx)
}

func CheckAssetPrecision(Tx *tx.Transaction) error {
	for k, outputs := range Tx.AssetOutputs {
		asset,err:= ledger.DefaultLedger.GetAsset(k)
		if err!= nil{
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
	if (len(Tx.AssetInputAmount) != len(Tx.AssetOutputAmount)){
		return  errors.New("The number of asset is not same between inputs and outputs.")
	}


	for k, v := range Tx.AssetInputAmount{
		if(v != Tx.AssetOutputAmount[k]){
			return  errors.New("The amount of asset is not same between inputs and outputs.")
		}
	}
	return nil
}

func CheckAttributeProgram(Tx *tx.Transaction) error {
	//TODO: implement CheckAttributeProgram
	return nil
}

func CheckTransactionContracts(Tx *tx.Transaction) error {
	return VerifySignableData(Tx)
}
