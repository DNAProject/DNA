package errors

import (
	"fmt"
	"bytes"
	"runtime"
)

type CallStack struct {
	Stacks []uintptr
	Len int
}

func GetCallStacks(err error) *CallStack {
	if dnaerr, ok := err.(dnaError); ok {
		return dnaerr.callstack
	}
	return nil
}


func CallStacksString(call *CallStack) string  {
	buf := bytes.Buffer{}
	if call == nil {
		return fmt.Sprintf("No call stack available")
	}

	for i := 0; i < call.Len; i++ {
		f := runtime.FuncForPC(call.Stacks[i])
		file, line := f.FileLine(call.Stacks[i])
		buf.WriteString(fmt.Sprintf("%s:%d - %s\n", file, line, f.Name()))
	}

	return fmt.Sprintf("%s", buf.Bytes())
}