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

package errors

import (
	"testing"
	"errors"
	"fmt"
)

var (
	TestRootError = errors.New("Test Root Error Msg.")
)



func TestNewDetailErr(t *testing.T) {
	e := NewDetailErr(TestRootError,ErrUnknown,"Test New Detail Error")
	if e == nil {
		t.Fatal("NewDetailErr should not return nil.")
	}
	fmt.Println(e.Error())

	msg := CallStacksString(GetCallStacks(e))

	fmt.Println(msg)

	if msg == ""{
		t.Errorf("CallStacksString should not return empty msg.")
	}

	rooterr := RootErr(e)
	fmt.Println("Root: ",rooterr.Error())

	code := ErrerCode(e)
	fmt.Println("Code: ",code.Error())

	fmt.Println("TestNewDetailErr End.")
}

