package transaction

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTxAttribute(t *testing.T) {
	url := []byte("http:\\www.onchain.com")
	tx := NewTxAttribute(DescriptionUrl, url)
	b := new(bytes.Buffer)
	tx.Serialize(b)
	fmt.Println("Serialize complete")

	ty := NewTxAttribute(ContractHash, nil)
	ty.Deserialize(b)
	fmt.Println("Deserialize complete.")
	fmt.Printf("Print: Usage= :0x%x,Url Date: %q\n", ty.Usage, ty.Date)
}
