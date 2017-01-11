package serialization

import (
	"bytes"
	"io/ioutil"
	"math"
	"strings"
	"testing"
)

func BenchmarkReadVarUint32(b *testing.B) {
	data := []byte{0xff, 0xff, 0xff, 0xff, 0x01}
	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		ReadVarUint32(r)
	}
}

func BenchmarkReadVarInt32(b *testing.B) {
	data := []byte{0xff, 0xff, 0xff, 0xff, 0x01}
	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		ReadVarUint32(r)
	}
}

func BenchmarkReadVarUint64(b *testing.B) {
	data := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}
	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		//r.Reset(data)
		ReadVarUint64(r)
	}
}

func BenchmarkReadVarInt64(b *testing.B) {
	data := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}
	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		//r.Reset(data)
		ReadVarUint64(r)
	}
}

func BenchmarkWriteVarUint32(b *testing.B) {
	n := uint32(math.MaxUint32)
	for i := 0; i < b.N; i++ {
		WriteVarUint32(ioutil.Discard, n)
	}
}

func BenchmarkWriteVarInt32(b *testing.B) {
	n := int32(math.MaxInt32)
	for i := 0; i < b.N; i++ {
		WriteVarInt32(ioutil.Discard, n)
	}
}

func BenchmarkWriteVarUint64(b *testing.B) {
	n := uint64(math.MaxUint64)
	for i := 0; i < b.N; i++ {
		WriteVarUint64(ioutil.Discard, n)
	}
}

func BenchmarkWriteVarInt64(b *testing.B) {
	n := int64(math.MaxInt64)
	for i := 0; i < b.N; i++ {
		WriteVarInt64(ioutil.Discard, n)
	}
}

func BenchmarkWriteVarBytes(b *testing.B) {
	s := []byte{10, 11, 12}
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		WriteVarBytes(buf, s)
	}
}

func BenchmarkWriteVarString(b *testing.B) {
	s := "jim"
	buf := new(bytes.Buffer)
	for i := 0; i < b.N; i++ {
		WriteVarString(buf, s)
	}
}

func TestVarUint32(t *testing.T) {
	cases := []struct {
		n       uint32
		want    []byte
		wantErr error
	}{
		{
			n:    0,
			want: []byte{0},
		},
		{
			n:    500,
			want: []byte{0xf4, 0x03},
		},
		{
			n:    math.MaxUint32,
			want: []byte{0xff, 0xff, 0xff, 0xff, 0x0f},
		},
	}

	for _, c := range cases {
		b := new(bytes.Buffer)
		n, err := WriteVarUint32(b, c.n)
		if c.wantErr != err {
			t.Errorf("WriteVarUint32(%d): err %v, want %v", c.n, err, c.wantErr)
			continue
		}
		if c.wantErr != nil {
			continue
		}
		if n != len(c.want) {
			t.Errorf("WriteVarUint32(%d): wrote %d byte(s), want %d", c.n, n, len(c.want))
		}
		if !bytes.Equal(c.want, b.Bytes()) {
			t.Errorf("WriteVarUint32(%d): got %x, want %x", c.n, b.Bytes(), c.want)
		}
		b = bytes.NewBuffer(b.Bytes())
		v, n, err := ReadVarUint32(b)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(c.want) {
			t.Errorf("ReadVarUint32 [c.n = %d] got %d bytes, want %d", c.n, n, len(c.want))
		}
		if uint32(v) != c.n {
			t.Errorf("ReadVarUint32 got %d, want %d", v, c.n)
		}
	}
}

func TestVarInt32(t *testing.T) {
	cases := []struct {
		n       int32
		want    []byte
		wantErr error
	}{
		{
			n:    0,
			want: []byte{0},
		},
		{
			n:    500,
			want: []byte{0xe8, 0x07},
		},
		{
			n:    math.MaxInt32,
			want: []byte{0xfe, 0xff, 0xff, 0xff, 0x0f},
		},
	}

	for _, c := range cases {
		b := new(bytes.Buffer)
		n, err := WriteVarInt32(b, c.n)
		if c.wantErr != err {
			t.Errorf("WriteVarInt32(%d): err %v, want %v", c.n, err, c.wantErr)
			continue
		}
		if c.wantErr != nil {
			continue
		}
		if n != len(c.want) {
			t.Errorf("WriteVarIint32(%d): wrote %d byte(s), want %d", c.n, n, len(c.want))
		}
		if !bytes.Equal(c.want, b.Bytes()) {
			t.Errorf("WriteVarIint32(%d): got %x, want %x", c.n, b.Bytes(), c.want)
		}
		b = bytes.NewBuffer(b.Bytes())
		v, n, err := ReadVarInt32(b)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(c.want) {
			t.Errorf("ReadVarInt32 [c.n = %d] got %d bytes, want %d", c.n, n, len(c.want))
		}
		if int32(v) != c.n {
			t.Errorf("ReadVarInt32 got %d, want %d", v, c.n)
		}
	}
}

