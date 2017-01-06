package errors

import (
	"runtime"
	"errors"
)



const callStackDepth = 10

type dnaError struct {
	 errmsg string
	 callstack *CallStack
	 stackLen int
	 root error
 }

func (e dnaError) Error() string {
	return e.errmsg
}

func  NewErr(errmsg string) error {
	return errors.New(errmsg)
}

func RootErr(err error) error {
	if err, ok := err.(dnaError); ok {
		return err.root
	}
	return err
}

func NewDetailErr(err error,errmsg string) error{
	if err == nil {return nil}


	dnaerr, ok := err.(dnaError)
	if !ok {
		dnaerr.root = err
		dnaerr.errmsg = err.Error()
		dnaerr.callstack = getCallStack(0, callStackDepth)

	}
	if errmsg != "" {
		dnaerr.errmsg = errmsg + ": " + dnaerr.errmsg
	}
	return dnaerr
}


func getCallStack(skip int, depth int) (*CallStack){
	stacks := make([]uintptr, depth)
	stacklen := runtime.Callers(skip, stacks)

	return &CallStack{
		Stacks: stacks,
		Len: stacklen,
	}
}

