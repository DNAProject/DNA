package transaction

import (
	"GoOnchain/common"
	"io"
)

type TxOutput struct {
	AssetID common.Uint256
	Value common.Fixed64
	ProgramHash common.Uint160
}


func (o *TxOutput) Serialize(w io.Writer)  {
	o.AssetID.Serialize(w)
	o.Value.Serialize(w)
	o.ProgramHash.Serialize(w)
}

func (tx *TxOutput) Deserialize(r io.Reader)  {
	//TODO；TxOutput Deserialize

}