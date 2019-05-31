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

package store

type IIterator interface {
	Next() bool
	Prev() bool
	First() bool
	Last() bool
	Seek(key []byte) bool
	Key() []byte
	Value() []byte
	Release()
}

type IStore interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	NewBatch() error
	BatchPut(key []byte, value []byte) error
	BatchDelete(key []byte) error
	BatchCommit() error
	Close() error
	NewIterator(prefix []byte) IIterator
}
