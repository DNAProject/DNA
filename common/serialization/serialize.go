package serialization

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"sync"
)

var bufPool = sync.Pool{New: func() interface{} { return new([10]byte) }}
var readerPool = sync.Pool{New: func() interface{} { return new(byteReader) }}
var ErrRange = errors.New("value out of range")

//SerializableData describe the data need be serialized.
type SerializableData interface {

	//Write data to writer
	Serialize(w io.Writer)

	//read data to reader
	Deserialize(r io.Reader)
}

func WriteDataList(w io.Writer, list []SerializableData) error {
	len := uint64(len(list))
	WriteVarUint64(w, len)

	for _, data := range list {
		data.Serialize(w)
	}

	return nil
}

func ReadVarInt32(r io.Reader) (int32, int, error) {
	br := getReader(r)
	defer readerPool.Put(br)
	val, err := binary.ReadVarint(br)
	if err != nil {
		return 0, br.n, err
	}
	if val > math.MaxInt32 || val < math.MinInt32 {
		return 0, br.n, ErrRange
	}
	return int32(val), br.n, nil
}

func ReadVarInt64(r io.Reader) (int64, int, error) {
	br := getReader(r)
	defer readerPool.Put(br)
	val, err := binary.ReadVarint(br)
	if err != nil {
		return 0, br.n, err
	}
	if val > math.MaxInt64 || val < math.MinInt64 {
		return 0, br.n, ErrRange
	}
	return int64(val), br.n, nil
}

func ReadVarUint32(r io.Reader) (uint32, int, error) {
	br := getReader(r)
	defer readerPool.Put(br)
	val, err := binary.ReadUvarint(br)
	if err != nil {
		return 0, br.n, err
	}
	if val > math.MaxUint32 {
		return 0, br.n, ErrRange
	}
	return uint32(val), br.n, nil
}

func ReadVarUint64(r io.Reader) (uint64, int, error) {
	br := getReader(r)
	defer readerPool.Put(br)
	val, err := binary.ReadUvarint(br)
	if err != nil {
		return 0, br.n, err
	}
	if val > math.MaxUint64 {
		return 0, br.n, ErrRange
	}
	return val, br.n, nil
}

func ReadVarBytes(r io.Reader) ([]byte, int, error) {
	len, n, err := ReadVarUint32(r)
	if err != nil {
		return nil, n, err
	}
	if len == 0 {
		return nil, n, nil
	}
	buf := make([]byte, len)
	n2, err := io.ReadFull(r, buf)
	return buf, n + n2, err
}

func ReadVarstring(r io.Reader) (string, int, error) {

	x, n, err := ReadVarBytes(r)
	if err != nil {
		return "", n, err
	}
	str := string(x)
	return str, n, err
}

func WriteVarUint32(w io.Writer, val uint32) (int, error) {
	if val > math.MaxUint32 {
		return 0, ErrRange
	}
	valx := uint64(val)
	buf := bufPool.Get().(*[10]byte)
	n := binary.PutUvarint(buf[:], valx)
	b, err := w.Write(buf[:n])
	bufPool.Put(buf)
	return b, err
}

func WriteVarInt32(w io.Writer, val int32) (int, error) {
	if val > math.MaxInt32 || val < math.MinInt32 {
		return 0, ErrRange
	}
	valx := int64(val)
	buf := bufPool.Get().(*[10]byte)
	n := binary.PutVarint(buf[:], valx)
	b, err := w.Write(buf[:n])
	bufPool.Put(buf)
	return b, err
}

func WriteVarInt64(w io.Writer, val int64) (int, error) {
	if val > math.MaxInt64 || val < math.MinInt64 {
		return 0, ErrRange
	}
	buf := bufPool.Get().(*[10]byte)
	n := binary.PutVarint(buf[:], val)
	b, err := w.Write(buf[:n])
	bufPool.Put(buf)
	return b, err
}

func WriteVarUint64(w io.Writer, val uint64) (int, error) {
	if val > math.MaxUint64 {
		return 0, ErrRange
	}
	buf := bufPool.Get().(*[10]byte)
	n := binary.PutUvarint(buf[:], val)
	b, err := w.Write(buf[:n])
	bufPool.Put(buf)
	return b, err
}

func WriteVarBytes(w io.Writer, str []byte) (int, error) {
	n, err := WriteVarUint32(w, uint32(len(str)))
	if err != nil {
		return n, err
	}
	n2, err := w.Write(str)
	return n + n2, err
}

func WriteVarString(w io.Writer, str string) (int, error) {
	buf := []byte(str)
	n, err := WriteVarBytes(w, buf)
	if err != nil {
		return n, err
	}
	return n, err
}

func getReader(r io.Reader) *byteReader {
	br := readerPool.Get().(*byteReader)
	br.reset(r)
	return br
}

type byteReader struct {
	n int
	r io.Reader
	e error
	b [1]byte
}

func (r *byteReader) reset(reader io.Reader) {
	*r = byteReader{n: 0, r: reader, e: nil}
}

func (r *byteReader) ReadByte() (byte, error) {
	if r.e != nil {
		return 0, r.e
	}
	n, err := r.r.Read(r.b[:])
	if n > 0 {
		r.e = err
		r.n++
		return r.b[0], nil
	}
	return 0, err
}
