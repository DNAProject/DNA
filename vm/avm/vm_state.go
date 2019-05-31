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

package avm

type VMState byte

const (
	NONE  VMState = 0
	HALT  VMState = 1 << 0
	FAULT VMState = 1 << 1
	BREAK VMState = 1 << 2

	INSUFFICIENT_RESOURCE VMState = 1 << 4
)
