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

package errors

import (
	"errors"
)

const callStackDepth = 10

type DetailError interface {
	error
	ErrCoder
	CallStacker
	GetRoot()  error
}


func  NewErr(errmsg string) error {
	return errors.New(errmsg)
}

func NewDetailErr(err error,errcode ErrCode,errmsg string) DetailError{
	if err == nil {return nil}

	dnaerr, ok := err.(dnaError)
	if !ok {
		dnaerr.root = err
		dnaerr.errmsg = err.Error()
		dnaerr.callstack = getCallStack(0, callStackDepth)
		dnaerr.code = errcode

	}
	if errmsg != "" {
		dnaerr.errmsg = errmsg + ": " + dnaerr.errmsg
	}


	return dnaerr
}

func RootErr(err error) error {
	if err, ok := err.(DetailError); ok {
		return err.GetRoot()
	}
	return err
}