func TestVarUint64(t *testing.T) {
	cases := []struct {
		n       uint64
		want    []byte
		wantErr error
	}{
		{
			n:    0,
			want: []byte{0},
		},
		{
			n:    500,
			want: []byte{0xf4, 0x03},
		},
		{
			n:    math.MaxUint64,
			want: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01},
		},
	}

	for _, c := range cases {
		b := new(bytes.Buffer)
		n, err := WriteVarUint64(b, c.n)
		if c.wantErr != err {
			t.Errorf("WriteVarUint64(%d): err %v, want %v", c.n, err, c.wantErr)
			continue
		}
		if c.wantErr != nil {
			continue
		}
		if n != len(c.want) {
			t.Errorf("WriteVarUint64(%d): wrote %d byte(s), want %d", c.n, n, len(c.want))
		}
		if !bytes.Equal(c.want, b.Bytes()) {
			t.Errorf("WriteVarUint64(%d): got %x, want %x", c.n, b.Bytes(), c.want)
		}
		b = bytes.NewBuffer(b.Bytes())
		v, n, err := ReadVarUint64(b)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(c.want) {
			t.Errorf("ReadVarUint64 [c.n = %d] got %d bytes, want %d", c.n, n, len(c.want))
		}
		if uint64(v) != c.n {
			t.Errorf("ReadVarUint64 got %d, want %d", v, c.n)
		}
	}
}

func TestVarInt64(t *testing.T) {
	cases := []struct {
		n       int64
		want    []byte
		wantErr error
	}{
		{
			n:    0,
			want: []byte{0},
		},
		{
			n:    500,
			want: []byte{0xe8, 0x07},
		},
		{
			n:    math.MaxInt64,
			want: []byte{0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01},
		},
	}

	for _, c := range cases {
		b := new(bytes.Buffer)
		n, err := WriteVarInt64(b, c.n)
		if c.wantErr != err {
			t.Errorf("WriteVarInt64(%d): err %v, want %v", c.n, err, c.wantErr)
			continue
		}
		if c.wantErr != nil {
			continue
		}
		if n != len(c.want) {
			t.Errorf("WriteVarInt64(%d): wrote %d byte(s), want %d", c.n, n, len(c.want))
		}
		if !bytes.Equal(c.want, b.Bytes()) {
			t.Errorf("WriteVarInt64(%d): got %x, want %x", c.n, b.Bytes(), c.want)
		}
		b = bytes.NewBuffer(b.Bytes())
		v, n, err := ReadVarInt64(b)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(c.want) {
			t.Errorf("ReadVarInt64 [c.n = %d] got %d bytes, want %d", c.n, n, len(c.want))
		}
		if int64(v) != c.n {
			t.Errorf("ReadVarInt64 got %d, want %d", v, c.n)
		}
	}
}

func TestVarBytes(t *testing.T) {
	s := []byte{10, 11, 12}
	b := new(bytes.Buffer)
	_, err := WriteVarBytes(b, s)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{3, 10, 11, 12}
	if !bytes.Equal(b.Bytes(), want) {
		t.Errorf("got %x, want %x", b.Bytes(), want)
	}
	b = bytes.NewBuffer(want)
	s, _, err = ReadVarBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	want = []byte{10, 11, 12}
	if !bytes.Equal(s, want) {
		t.Errorf("got %x, expected %x", s, want)
	}
}

func TestVarString(t *testing.T) {
	s := "hello"
	b := new(bytes.Buffer)
	_, err := WriteVarString(b, s)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{5, 104, 101, 108, 108, 111}
	if !bytes.Equal(b.Bytes(), want) {
		t.Errorf("got %x, want %x", b.Bytes(), want)
	}
	b = bytes.NewBuffer(want)
	s, _, err = ReadVarstring(b)
	if err != nil {
		t.Fatal(err)
	}
	want2 := "hello"
	if !strings.EqualFold(s, want2) {
		t.Errorf("got %x, expected %x", s, want2)
	}

}

func TestEmptyVarBytes(t *testing.T) {
	s := []byte{}
	b := new(bytes.Buffer)
	_, err := WriteVarBytes(b, s)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x00}
	if !bytes.Equal(b.Bytes(), want) {
		t.Errorf("got %x, want %x", b.Bytes(), want)
	}

	b = bytes.NewBuffer(want)
	s, _, err = ReadVarBytes(b)
	if err != nil {
		t.Fatal(err)
	}
	want = nil // we deliberately return nil for empty strings to avoid unnecessary byteslice allocation
	if !bytes.Equal(s, want) {
		t.Errorf("got %x, expected %x", s, want)
	}
}
