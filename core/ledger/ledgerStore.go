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

package ledger

import (
	. "DNA/common"
	"DNA/core/account"
	. "DNA/core/asset"
	tx "DNA/core/transaction"
	"DNA/crypto"
	"DNA/smartcontract/states"
)

// ILedgerStore provides func with store package.
type ILedgerStore interface {
	//TODO: define the state store func
	SaveBlock(b *Block, ledger *Ledger) error
	GetBlock(hash Uint256) (*Block, error)
	BlockInCache(hash Uint256) bool
	GetBlockHash(height uint32) (Uint256, error)
	InitLedgerStore(ledger *Ledger) error
	IsDoubleSpend(tx *tx.Transaction) bool

	//SaveHeader(header *Header,ledger *Ledger) error
	AddHeaders(headers []Header, ledger *Ledger) error
	GetHeader(hash Uint256) (*Header, error)

	GetTransaction(hash Uint256) (*tx.Transaction, error)

	SaveAsset(assetid Uint256, asset *Asset) error
	GetAsset(hash Uint256) (*Asset, error)

	GetContract(codeHash Uint160) ([]byte, error)
	GetStorage(key []byte) ([]byte, error)
	GetAccount(programHash Uint160) (*account.AccountState, error)
	GetAssetState(assetId Uint256) (*states.AssetState, error)

	GetCurrentBlockHash() Uint256
	GetCurrentHeaderHash() Uint256
	GetHeaderHeight() uint32
	GetHeight() uint32
	GetHeaderHashByHeight(height uint32) Uint256

	GetBookKeeperList() ([]*crypto.PubKey, []*crypto.PubKey, error)
	InitLedgerStoreWithGenesisBlock(genesisblock *Block, defaultBookKeeper []*crypto.PubKey) (uint32, error)

	GetQuantityIssued(assetid Uint256) (Fixed64, error)

	GetUnspent(txid Uint256, index uint16) (*tx.TxOutput, error)
	ContainsUnspent(txid Uint256, index uint16) (bool, error)
	GetUnspentFromProgramHash(programHash Uint160, assetid Uint256) ([]*tx.UTXOUnspent, error)
	GetUnspentsFromProgramHash(programHash Uint160) (map[Uint256][]*tx.UTXOUnspent, error)
	GetAssets() map[Uint256]*Asset

	IsTxHashDuplicate(txhash Uint256) bool
	IsBlockInStore(hash Uint256) bool
	Close()
}
