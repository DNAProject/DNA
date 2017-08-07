package payload

import (
	"DNA/common/serialization"
	"io"
)

const IssueAssetPayloadVersion byte = 0x00

type IssueAsset struct {
	Nonce uint64
}

func (a *IssueAsset) Data(version byte) []byte {
	//TODO: implement IssueAsset.Data()
	return []byte{0}

}

func (a *IssueAsset) Serialize(w io.Writer, version byte) error {
	return serialization.WriteUint64(w, a.Nonce)
}

func (a *IssueAsset) Deserialize(r io.Reader, version byte) error {
	var err error
	a.Nonce, err = serialization.ReadUint64(r)
	return err
}
