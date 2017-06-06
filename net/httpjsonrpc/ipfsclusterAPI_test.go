package httpjsonrpc

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestAddFileIPFS(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=test bs=1024 count=1000")
	cmd.Run()
	ref, err := AddFileIPFS("test", nil, true)
	if err != nil {
		t.Fatalf("AddFileIPFS error:%s", err.Error())
	}
	os.Remove("test")
	fmt.Printf("ipfs path=%s\n", ref)
}
func TestGetFileIPFS(t *testing.T) {
	ref := "QmVHzLjYvp4bposJDD2PNeJ9PAFixyQu3oFj6gqipgsukX"
	err := GetFileIPFS(ref, "testOut")
	if err != nil {
		t.Fatalf("GetFileIPFS error:%s", err.Error())
	}
	//os.Remove("testOut")
}
