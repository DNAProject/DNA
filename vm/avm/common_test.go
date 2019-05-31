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

package avm

import (
	"math/big"
	"testing"
	"fmt"
)

func TestCommon(t *testing.T) {
	i := ToBigInt(big.NewInt(1))
	t.Log("i", i)

	fmt.Println(ToArrayReverse([]byte{1, 2, 3}))
}

func ToArrayReverse(arr []byte) []byte {
	l := len(arr)
	x := make([]byte, 0)
	for i := l - 1; i >= 0 ;i--{
		x = append(x, arr[i])
	}
	return x
}
