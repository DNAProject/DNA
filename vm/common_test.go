package vm

import (
	"testing"
	"fmt"
	"math/big"
)

func TestCommon(t *testing.T) {
	i := ToBigInt(big.NewInt(1))
	fmt.Println("i", i)
}
