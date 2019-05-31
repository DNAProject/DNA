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

package util

import (
	"DNA/crypto/sm3"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	//"math/big"
)

const (
	HASHLEN       = 32
	PRIVATEKEYLEN = 32
	PUBLICKEYLEN  = 32
	SIGNRLEN      = 32
	SIGNSLEN      = 32
	SIGNATURELEN  = 64
	NEGBIGNUMLEN  = 33
)

type CryptoAlgSet struct {
	EccParams elliptic.CurveParams
	Curve     elliptic.Curve
}

// RandomNum Generate the "real" random number which can be used for crypto algorithm
func RandomNum(n int) ([]byte, error) {
	// TODO Get the random number from System urandom
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}
	return b, nil
}

func Hash(data []byte) [HASHLEN]byte {
	return sha256.Sum256(data)
}

func SM3(data []byte) [HASHLEN]byte {
	return sm3.Sum(data)
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func RIPEMD160(value []byte) []byte {
	//TODO: implement RIPEMD160

	return nil
}
