package protocol

import (
	"bytes"
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

// 终端参数
type Param struct {
	id         uint32
	serialized []byte
}

// 参数ID
func (param *Param) ID() uint32 {
	return param.id
}

// 设为Byte
func (param *Param) SetByte(id uint32, b byte) *Param {
	param.id = id
	param.serialized = []byte{b}
	return param
}

// 设为Bytes
func (param *Param) SetBytes(id uint32, b []byte) *Param {
	param.id = id
	buffer := make([]byte, len(b))
	copy(buffer, b)
	param.serialized = buffer
	return param
}

// 设为Uint16
func (param *Param) SetUint16(id uint32, n uint16) *Param {
	param.id = id
	var buffer [2]byte
	binary.BigEndian.PutUint16(buffer[:], n)
	param.serialized = buffer[:]
	return param
}

// 设为Uint32
func (param *Param) SetUint32(id uint32, n uint32) *Param {
	param.id = id
	var buffer [4]byte
	binary.BigEndian.PutUint32(buffer[:], n)
	param.serialized = buffer[:]
	return param
}

// 设为字符串
func (param *Param) SetString(id uint32, s string) *Param {
	if len(s) == 0 {
		return param.SetBytes(id, nil)
	}
	data, _ := ioutil.ReadAll(transform.NewReader(
		bytes.NewReader([]byte(s)), simplifiedchinese.GB18030.NewEncoder()))
	return param.SetBytes(id, data)
}

// 读取Byte
func (param *Param) GetByte() (byte, error) {
	if len(param.serialized) < 1 {
		return 0, ErrInvalidBody
	}
	return param.serialized[0], nil
}

// 读取Bytes
func (param *Param) GetBytes() ([]byte, error) {
	data := make([]byte, len(param.serialized))
	copy(data, param.serialized)
	return data, nil
}

// 读取Uint16
func (param *Param) GetUint16() (uint16, error) {
	if len(param.serialized) < 2 {
		return 0, ErrInvalidBody
	}
	return binary.BigEndian.Uint16(param.serialized[:2]), nil
}

// 读取Uint32
func (param *Param) GetUint32() (uint32, error) {
	if len(param.serialized) < 4 {
		return 0, ErrInvalidBody
	}
	return binary.BigEndian.Uint32(param.serialized[:4]), nil
}

// 读取字符串
func (param *Param) GetString() (string, error) {
	data, err := ioutil.ReadAll(transform.NewReader(
		bytes.NewReader(param.serialized), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(data), nil
}
