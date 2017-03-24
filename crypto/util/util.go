package util

import (
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const (
	HASHLEN       = 32
	PRIVATEKEYLEN = 32
	PUBLICKEYLEN  = 32
	SIGNRLEN      = 32
	SIGNSLEN      = 32
	SIGNATURELEN  = 64
)

// InterfaceCrypto ---
type InterfaceCrypto struct {
	EccParams elliptic.CurveParams
	EccParamA *big.Int
	Curve     elliptic.Curve
}

// RandomNum Generate the "real" random number which can be used for crypto algorithm
func RandomNum(n int) ([]byte, error) {
	// TODO Get the random number from System urandom
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}
	return b, nil
}

func Hash256(value []byte) [HASHLEN]byte {
	var data [HASHLEN]byte
	digest := sha256.Sum256(value)
	copy(data[0:HASHLEN], digest[0:32])
	return data
}

func DoubleHash256(value []byte) [HASHLEN]byte {
	var data [HASHLEN]byte
	digest1 := sha256.Sum256(value)
	digest2 := sha256.Sum256(digest1[0:HASHLEN])
	copy(data[0:HASHLEN], digest2[0:HASHLEN])
	for i := 0; i < HASHLEN; i++ {
		digest1[i] = 0
		digest2[i] = 0
	}
	return data
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func Hash160(value []byte) []byte {
	md := ripemd160.New()
	io.WriteString(md, string(value))
	f := md.Sum(nil)
	return f
}
