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

package crypto

import (
	. "DNA/common"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {

	var data []Uint256
	a1 := Uint256(sha256.Sum256([]byte("a")))
	a2 := Uint256(sha256.Sum256([]byte("b")))
	a3 := Uint256(sha256.Sum256([]byte("c")))
	a4 := Uint256(sha256.Sum256([]byte("d")))
	a5 := Uint256(sha256.Sum256([]byte("e")))
	data = append(data, a1)
	data = append(data, a2)
	data = append(data, a3)
	data = append(data, a4)
	data = append(data, a5)
	x, _ := ComputeRoot(data)
	fmt.Printf("[Root Hash]:%x\n", x)

}
