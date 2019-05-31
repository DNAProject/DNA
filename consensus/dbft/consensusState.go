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

package dbft

type ConsensusState byte

const (
	Initial         ConsensusState = 0x00
	Primary         ConsensusState = 0x01
	Backup          ConsensusState = 0x02
	RequestSent     ConsensusState = 0x04
	RequestReceived ConsensusState = 0x08
	SignatureSent   ConsensusState = 0x10
	BlockGenerated  ConsensusState = 0x20
)

func (state ConsensusState) HasFlag(flag ConsensusState) bool {
	return (state & flag) == flag
}
