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

package httprestful

import (
	. "DNA/common/config"
	"DNA/core/ledger"
	"DNA/events"
	"DNA/net/httprestful/common"
	. "DNA/net/httprestful/restful"
	. "DNA/net/protocol"
	"strconv"
)

func StartServer(n Noder) {
	common.SetNode(n)
	ledger.DefaultLedger.Blockchain.BCEvents.Subscribe(events.EventBlockPersistCompleted, SendBlock2NoticeServer)
	func() {
		rest := InitRestServer(common.CheckAccessToken)
		go rest.Start()
	}()
}

func SendBlock2NoticeServer(v interface{}) {

	if len(Parameters.NoticeServerUrl) == 0 || !common.CheckPushBlock() {
		return
	}
	go func() {
		req := make(map[string]interface{})
		req["Height"] = strconv.FormatInt(int64(ledger.DefaultLedger.Blockchain.BlockHeight), 10)
		req = common.GetBlockByHeight(req)

		repMsg, _ := common.PostRequest(req, Parameters.NoticeServerUrl)
		if repMsg[""] == nil {
			//TODO
		}
	}()
}
