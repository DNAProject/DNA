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
	"testing"
	"math/big"
)

func TestOpAdd(t *testing.T) {
	evm := NewExecutionEngine(nil, nil, nil)
	evm.stack.push(big.NewInt(2))
	evm.stack.push(big.NewInt(3))
	opAdd(evm)
	t.Log("add execute resutl:", evm.stack.pop())
}

func TestOpMul(t *testing.T) {
	evm := NewExecutionEngine(nil, nil, nil)
	evm.stack.push(big.NewInt(2))
	evm.stack.push(big.NewInt(3))
	opMul(evm)
	t.Log("Mul execute resutl:", evm.stack.pop())
}

func TestOpSub(t *testing.T) {
	evm := NewExecutionEngine(nil, nil, nil)
	evm.stack.push(big.NewInt(2))
	evm.stack.push(big.NewInt(3))
	opSub(evm)
	t.Log("Mul execute resutl:", evm.stack.pop())
}

func TestOpDiv(t *testing.T) {
	evm := NewExecutionEngine(nil, nil, nil)
	evm.stack.push(big.NewInt(-2))
	evm.stack.push(big.NewInt(3))
	opDiv(evm)
	t.Log("Mul execute resutl:", evm.stack.pop())
}

func TestOpSdiv(t *testing.T) {
	evm := NewExecutionEngine(nil, nil, nil)
	evm.stack.push(big.NewInt(-2))
	evm.stack.push(big.NewInt(3))
	opSdiv(evm)
	t.Log("Mul execute resutl:", evm.stack.pop())
}