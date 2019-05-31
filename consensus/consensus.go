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

package consensus

import (
	"fmt"
	"DNA/common/log"
	"time"
)

type ConsensusService interface {
	Start() error
	Halt() error
}


func Log(message string){
	logMsg := fmt.Sprintf("[%s] %s" ,time.Now().Format("02/01/2006 15:04:05"),message)
	fmt.Println(logMsg)
	log.Info(logMsg)
}

