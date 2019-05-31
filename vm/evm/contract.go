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

package evm

import (
	"DNA/common"
	"math/big"
)


type Contract struct {
	Caller common.Uint160
	Code []byte
	CodeHash common.Uint160
	Input []byte
	value *big.Int
	jumpdest destinations
}

func NewContract(caller common.Uint160) *Contract {
	var contract Contract
	contract.value = new(big.Int)
	contract.jumpdest = make(destinations, 0)
	contract.Caller = caller
	return &contract
}

func (c *Contract) GetOp(n uint64) OpCode {
	return OpCode(c.GetByte(n))
}

func (c *Contract) GetByte(n uint64) byte {
	if n < uint64(len(c.Code)) {
		return c.Code[n]
	}
	return 0
}

func (c *Contract) SetCode(code []byte, codeHash common.Uint160) {
	c.Code = code
	c.CodeHash = codeHash
}

func (c *Contract) SetCallCode(code, input []byte, codeHash common.Uint160) {
	c.Code = code
	c.Input = input
	c.CodeHash = codeHash
}
