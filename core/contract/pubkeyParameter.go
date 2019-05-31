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

package contract

type PubkeyParameter struct {
	PubKey    string
	Parameter string
}

type ParameterIndex struct {
	Parameter []byte
	Index     int
}

type ParameterIndexSlice []ParameterIndex

func (p ParameterIndexSlice) Len() int           { return len(p) }
func (p ParameterIndexSlice) Less(i, j int) bool { return p[i].Index < p[j].Index }
func (p ParameterIndexSlice) Swap(i, j int)      { p[i].Index, p[j].Index = p[j].Index, p[i].Index }
