package utils

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"

)

type VmReader struct {
	reader     *bytes.Reader
	BaseStream []byte
}

func NewVmReader(b []byte) *VmReader {
	var vmreader VmReader
	vmreader.reader = bytes.NewReader(b)
	vmreader.BaseStream = b
	return &vmreader
}

func (r *VmReader) Reader() *bytes.Reader{
	return r.reader
}

func (r *VmReader) ReadByte() (byte,error) {
	byte, err := r.reader.ReadByte()
	return byte,err
}

func (r *VmReader) ReadBytes(count int) []byte {
	var bytes []byte
	for i := 0; i < count; i++ {
		d,_ := r.ReadByte()
		bytes = append(bytes, d)
	}
	return bytes
}

func (r *VmReader) ReadUint16() uint16 {
	b := r.ReadBytes(2)
	return binary.LittleEndian.Uint16(b)
}

func (r *VmReader) ReadUInt32() uint32 {
	b := r.ReadBytes(4)
	return binary.LittleEndian.Uint32(b)
}

func (r *VmReader) ReadUInt64() uint64 {
	b := r.ReadBytes(8)
	return binary.LittleEndian.Uint64(b)
}

func (r *VmReader) ReadInt16() int16 {
	b := r.ReadBytes(2)
	bytesBuffer := bytes.NewBuffer(b)
	var vi int16
	binary.Read(bytesBuffer, binary.LittleEndian, &vi)
	return vi

}

func (r *VmReader) ReadInt32() int32 {
	b := r.ReadBytes(4)
	bytesBuffer := bytes.NewBuffer(b)
	var vi int32
	binary.Read(bytesBuffer, binary.LittleEndian, &vi)
	return vi
}

func (r *VmReader) Position() int {
	return int(r.reader.Size()) - r.reader.Len()
}

func (r *VmReader) Length() int {
	return r.reader.Len()
}

func (r *VmReader) Seek(offset int64, whence int) (int64, error) {
	return r.reader.Seek(offset, whence)
}

func (r *VmReader) ReadVarBytes(max int) []byte {
	n := int(r.ReadVarInt(uint64(max)))
	return r.ReadBytes(n)
}

func (r *VmReader) ReadVarInt(max uint64) uint64 {
	fb,_ := r.ReadByte()
	var value uint64

	switch fb {
	case 0xFD:
		value = uint64(r.ReadInt16())
	case 0xFE:
		value = uint64(r.ReadUInt32())
	case 0xFF:
		value = uint64(r.ReadUInt64())
	default:
		value = uint64(fb)
	}
	if value > max {
		return 0
	}
	return value
}

func (r *VmReader) ReadVarString() string{
	bs := r.ReadVarBytes(0X7fffffc7)
	return string(bs)
	////b := r.ReadBytes(252)
	//
	//bytes := [4]byte{1,2,3,4}
	//str := convert(bytes[:])
	//fmt.Println(str)
	//return hex.EncodeToString(bs[:])

}

func convert( b []byte ) string {
	s := make([]string,len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s,",")
}

