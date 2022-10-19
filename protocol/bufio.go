package protocol

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"time"
)

type Writer struct {
	b *bytes.Buffer
}

func NewWriter() Writer {
	return Writer{b: bytes.NewBuffer(nil)}
}

func (writer *Writer) Bytes() []byte {
	return writer.b.Bytes()
}

func (writer *Writer) Write(p []byte, size ...int) *Writer {
	if len(size) == 0 {
		writer.b.Write(p)
		return writer
	}

	if len(p) >= size[0] {
		writer.b.Write(p[:size[0]])
	} else {
		writer.b.Write(p)
		end := size[0] - len(p)
		for i := 0; i < end; i++ {
			writer.b.WriteByte(0)
		}
	}
	return writer
}

func (writer *Writer) WriteByte(b byte) *Writer {
	writer.b.WriteByte(b)
	return writer
}

func (writer *Writer) WriteUint16(n uint16) *Writer {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], n)
	writer.b.Write(buf[:])
	return writer
}

func (writer *Writer) WriteUint32(n uint32) *Writer {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], n)
	writer.b.Write(buf[:])
	return writer
}

func (writer *Writer) WriteBcdTime(t time.Time) *Writer {
	writer.b.Write(toBCDTime(t))
	return writer
}

func (writer *Writer) WritString(str string, size ...int) error {
	reader := bytes.NewReader([]byte(str))
	data, err := ioutil.ReadAll(
		transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder()))
	if err != nil {
		return err
	}
	writer.Write(data, size...)
	return nil
}

type Reader struct {
	d []byte
	r *bytes.Reader
}

func NewReader(data []byte) Reader {
	return Reader{d: data, r: bytes.NewReader(data)}
}

func (reader *Reader) Len() int {
	return reader.r.Len()
}

func (reader *Reader) Read(size ...int) ([]byte, error) {
	num := reader.r.Len()
	if len(size) > 0 {
		num = size[0]
	}

	if num > reader.r.Len() {
		return nil, io.ErrUnexpectedEOF
	}

	curr := len(reader.d) - reader.r.Len()
	buf := reader.d[curr : curr+num]
	reader.r.Seek(int64(num), io.SeekCurrent)
	return buf, nil
}

func (reader *Reader) ReadByte() (byte, error) {
	return reader.r.ReadByte()
}

func (reader *Reader) ReadUint16() (uint16, error) {
	if reader.r.Len() < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	var buf [2]byte
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != len(buf) {
		return 0, io.ErrUnexpectedEOF
	}
	return binary.BigEndian.Uint16(buf[:]), nil
}

func (reader *Reader) ReadUint32() (uint32, error) {
	if reader.r.Len() < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	var buf [4]byte
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != len(buf) {
		return 0, io.ErrUnexpectedEOF
	}
	return binary.BigEndian.Uint32(buf[:]), nil
}

//LittleEndian高位到低位读取
func (reader *Reader) ReadUint64() (uint64, error) {
	if reader.r.Len() < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	var buf [8]byte
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != len(buf) {
		return 0, io.ErrUnexpectedEOF
	}
	return binary.BigEndian.Uint64(buf[:]), nil
}

func (reader *Reader) ReadBcdTime() (time.Time, error) {
	if reader.r.Len() < 6 {
		return time.Time{}, io.ErrUnexpectedEOF
	}

	var buf [6]byte
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return time.Time{}, err
	}
	if n != len(buf) {
		return time.Time{}, io.ErrUnexpectedEOF
	}
	return fromBCDTime(buf[:])
}

func (reader *Reader) ReadString(size ...int) (string, error) {
	data, err := reader.Read(size...)
	if err != nil {
		return "", err
	}

	text, err := ioutil.ReadAll(transform.NewReader(
		bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return "", err
	}
	return BytesToString(text), nil

}
