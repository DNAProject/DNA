package payload

import "io"

type DestroyUTXO struct {
}

func (a *DestroyUTXO) Data(version byte) []byte {
	//TODO: implement TransferAsset.Data()
	return []byte{0}

}

func (a *DestroyUTXO) Serialize(w io.Writer, version byte) error {
	return nil
}

func (a *DestroyUTXO) Deserialize(r io.Reader, version byte) error {
	return nil
}
