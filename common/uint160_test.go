package common

import (
	"fmt"
	"testing"
)

func Test_ToScriptHash(t *testing.T) {
	ph, _ := ToScriptHash("AGX8Xw1HbWGFwozH4tH1ej14FXVdkeTTau")
	phstr := fmt.Sprintf("%x", ph.ToArray())
	if phstr != "082e502f35ec5cf8cc1209d0de00c550578911a7" {
		t.Error("ToScriptHash error!")
	}

	addr, _ := ph.ToAddress()
	if addr != "AGX8Xw1HbWGFwozH4tH1ej14FXVdkeTTau" {
		t.Error("ToScriptHash error!")
	}

	ph, _ = ToScriptHash("AQSqZUpRbf5KpxtATT5NFTLE2rE4Cd2MJq")
	phstr = fmt.Sprintf("%x", ph.ToArray())
	if phstr != "5f1f7caef4c5712fe63d09d81979ab7b36b6eda1" {
		t.Error("ToScriptHash error!")
	}

	addr, _ = ph.ToAddress()
	if addr != "AQSqZUpRbf5KpxtATT5NFTLE2rE4Cd2MJq" {
		t.Error("ToScriptHash error!")
	}
}
