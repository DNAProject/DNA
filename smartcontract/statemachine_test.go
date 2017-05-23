package smartcontract

import (
	"testing"
)

func TestNewStateMachine(t *testing.T) {
	stateMachine := NewStateMachine()
	m := stateMachine.GetServiceMap()

	for k, v := range m {
		t.Log(k, v)
	}
}
