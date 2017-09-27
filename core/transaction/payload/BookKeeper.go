package payload

import (
	"DNA/common/serialization"
	"DNA/crypto"
	. "DNA/errors"
	"bytes"
	"io"
)

const BookKeeperPayloadVersion byte = 0x00

type BookKeeperAction byte

const (
	BookKeeperAction_ADD BookKeeperAction = 0
	BookKeeperAction_SUB BookKeeperAction = 1
)

type BookKeeper struct {
	PubKey *crypto.PubKey
	Action BookKeeperAction
	Cert   []byte
	Issuer *crypto.PubKey
}

func (self *BookKeeper) Data(version byte) []byte {
	var buf bytes.Buffer
	self.PubKey.Serialize(&buf)
	buf.WriteByte(byte(self.Action))
	serialization.WriteVarBytes(&buf, self.Cert)
	self.Issuer.Serialize(&buf)

	return buf.Bytes()
}

func (self *BookKeeper) Serialize(w io.Writer, version byte) error {
	_, err := w.Write(self.Data(version))

	return err
}

func (self *BookKeeper) Deserialize(r io.Reader, version byte) error {
	self.PubKey = new(crypto.PubKey)
	err := self.PubKey.DeSerialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[BookKeeper], PubKey Deserialize failed.")
	}
	var p [1]byte
	n, err := r.Read(p[:])
	if n == 0 {
		return NewDetailErr(err, ErrNoCode, "[BookKeeper], Action Deserialize failed.")
	}
	self.Action = BookKeeperAction(p[0])
	self.Cert, err = serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[BookKeeper], Cert Deserialize failed.")
	}
	self.Issuer = new(crypto.PubKey)
	err = self.Issuer.DeSerialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[BookKeeper], Issuer Deserialize failed.")
	}

	return nil
}
