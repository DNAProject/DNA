package sm2

import (
	"GoOnchain/crypto/util"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

// HASHLEN ---
const (
	HASHLEN       = 32
	PRIVATEKEYLEN = 32
	PUBLICKEYLEN  = 32
	SIGNATURELEN  = 64
)

var Sm2ParamA *big.Int
var BaseG *ECPoint
var Infinity *ECPoint
var EcParams *elliptic.CurveParams

// Init ---
func Init(Crypto *util.InterfaceCrypto) {
	EcParams = new(elliptic.CurveParams)
	EcParams.P, _ = new(big.Int).SetString("8542D69E4C044F18E8B92435BF6FF7DE457283915C45517D722EDB8B08F1DFC3", 16)
	EcParams.N, _ = new(big.Int).SetString("8542D69E4C044F18E8B92435BF6FF7DD297720630485628D5AE74EE7C32E79B7", 16)
	Sm2ParamA, _ = new(big.Int).SetString("787968B4FA32C3FD2417842E73BBFEFF2F3C848B6831D7E0EC65228B3937E498", 16)
	EcParams.B, _ = new(big.Int).SetString("63E4C6D3B23B0C849CF84241484BFE48F61D59A5B16BA06E6E12D1DA27C5249A", 16)

	EcParams.Gx, _ = new(big.Int).SetString("421DEBD61B62EAB6746434EBC3CC315E32220B3BADD50BDC4C4E6C147FEDD43D", 16)
	EcParams.Gy, _ = new(big.Int).SetString("0680512BCBB42C07D47349D2153B70C4E5D7FDFCBFA36EA1A85841B9E46E09A2", 16)
	EcParams.BitSize = 256
	EcParams.Name = "sm2"

	Crypto.EccParamA.Set(Sm2ParamA)
	Crypto.EccParams = *EcParams

	Infinity = &ECPoint{nil, nil, EcParams}

	GX := &ECFieldElement{EcParams.Gx, EcParams}
	GY := &ECFieldElement{EcParams.Gy, EcParams}
	BaseG = &ECPoint{GX, GY, EcParams}

	return
}

// GenKeyPair generates a public & private key pair
func GenKeyPair(Crypto *util.InterfaceCrypto) ([]byte, *big.Int, *big.Int, error) {
	pubKey := NewECPoint()

	dBytes, _ := util.RandomNum(PRIVATEKEYLEN)
	pubKey.Mul(BaseG, dBytes)

	//PrintHex("prikey", dBytes, len(dBytes))
	//PrintHex("pbkeyX", pubKey.X.value.Bytes(), len(pubKey.X.value.Bytes()))
	//PrintHex("pbkeyY", pubKey.Y.value.Bytes(), len(pubKey.Y.value.Bytes()))

	return dBytes, pubKey.X.value, pubKey.Y.value, nil
}

// CaculateE ---
func CaculateE(curveN *big.Int, msg []byte) *big.Int {
	msgBitLen := len(msg) * 8

	trunc := new(big.Int).SetBytes(msg)

	if curveN.BitLen() < msgBitLen {
		trunc.Rsh(trunc, uint(msgBitLen-curveN.BitLen()))
	}
	return trunc
}

// Sign process:
// 1. choose an integer num k between 1 and n - 1.
// 2. compute point = k * BaseG.
// 3. compute r = point.X mod n, if r is equal 0 goto step 1.
// 4. compute s = (z + r × d) / k mod n, if s is equal 0 goto step 1.
//    z is the message, d is private key.
// 5. (r, s) is the signature.
// Notice: In step 4, if the resault is fraction(it is usally be in actually),
// It must be numerator multiply the denominator's Inverse
func Sign(Crypto *util.InterfaceCrypto, priKey []byte, data []byte) (*big.Int, *big.Int, error) {
	if nil == priKey {
		fmt.Println("prikey is nil")
	}

	e := CaculateE(EcParams.N, data)
	priD := new(big.Int).SetBytes(priKey)

	k := big.NewInt(0)
	r := big.NewInt(0)
	s := big.NewInt(0)

	for {
		for {
			for {
				randK := make([]byte, EcParams.BitSize/8)
				_, err := rand.Read(randK)
				if err != nil {
					return nil, nil, err
				}
				//PrintHex("ranK", randK, len(randK))
				k.SetBytes(randK)
				if k.Sign() != 0 && k.Cmp(EcParams.N) < 0 {
					break
				}
			}

			kG := NewECPoint()
			kG.Mul(BaseG, k.Bytes())

			r.Mod(kG.X.value, EcParams.N)
			if r.Sign() != 0 {
				break
			}
		}

		tmp := big.NewInt(0)
		tmp.Mul(priD, r)
		tmp.Add(tmp, e)

		s.ModInverse(k, EcParams.N)
		s.Mul(s, tmp)
		s.Mod(s, EcParams.N)

		tmp.Div(EcParams.N, big.NewInt(2))
		if 1 == s.Cmp(tmp) {
			s.Sub(EcParams.N, s)
		}

		if s.Sign() != 0 {
			break
		}
	}
	retR := big.NewInt(0)
	retS := big.NewInt(0)

	// r and s must between 1 and N - 1
	if r.Sign() < 1 {
		retR.Add(EcParams.P, r)
	} else {
		retR.Set(r)
	}

	if s.Sign() < 1 {
		retS.Add(EcParams.P, s)
	} else {
		retS.Set(s)
	}

	//PrintHex("r", sig[0].Bytes(), len(sig[0].Bytes()))
	//PrintHex("s", sig[1].Bytes(), len(sig[1].Bytes()))
	return retR, retS, nil
}

// Verify process:
// 1. verify that r and s are between 1 and n - 1.
// 2. compute w = s mod n.
// 3. compute u = z × w mod n.
// 4. compute v = r × w mod n.
// 5. compute point = uG + vQ (Q is the publickey point).
// 6. check that x mod n == r, otherwise verify failed.
func Verify(Crypto *util.InterfaceCrypto, X *big.Int, Y *big.Int, data []byte, r, s *big.Int) (bool, error) {
	if r.Sign() < 1 || s.Sign() < 1 || r.Cmp(EcParams.N) >= 0 || s.Cmp(EcParams.N) >= 0 {
		return false, errors.New("signature is invalid")
	}

	w := big.NewInt(0)
	u := big.NewInt(0)
	v := big.NewInt(0)

	e := CaculateE(EcParams.N, data)
	w.ModInverse(s, EcParams.N)

	u.Mul(e, w)
	u.Mod(u, EcParams.N)

	v.Mul(r, w)
	v.Mod(v, EcParams.N)

	pub := NewECPoint()
	pub.X.value.Set(X)
	pub.Y.value.Set(Y)

	point := SumOfTwoMultiplies(BaseG, u, pub, v)

	verR := big.NewInt(0)
	verR.Mod(point.X.value, EcParams.N)

	return (0 == verR.Cmp(r)), nil
}
