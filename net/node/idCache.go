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

package node

import (
	"DNA/common"
	"sync"
)

type idCache struct {
	sync.RWMutex
	list map[common.Uint256]bool
}

func (c *idCache) init() {
}

func (c *idCache) add() {
}

func (c *idCache) del() {
}

func (c *idCache) ExistedID(id common.Uint256) bool {
	// TODO
	return false
}
