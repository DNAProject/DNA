package vm

import (
	"testing"
	"math/big"
)

func TestCommon(t *testing.T) {
	i := ToBigInt(big.NewInt(1))
	t.Log("i", i)
}
