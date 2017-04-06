package signature

import (
	"DNA/common"
	"DNA/common/log"
	"DNA/core/contract/program"
	"DNA/crypto"
	. "DNA/errors"
	"DNA/vm/interfaces"
	"bytes"
	"crypto/sha256"
	"io"
)

//SignableData describe the data need be signed.
type SignableData interface {
	interfaces.ISignableObject

	//Get the the SignableData's program hashes
	GetProgramHashes() ([]common.Uint160, error)

	SetPrograms([]*program.Program)

	GetPrograms() []*program.Program

	//TODO: add SerializeUnsigned
	SerializeUnsigned(io.Writer) error
}

func SignBySigner(data SignableData, signer Signer) ([]byte, error) {
	log.Trace()
	//fmt.Println("data",data)
	rtx, err := Sign(data, signer.PrivKey())
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Signature],SignBySigner failed.")
	}
	return rtx, nil
}

func GetHashData(data SignableData) []byte {
	b_buf := new(bytes.Buffer)
	data.SerializeUnsigned(b_buf)
	return b_buf.Bytes()
}

func GetHashForSigning(data SignableData) []byte {
	//TODO: GetHashForSigning
	temp := sha256.Sum256(GetHashData(data))
	return temp[:]
}

func Sign(data SignableData, prikey []byte) ([]byte, error) {
	// FIXME ignore the return error value
	signature, err := crypto.Sign(prikey, GetHashForSigning(data))
	if err != nil {
		return nil, NewDetailErr(err, ErrNoCode, "[Signature],Sign failed.")
	}
	return signature, nil
}
