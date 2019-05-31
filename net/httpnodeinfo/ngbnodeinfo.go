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

package httpnodeinfo

import "strings"

type NgbNodeInfo struct {
	NgbId         string
	NgbType       string
	NgbAddr       string
	HttpInfoAddr  string
	HttpInfoPort  int
	HttpInfoStart bool
}

type NgbNodeInfoSlice []NgbNodeInfo

func (n NgbNodeInfoSlice) Len() int {
	return len(n)
}

func (n NgbNodeInfoSlice) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n NgbNodeInfoSlice) Less(i, j int) bool {
	if 0 <= strings.Compare(n[i].HttpInfoAddr, n[j].HttpInfoAddr) {
		return false
	} else {
		return true
	}
}
