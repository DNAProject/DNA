package ledger

import (
	"io"
	"bytes"
)

type Header struct {
	Blockdata *Blockdata
}

//Serialize the blockheader
func (h *Header) Serialize(w io.Writer) {
	h.Blockdata.Serialize(w)
	w.Write([]byte{'0'})

}

func (h *Header) Deserialize(r io.Reader) error {
	header := new(Blockdata)
	header.Deserialize(r)
	h.Blockdata = header
	var headerFlag [1]byte
	_, err := io.ReadFull(r, headerFlag[:])
	if err != nil {
		return err
	}

	return nil
}

func (h *Header) ToArray() []byte {
	b := new(bytes.Buffer)
	h.Serialize(b)
	return b.Bytes()
}
