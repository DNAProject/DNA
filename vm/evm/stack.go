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

package evm

import (
	"math/big"
	"fmt"
)

type Stack struct {
	data []*big.Int
}


func newstack() *Stack {
	return &Stack{data: make([]*big.Int, 0, 1024)}
}

func (s *Stack) Data() []*big.Int {
	return s.data
}

func (s *Stack) push(b *big.Int) {
	s.data = append(s.data, b)
}

func (s *Stack) dup(n int) {
	s.push(new(big.Int).Set(s.data[s.len() - n]))
}

func (s *Stack) pushN(bs ...*big.Int) {
	s.data = append(s.data, bs...)
}

func (s *Stack) pop() (ret *big.Int) {
	l := s.len() - 1
	ret = s.data[l]
	s.data = s.data[:l]
	return
}

func (s *Stack) len() int {
	return len(s.data)
}

func (s *Stack) swap(n int) {
	l := s.len()
	l1 := l-n-1
	l2 := l-1
	s.data[l1], s.data[l2] = s.data[l2], s.data[l1]
}

func (s *Stack) peek() *big.Int {
	return s.data[s.len() - 1]
}

func (s *Stack) Back(n int) *big.Int {
	return s.data[s.len() - n - 1]
}

func (s *Stack) require(n int) error {
	if s.len() < n {
		return fmt.Errorf("stack underflow (%d <=> %d)", len(s.data), n)
	}
	return nil
}

func (s *Stack) Print() {
	fmt.Println("### stack ###")
	if len(s.data) > 0 {
		for i, val := range s.Data() {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
