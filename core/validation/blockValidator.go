package validation

import (
	"DNA/core/ledger"
	tx "DNA/core/transaction"
	. "DNA/errors"
	"errors"
	"fmt"
)

func VerifyBlock(block *ledger.Block, ld *ledger.Ledger, completely bool) error {
	if block.Blockdata.Height == 0 {
		return nil
	}
	err := VerifyBlockData(block.Blockdata, ld)
	if err != nil {
		return err
	}

	flag, err := VerifySignableData(block)
	if flag == false || err != nil {
		return err
	}

	if block.Transactions == nil {
		return errors.New(fmt.Sprintf("No Transactions Exist in Block."))
	}
	if block.Transactions[0].TxType != tx.BookKeeping {
		return errors.New(fmt.Sprintf("Blockdata Verify failed first Transacion in block is not BookKeeping type."))
	}
	for index, v := range block.Transactions {
		if v.TxType == tx.BookKeeping && index != 0 {
			return errors.New(fmt.Sprintf("This Block Has BookKeeping transaction after first transaction in block."))
		}
	}

	//verfiy block's transactions
	if completely {
		/*
			//TODO: NextMiner Check.
			mineraddress, err := ledger.GetMinerAddress(ld.Blockchain.GetMinersByTXs(block.Transactions))
			if err != nil {
				return errors.New(fmt.Sprintf("GetMinerAddress Failed."))
			}
			if block.Blockdata.NextMiner != mineraddress {
				return errors.New(fmt.Sprintf("Miner is not validate."))
			}
		*/
		for _, txVerify := range block.Transactions {
			if err := VerifyTransaction(txVerify); err != nil {
				return errors.New(fmt.Sprintf("VerifyTransaction error when verifiy block"))
			}
			if err := VerifyTransactionWithLedger(txVerify, ledger.DefaultLedger); err != nil {
				return errors.New(fmt.Sprintf("VerifyTransactionWithLedger error when verifiy block"))
			}
		}
	}

	return nil
}

func VerifyHeader(bd *ledger.Header, ledger *ledger.Ledger) error {
	return VerifyBlockData(bd.Blockdata, ledger)
}

func VerifyBlockData(bd *ledger.Blockdata, ledger *ledger.Ledger) error {
	if bd.Height == 0 {
		return nil
	}

	prevHeader, err := ledger.Blockchain.GetHeader(bd.PrevBlockHash)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[BlockValidator], Cannnot find prevHeader..")
	}
	if prevHeader == nil {
		return NewDetailErr(err, ErrNoCode, "[BlockValidator], Cannnot find previous block.")
	}

	if prevHeader.Blockdata.Height+1 != bd.Height {
		return NewDetailErr(err, ErrNoCode, "[BlockValidator], block height is incorrect.")
	}

	if prevHeader.Blockdata.Timestamp >= bd.Timestamp {
		return NewDetailErr(err, ErrNoCode, "[BlockValidator], block timestamp is incorrect.")
	}

	return nil
}
