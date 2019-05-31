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

package events

import (
	"testing"
	"fmt"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent()

	var subscriber1 EventFunc = func(v interface{}){
		fmt.Println("subscriber1 event func.")
	}

	var subscriber2 EventFunc = func(v interface{}){
		fmt.Println("subscriber2 event func.")
	}

	fmt.Println("Subscribe...")
	sub1 := event.Subscribe(EventReplyTx,subscriber1)
	event.Subscribe(EventSaveBlock,subscriber2)

	fmt.Println("Notify...")
	event.Notify(EventReplyTx,nil)

	fmt.Println("Notify All...")
	event.NotifyAll()

	event.UnSubscribe(EventReplyTx,sub1)
	fmt.Println("Notify All after unsubscribe sub1...")
	event.NotifyAll()

}
