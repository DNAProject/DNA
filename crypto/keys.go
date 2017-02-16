package crypto

import (
	"io"
	"errors"
	"math/big"
	"crypto/rand"
	"crypto/ecdsa"
	"fmt"
)

type PubKey ECPoint

func (e *PubKey) Serialize(w io.Writer) {
	//TODO: implement PubKey.serialize
}

func (e *PubKey) DeSerialize(r io.Reader) error {
	//TODO
	return nil
}


func Reverse(data []byte) {

	len1 := len(data)

	for i := 0; i < len1/2; i++ {
		Tmp := data[i]
		data[i] = data[len1-1-i]
		data[len1-1-i] = Tmp
	}
}

func IsEven(k *big.Int) bool {
	z := big.NewInt(0)
	z.Mod(k, big.NewInt(2))
	if z.Int64() == 0 {
		return true
	} else {
		return false
	}
}

func (ep *PubKey) EncodePoint(commpressed bool) []byte {

	if ep.X == nil && ep.Y == nil {
		infinity := make([]byte, 1)
		fmt.Println("IsInfinity")
		return infinity
	}

	var data []byte

	if commpressed {
		data = make([]byte, 33)
	} else {
		data = make([]byte, 65)

		yBytes := ep.Y.Bytes()
		//ep.Y.value.Bytes()

		tmp := make([]byte, len(yBytes))
		copy(tmp, yBytes)
		Reverse(tmp)

		copy(data[65-len(yBytes):], tmp)
	}

	xBytes := ep.X.Bytes()

	tmp := make([]byte, len(xBytes))
	copy(tmp, xBytes)
	Reverse(tmp)

	copy(data[33-len(tmp):], tmp)

	if !commpressed {
		data[0] = 0x04
	} else {
		if IsEven(ep.Y) {
			data[0] = 0x02
		} else {
			data[0] = 0x03
		}
	}

	return data
}


func NewPubKey(prikey []byte) *PubKey{
       //TODO: NewPubKey
       return nil
}

func GenPrivKey() []byte {
	return nil
}

//FIXME, does the privkey need base58 encoding?
//This generates a public & private key pair
func GenKeyPair() ([]byte, PubKey, error) {
	pubkey := new(PubKey)
	privatekey := new(ecdsa.PrivateKey)
	privatekey, err := ecdsa.GenerateKey(Crypto.curve, rand.Reader)
	if err != nil {
		return nil, *pubkey, errors.New("Generate key pair error")
	}

	privkey, err := privatekey.D.MarshalText()
	pubkey.X = privatekey.PublicKey.X
	pubkey.Y = privatekey.PublicKey.Y
	return privkey, *pubkey, nil
}

func DecodePoint(encoded []byte) *PubKey{
	//TODO: DecodePoint
	return nil
}

type PubKeySlice []*PubKey

func (p PubKeySlice) Len() int           { return len(p) }
func (p PubKeySlice) Less(i, j int) bool {
	//TODO:PubKeySlice Less
	return false
}
func (p PubKeySlice) Swap(i, j int) {
	//TODO:PubKeySlice Swap
}
