package errors

import (
	"testing"
	"errors"
	"fmt"
)

var (
	TestRootError = errors.New("Test Root Error Msg.")
)



func TestNewDetailErr(t *testing.T) {
	e := NewDetailErr(TestRootError,"Test New Detail Error")
	if e == nil {
		t.Fatal("NewDetailErr should not return nil.")
	}
	fmt.Println(e.Error())

	msg := CallStacksString(GetCallStacks(e))

	fmt.Println(msg)

	if msg == ""{
		t.Errorf("CallStacksString should not return empty msg.")
	}

	rooterr := RootErr(e)

	fmt.Println("Root: ",rooterr.Error())

	fmt.Println("TestNewDetailErr End.")
}

