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
	. "DNA/vm/evm/common"
	"DNA/vm/evm/crypto"
	"crypto/sha256"
	"github.com/golang/crypto/ripemd160"
)

type PrecompiledContract interface {
	Run(input []byte) ([]byte, error)
}

func allZero(b []byte) bool {
	for _, byte := range b {
		if byte != 0 {
			return false
		}
	}
	return true
}

var PrecompiledContracts = map[common.Uint160]PrecompiledContract{
	common.BytesToUint160([]byte{1}): &ecrecover{},
	common.BytesToUint160([]byte{2}): &sha256hash{},
	common.BytesToUint160([]byte{3}): &ripemd160hash{},
	common.BytesToUint160([]byte{4}): &dataCopy{},
}

type ecrecover struct{}

func (c *ecrecover) Run(in []byte) ([]byte, error) {
	const ecRecoverInputLength = 128
	in = RightPadBytes(in, ecRecoverInputLength)

	r := new(big.Int).SetBytes(in[64:96])
	s := new(big.Int).SetBytes(in[96:128])
	v := in[63] - 27

	if !allZero(in[32:63]) || !crypto.ValidateSignatureValues(v, r, s, false) {
		return nil, nil
	}
	pubKey, err := crypto.Ecrecover(in[:32], append(in[64:128], v))
	if err != nil {
		return nil, nil
	}

	return LeftPadBytes(crypto.Keccak256(pubKey[1:])[12:], 32), nil
}

type sha256hash struct{}

func (c *sha256hash) Run(in []byte) ([]byte, error) {
	h := sha256.Sum256(in)
	return h[:], nil
}

type ripemd160hash struct{}

func (c *ripemd160hash) Run(in []byte) ([]byte, error) {
	ripemd := ripemd160.New()
	ripemd.Write(in)
	return LeftPadBytes(ripemd.Sum(nil), 32), nil
}

type dataCopy struct{}
func (c *dataCopy) Run(in []byte) ([]byte, error) {
	return in, nil
}


