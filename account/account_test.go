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

package account

import (
	"DNA/crypto"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestClient(t *testing.T) {
	t.Log("created client start!")
	crypto.SetAlg(crypto.P256R1)
	dir := "./data/"
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		t.Log("create dir ", dir, " error: ", err)
	} else {
		t.Log("create dir ", dir, " success!")
	}
	for i := 0; i < 10000; i++ {
		p := path.Join(dir, fmt.Sprintf("wallet%d.txt", i))
		fmt.Println("client path", p)
		CreateClient(p, []byte(DefaultPin))
	}
}
