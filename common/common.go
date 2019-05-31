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

import (
	"DNA/common/log"
	. "DNA/errors"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/golang/crypto/ripemd160"
	"io"
	"math/rand"
	"os"
)

func ToCodeHash(code []byte) (Uint160, error) {
	temp := sha256.Sum256(code)
	md := ripemd160.New()
	io.WriteString(md, string(temp[:]))
	f := md.Sum(nil)

	hash, err := Uint160ParseFromBytes(f)
	if err != nil {
		return Uint160{}, NewDetailErr(errors.New("[Common] , ToCodeHash err."), ErrNoCode, "")
	}
	return hash, nil
}

func GetNonce() uint64 {
	log.Debug()
	// Fixme replace with the real random number generator
	nonce := uint64(rand.Uint32())<<32 + uint64(rand.Uint32())
	return nonce
}

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, tmp)
	return bytesBuffer.Bytes()
}

func BytesToInt16(b []byte) int16 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int16
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int16(tmp)
}

func BytesToInt(b [] byte) []int {
	i := make([]int, len(b))
	for k,v := range b {
		i[k] = int(v)
	}
	return i
}

func IsEqualBytes(b1 []byte, b2 []byte) bool {
	len1 := len(b1)
	len2 := len(b2)
	if len1 != len2 {
		return false
	}

	for i := 0; i < len1; i++ {
		if b1[i] != b2[i] {
			return false
		}
	}

	return true
}

func ToHexString(data []byte) string {
	return hex.EncodeToString(data)
}

func HexToBytes(value string) ([]byte, error) {
	return hex.DecodeString(value)
}

func ToArrayReverse(arr []byte) []byte {
	l := len(arr)
	x := make([]byte, 0)
	for i := l - 1; i >= 0 ;i--{
		x = append(x, arr[i])
	}
	return x
}

func BytesReverse(u []byte) []byte {
	for i, j := 0, len(u)-1; i < j; i, j = i+1, j-1 {
		u[i], u[j] = u[j], u[i]
	}
	return u
}

func HexToBytesReverse(value string) ([]byte, error) {
	u, err := hex.DecodeString(value)
	if err != nil {
		return u, err
	}
	return BytesReverse(u), err
}

func ClearBytes(arr []byte, len int) {
	for i := 0; i < len; i++ {
		arr[i] = 0
	}
}

func CompareHeight(blockHeight uint64, heights []uint64) bool {
	for _, height := range heights {
		if blockHeight < height {
			return false
		}
	}
	return true
}

func GetUint16Array(source []byte) ([]uint16, error) {
	if source == nil {
		return nil, NewDetailErr(errors.New("[Common] , GetUint16Array err, source = nil"), ErrNoCode, "")
	}

	if len(source)%2 != 0 {
		return nil, NewDetailErr(errors.New("[Common] , GetUint16Array err, length of source is odd."), ErrNoCode, "")
	}

	dst := make([]uint16, len(source)/2)
	for i := 0; i < len(source)/2; i++ {
		dst[i] = uint16(source[i*2]) + uint16(source[i*2+1])*256
	}

	return dst, nil
}

func ToByteArray(source []uint16) []byte {
	dst := make([]byte, len(source)*2)
	for i := 0; i < len(source); i++ {
		dst[i*2] = byte(source[i] % 256)
		dst[i*2+1] = byte(source[i] / 256)
	}

	return dst
}

func SliceRemove(slice []uint32, h uint32) []uint32 {
	for i, v := range slice {
		if v == h {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func FileExisted(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
