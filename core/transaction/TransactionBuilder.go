package transaction

import (
	"DNA/common"
	"DNA/core/asset"
	"DNA/core/contract/program"
	"DNA/core/transaction/payload"
	"DNA/crypto"
)

//initial a new transaction with asset registration payload
func NewRegisterAssetTransaction(asset *asset.Asset, amount common.Fixed64, issuer *crypto.PubKey, conroller common.Uint160, height uint32) (*Transaction, error) {

	//TODO: check arguments

	assetRegPayload := &payload.RegisterAsset{
		Asset:  asset,
		Amount: amount,
		//Precision: precision,
		Issuer:     issuer,
		Controller: conroller,
	}

	tx := &Transaction{
		//nonce uint64 //TODO: genenrate nonce
		UTXOInputs:      []*UTXOTxInput{},
		BalanceInputs:   []*BalanceTxInput{},
		Attributes:      []*TxAttribute{},
		TxType:          RegisterAsset,
		Payload:         assetRegPayload,
		Programs:        []*program.Program{},
		CurrBlockHeight: height,
	}
	tx.SetTransactionVersion(1)
	return tx, nil
}

//initial a new transaction with asset registration payload
func NewBookKeeperTransaction(pubKey *crypto.PubKey, isAdd bool, cert []byte, issuer *crypto.PubKey, height uint32) (*Transaction, error) {

	bookKeeperPayload := &payload.BookKeeper{
		PubKey: pubKey,
		Action: payload.BookKeeperAction_SUB,
		Cert:   cert,
		Issuer: issuer,
	}

	if isAdd {
		bookKeeperPayload.Action = payload.BookKeeperAction_ADD
	}

	tx := &Transaction{
		TxType:          BookKeeper,
		Payload:         bookKeeperPayload,
		UTXOInputs:      []*UTXOTxInput{},
		BalanceInputs:   []*BalanceTxInput{},
		Attributes:      []*TxAttribute{},
		Programs:        []*program.Program{},
		CurrBlockHeight: height,
	}
	tx.SetTransactionVersion(1)
	return tx, nil
}

func NewIssueAssetTransaction(outputs []*TxOutput, height uint32) (*Transaction, error) {
	assetRegPayload := &payload.IssueAsset{}

	tx := &Transaction{
		TxType:          IssueAsset,
		Payload:         assetRegPayload,
		Attributes:      []*TxAttribute{},
		BalanceInputs:   []*BalanceTxInput{},
		Outputs:         outputs,
		Programs:        []*program.Program{},
		CurrBlockHeight: height,
	}
	tx.SetTransactionVersion(1)
	return tx, nil
}

func NewTransferAssetTransaction(inputs []*UTXOTxInput, outputs []*TxOutput, height uint32) (*Transaction, error) {
	assetRegPayload := &payload.TransferAsset{}

	tx := &Transaction{
		TxType:          TransferAsset,
		Payload:         assetRegPayload,
		Attributes:      []*TxAttribute{},
		UTXOInputs:      inputs,
		BalanceInputs:   []*BalanceTxInput{},
		Outputs:         outputs,
		Programs:        []*program.Program{},
		CurrBlockHeight: height,
	}
	tx.SetTransactionVersion(1)
	return tx, nil
}

func NewRecordTransaction(recordType string, recordData []byte, height uint32) (*Transaction, error) {
	recordPayload := &payload.Record{
		RecordType: recordType,
		RecordData: recordData,
	}

	tx := &Transaction{
		TxType:          Record,
		Payload:         recordPayload,
		Attributes:      []*TxAttribute{},
		UTXOInputs:      []*UTXOTxInput{},
		BalanceInputs:   []*BalanceTxInput{},
		Programs:        []*program.Program{},
		CurrBlockHeight: height,
	}
	tx.SetTransactionVersion(1)
	return tx, nil
}

func NewPrivacyPayloadTransaction(fromPrivKey []byte, fromPubkey *crypto.PubKey, toPubkey *crypto.PubKey, payloadType payload.EncryptedPayloadType, data []byte, height uint32) (*Transaction, error) {
	privacyPayload := &payload.PrivacyPayload{
		PayloadType: payloadType,
		EncryptType: payload.ECDH_AES256,
		EncryptAttr: &payload.EcdhAes256{
			FromPubkey: fromPubkey,
			ToPubkey:   toPubkey,
		},
	}
	privacyPayload.Payload, _ = privacyPayload.EncryptAttr.Encrypt(data, fromPrivKey)

	tx := &Transaction{
		TxType:          PrivacyPayload,
		Payload:         privacyPayload,
		Attributes:      []*TxAttribute{},
		UTXOInputs:      []*UTXOTxInput{},
		BalanceInputs:   []*BalanceTxInput{},
		Programs:        []*program.Program{},
		CurrBlockHeight: height,
	}
	tx.SetTransactionVersion(1)
	return tx, nil
}
func NewDataFileTransaction(path string, fileName string, note string, issuer *crypto.PubKey, height uint32) (*Transaction, error) {
	//TODO: check arguments
	DataFilePayload := &payload.DataFile{
		IPFSPath: path,
		Filename: fileName,
		Note:     note,
		Issuer:   issuer,
	}

	tx := &Transaction{
		TxType:          DataFile,
		Payload:         DataFilePayload,
		Attributes:      []*TxAttribute{},
		UTXOInputs:      []*UTXOTxInput{},
		BalanceInputs:   []*BalanceTxInput{},
		Programs:        []*program.Program{},
		CurrBlockHeight: height,
	}
	tx.SetTransactionVersion(1)
	return tx, nil
}
