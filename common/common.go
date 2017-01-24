package common

import (
	"bytes"
	"encoding/binary"
)

func ToCodeHash(code []byte) Uint160{
	//TODO: ToCodeHash
	return Uint160{}
}

func  GetNonce() uint64{
	//TODO: GetNonce()
	return 0;
}

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, tmp)
	return bytesBuffer.Bytes()
}

func IsEqualBytes(b1 []byte,b2 []byte) bool {
	len1 := len(b1)
	len2 := len(b2)
	if len1 != len2 {return false}

	for i:=0; i<len1; i++ {
		if b1[i] != b2[i] {return false}
	}

	return true
}

func ToHexString(data []byte) string{
	//TODO: ToHexString
	return string(data)
}

func HexToBytes(value string) []byte{
	//TODO: HexToBytes
	return nil
}