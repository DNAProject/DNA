// Copyright 2016 DNA Dev team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transaction

import (
	"DNA/common"
	"DNA/core/asset"
	"DNA/core/code"
	"DNA/core/contract/program"
	"DNA/core/transaction/payload"
	"DNA/crypto"
	"DNA/smartcontract/types"
)

//initial a new transaction with asset registration payload
func NewRegisterAssetTransaction(asset *asset.Asset, amount common.Fixed64, issuer *crypto.PubKey, conroller common.Uint160) (*Transaction, error) {

	//TODO: check arguments

	assetRegPayload := &payload.RegisterAsset{
		Asset:  asset,
		Amount: amount,
		//Precision: precision,
		Issuer:     issuer,
		Controller: conroller,
	}

	return &Transaction{
		//nonce uint64 //TODO: genenrate nonce
		UTXOInputs:    []*UTXOTxInput{},
		BalanceInputs: []*BalanceTxInput{},
		Attributes:    []*TxAttribute{},
		TxType:        RegisterAsset,
		Payload:       assetRegPayload,
		Programs:      []*program.Program{},
	}, nil
}

//initial a new transaction with asset registration payload
func NewBookKeeperTransaction(pubKey *crypto.PubKey, isAdd bool, cert []byte, issuer *crypto.PubKey) (*Transaction, error) {

	bookKeeperPayload := &payload.BookKeeper{
		PubKey: pubKey,
		Action: payload.BookKeeperAction_SUB,
		Cert:   cert,
		Issuer: issuer,
	}

	if isAdd {
		bookKeeperPayload.Action = payload.BookKeeperAction_ADD
	}

	return &Transaction{
		TxType:        BookKeeper,
		Payload:       bookKeeperPayload,
		UTXOInputs:    []*UTXOTxInput{},
		BalanceInputs: []*BalanceTxInput{},
		Attributes:    []*TxAttribute{},
		Programs:      []*program.Program{},
	}, nil
}

func NewIssueAssetTransaction(outputs []*TxOutput) (*Transaction, error) {

	assetRegPayload := &payload.IssueAsset{}

	return &Transaction{
		TxType:        IssueAsset,
		Payload:       assetRegPayload,
		Attributes:    []*TxAttribute{},
		BalanceInputs: []*BalanceTxInput{},
		Outputs:       outputs,
		Programs:      []*program.Program{},
	}, nil
}

func NewTransferAssetTransaction(inputs []*UTXOTxInput, outputs []*TxOutput) (*Transaction, error) {

	//TODO: check arguments

	assetRegPayload := &payload.TransferAsset{}

	return &Transaction{
		TxType:        TransferAsset,
		Payload:       assetRegPayload,
		Attributes:    []*TxAttribute{},
		UTXOInputs:    inputs,
		BalanceInputs: []*BalanceTxInput{},
		Outputs:       outputs,
		Programs:      []*program.Program{},
	}, nil
}

//initial a new transaction with record payload
func NewRecordTransaction(recordType string, recordData []byte) (*Transaction, error) {
	//TODO: check arguments
	recordPayload := &payload.Record{
		RecordType: recordType,
		RecordData: recordData,
	}

	return &Transaction{
		TxType:        Record,
		Payload:       recordPayload,
		Attributes:    []*TxAttribute{},
		UTXOInputs:    []*UTXOTxInput{},
		BalanceInputs: []*BalanceTxInput{},
		Programs:      []*program.Program{},
	}, nil
}

func NewPrivacyPayloadTransaction(fromPrivKey []byte, fromPubkey *crypto.PubKey, toPubkey *crypto.PubKey, payloadType payload.EncryptedPayloadType, data []byte) (*Transaction, error) {
	privacyPayload := &payload.PrivacyPayload{
		PayloadType: payloadType,
		EncryptType: payload.ECDH_AES256,
		EncryptAttr: &payload.EcdhAes256{
			FromPubkey: fromPubkey,
			ToPubkey:   toPubkey,
		},
	}
	privacyPayload.Payload, _ = privacyPayload.EncryptAttr.Encrypt(data, fromPrivKey)

	return &Transaction{
		TxType:        PrivacyPayload,
		Payload:       privacyPayload,
		Attributes:    []*TxAttribute{},
		UTXOInputs:    []*UTXOTxInput{},
		BalanceInputs: []*BalanceTxInput{},
		Programs:      []*program.Program{},
	}, nil
}
func NewDataFileTransaction(path string, fileName string, note string, issuer *crypto.PubKey) (*Transaction, error) {
	//TODO: check arguments
	DataFilePayload := &payload.DataFile{
		IPFSPath: path,
		Filename: fileName,
		Note:     note,
		Issuer:   issuer,
	}

	return &Transaction{
		TxType:        DataFile,
		Payload:       DataFilePayload,
		Attributes:    []*TxAttribute{},
		UTXOInputs:    []*UTXOTxInput{},
		BalanceInputs: []*BalanceTxInput{},
		Programs:      []*program.Program{},
	}, nil
}

//initial a new transaction with publish payload
func NewDeployTransaction(fc *code.FunctionCode, programHash common.Uint160, name, codeversion, author, email, desp string, language types.LangType) (*Transaction, error) {
	//TODO: check arguments
	DeployCodePayload := &payload.DeployCode{
		Code:        fc,
		Name:        name,
		CodeVersion: codeversion,
		Author:      author,
		Email:       email,
		Description: desp,
		Language:    language,
		ProgramHash: programHash,
	}

	return &Transaction{
		TxType:        DeployCode,
		Payload:       DeployCodePayload,
		Attributes:    []*TxAttribute{},
		UTXOInputs:    []*UTXOTxInput{},
		BalanceInputs: []*BalanceTxInput{},
		Programs:      []*program.Program{},
	}, nil
}

//initial a new transaction with invoke payload
func NewInvokeTransaction(fc []byte, codeHash common.Uint160, programhash common.Uint160) (*Transaction, error) {
	//TODO: check arguments
	InvokeCodePayload := &payload.InvokeCode{
		Code:        fc,
		CodeHash:    codeHash,
		ProgramHash: programhash,
	}

	return &Transaction{
		TxType:        InvokeCode,
		Payload:       InvokeCodePayload,
		Attributes:    []*TxAttribute{},
		UTXOInputs:    []*UTXOTxInput{},
		BalanceInputs: []*BalanceTxInput{},
		Programs:      []*program.Program{},
	}, nil
}
