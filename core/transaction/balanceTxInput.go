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
)


type BalanceTxInput struct {
	AssetID common.Uint256
	Value common.Fixed64
	ProgramHash common.Uint160
}

func (bi *BalanceTxInput) Serialize(w io.Writer)  {
	bi.AssetID.Serialize(w)
	bi.Value.Serialize(w)
	bi.ProgramHash.Serialize(w)
}

func (bi *BalanceTxInput) Deserialize(r io.Reader) error  {
	err := bi.AssetID.Deserialize(r)
	if err != nil {return err}

	err = bi.Value.Deserialize(r)
	if err != nil {return err}

	err = bi.ProgramHash.Deserialize(r)
	if err != nil {return err}

	return nil
}
