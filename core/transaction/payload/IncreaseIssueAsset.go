package payload

import (
	"DNA/common"
	. "DNA/errors"
	"bytes"
	"io"
)

type IncreaseIssueAsset struct {
	AssetID common.Uint256
	Amount  common.Fixed64
}

func (self *IncreaseIssueAsset) Data() []byte {
	var buf bytes.Buffer
	self.AssetID.Serialize(&buf)
	self.Amount.Serialize(&buf)

	return buf.Bytes()
}

func (self *IncreaseIssueAsset) Serialize(w io.Writer) error {
	_, err := w.Write(self.Data())

	return err
}

func (self *IncreaseIssueAsset) Deserialize(r io.Reader) error {
	err := self.AssetID.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[IncreaseIssueAsset], AssertID Deserialize failed.")
	}
	err = self.Amount.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[IncreaseIssueAsset], Amount Deserialize failed.")
	}

	return nil
}
