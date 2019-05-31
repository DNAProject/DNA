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

import (
	"io"
	ser "DNA/common/serialization"
)

type ChangeView struct {
	msgData ConsensusMessageData
	NewViewNumber byte
}

func (cv *ChangeView) Serialize(w io.Writer)error{
	cv.msgData.Serialize(w)
	w.Write([]byte{cv.NewViewNumber})
	return nil
}

//read data to reader
func (cv *ChangeView) Deserialize(r io.Reader) error{
	 cv.msgData.Deserialize(r)
	viewNum,err := ser.ReadBytes(r,1)
	if err != nil {
		return err
	}
	cv.NewViewNumber = viewNum[0]
	return nil
}

func (cv *ChangeView) Type() ConsensusMessageType{
	return cv.ConsensusMessageData().Type
}

func (cv *ChangeView) ViewNumber() byte{
	return cv.msgData.ViewNumber
}

func (cv *ChangeView) ConsensusMessageData() *ConsensusMessageData{
	return &(cv.msgData)
}

