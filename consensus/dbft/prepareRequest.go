package dbft

import (
	. "DNA/common"
	"DNA/common/log"
	ser "DNA/common/serialization"
	tx "DNA/core/transaction"
	. "DNA/errors"
	"fmt"
	"io"
)

type PrepareRequest struct {
	msgData ConsensusMessageData

	Nonce             uint64
	NextMiner         Uint160
	TransactionHashes []Uint256
	Transactions      []*tx.Transaction
	Signature         []byte
}

func (pr *PrepareRequest) Serialize(w io.Writer) error {
	log.Trace()
	pr.msgData.Serialize(w)
	err := ser.WriteVarUint(w, pr.Nonce)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "PrepareRequest Execute WriteVarUint failed.")
	}
	_, err = pr.NextMiner.Serialize(w)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "PrepareRequest Execute NextMiner.Serialize failed.")
	}
	//Serialize  Transaction's hashes
	txNum := uint64(len(pr.TransactionHashes))
	err = ser.WriteVarUint(w, txNum)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "PrepareRequest Execute WriteVarUint. failed.")
	}
	for _, txHash := range pr.TransactionHashes {
		_, err = txHash.Serialize(w)
		if err != nil {
			return NewDetailErr(err, ErrNoCode, "PrepareRequest Execute txHash.Serialize. failed.")
		}
	}

	for _, t := range pr.Transactions {
		err = t.Serialize(w)
		if err != nil {
			return NewDetailErr(err, ErrNoCode, "PrepareRequest Execute t.Serialize. failed.")
		}
	}

	err = ser.WriteVarBytes(w, pr.Signature)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "PrepareRequest Execute ser.WriteVarBytes failed.")
	}
	return nil
}

//read data to reader
func (pr *PrepareRequest) Deserialize(r io.Reader) error {
	log.Trace()
	pr.msgData = ConsensusMessageData{}
	pr.msgData.Deserialize(r)
	pr.Nonce, _ = ser.ReadVarUint(r, 0)
	pr.NextMiner = Uint160{}
	pr.NextMiner.Deserialize(r)

	//TransactionHashes
	len, err := ser.ReadVarUint(r, 0)
	if err != nil {
		return err
	}

	if len == 0 {
		fmt.Printf("The hash len at consensus payload is 0\n")
	} else {

		// deser hashes
		pr.TransactionHashes = make([]Uint256, len)
		for i := uint64(0); i < len; i++ {
			hash := new(Uint256)
			err = hash.Deserialize(r)
			if err != nil {
				return err
			}
			pr.TransactionHashes[i] = *hash
		}

		// deser txs
		pr.Transactions = make([]*tx.Transaction, len)
		for i := uint64(0); i < len; i++ {
			var t tx.Transaction
			err = t.Deserialize(r)
			if err != nil {
				return err
			}
			pr.Transactions[i] = &t
		}

		//if pr.BookkeepingTransaction.Hash() != pr.TransactionHashes[0] {
		//	log.Debug("pr.BookkeepingTransaction.Hash()=", pr.BookkeepingTransaction.Hash())
		//	log.Debug("pr.TransactionHashes[0]=", pr.TransactionHashes[0])
		//	return NewDetailErr(nil, ErrNoCode, "The Bookkeeping Transaction data is incorrect.")

		//}
	}

	pr.Signature, err = ser.ReadVarBytes(r)
	if err != nil {
		fmt.Printf("Parse the Signature error\n")
		return err
	}

	return nil
}

func (pr *PrepareRequest) Type() ConsensusMessageType {
	log.Trace()
	return pr.ConsensusMessageData().Type
}

func (pr *PrepareRequest) ViewNumber() byte {
	log.Trace()
	return pr.msgData.ViewNumber
}

func (pr *PrepareRequest) ConsensusMessageData() *ConsensusMessageData {
	log.Trace()
	return &(pr.msgData)
}
