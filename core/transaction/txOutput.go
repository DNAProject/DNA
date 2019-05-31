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

package transaction

import (
	"DNA/common"
	"io"
	"bytes"
)

type TxOutput struct {
	AssetID     common.Uint256
	Value       common.Fixed64
	ProgramHash common.Uint160
}

func (o *TxOutput) Serialize(w io.Writer) {
	o.AssetID.Serialize(w)
	o.Value.Serialize(w)
	o.ProgramHash.Serialize(w)
}

func (o *TxOutput) Deserialize(r io.Reader) {
	o.AssetID.Deserialize(r)
	o.Value.Deserialize(r)
	o.ProgramHash.Deserialize(r)
}

func (ui *TxOutput) ToArray() ([]byte) {
	b := new(bytes.Buffer)
	ui.Serialize(b)
	return b.Bytes()
}