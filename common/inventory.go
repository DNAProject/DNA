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

package common

type InventoryType byte

const (
	TRANSACTION	InventoryType = 0x01
	BLOCK		InventoryType = 0x02
	CONSENSUS	InventoryType = 0xe0
)

//TODO: temp inventory
type Inventory interface {
	//sig.SignableData
	Hash() Uint256
	Verify() error
	Type() InventoryType
}

