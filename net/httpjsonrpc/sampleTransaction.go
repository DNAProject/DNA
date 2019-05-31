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

package httpjsonrpc

import (
	. "DNA/account"
	. "DNA/common"
	"DNA/common/log"
	. "DNA/core/asset"
	"DNA/core/contract"
	"DNA/core/signature"
	"DNA/core/transaction"
	"strconv"
)

const (
	ASSETPREFIX = "DNA"
)

func NewRegTx(rand string, index int, admin, issuer *Account) *transaction.Transaction {
	name := ASSETPREFIX + "-" + strconv.Itoa(index) + "-" + rand
	description := "description"
	asset := &Asset{name, description, byte(MaxPrecision), AssetType(Share), UTXO}
	amount := Fixed64(1000)
	controller, _ := contract.CreateSignatureContract(admin.PubKey())
	tx, _ := transaction.NewRegisterAssetTransaction(asset, amount, issuer.PubKey(), controller.ProgramHash)
	return tx
}

func SignTx(admin *Account, tx *transaction.Transaction) {
	signdate, err := signature.SignBySigner(tx, admin)
	if err != nil {
		log.Error(err, "signdate SignBySigner failed")
	}
	transactionContract, _ := contract.CreateSignatureContract(admin.PublicKey)
	transactionContractContext := contract.NewContractContext(tx)
	transactionContractContext.AddContract(transactionContract, admin.PublicKey, signdate)
	tx.SetPrograms(transactionContractContext.GetPrograms())
}
