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

package httpjsonrpc

import (
	"DNA/common/log"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestAddFileIPFS(t *testing.T) {
	var path string = "./Log/"
	log.CreatePrintLog(path)
	cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=test bs=1024 count=1000")
	cmd.Run()
	ref, err := AddFileIPFS("test", true)
	if err != nil {
		t.Fatalf("AddFileIPFS error:%s", err.Error())
	}
	os.Remove("test")
	fmt.Printf("ipfs path=%s\n", ref)
}
func TestGetFileIPFS(t *testing.T) {
	var path string = "./Log/"
	log.CreatePrintLog(path)
	ref := "QmVHzLjYvp4bposJDD2PNeJ9PAFixyQu3oFj6gqipgsukX"
	err := GetFileIPFS(ref, "testOut")
	if err != nil {
		t.Fatalf("GetFileIPFS error:%s", err.Error())
	}
	//os.Remove("testOut")
}
