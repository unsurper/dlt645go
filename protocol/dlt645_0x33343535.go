package protocol

import (
	"encoding/hex"
	"strconv"
)

type Dlt_0x33343535 struct {
	//接收表号
	Devicename []byte
	//表号
	DeviceID string
	//当前A相电流
	Current_a float64
}

func (entity *Dlt_0x33343535) MsgID() MsgID {
	return Msgdlt_0x33343535
}

func (entity *Dlt_0x33343535) Encode() ([]byte, error) {
	writer := NewWriter()

	// 接收符号
	writer.Write([]byte{0xFE, 0xFE, 0x68})
	writer.Write(entity.Devicename)
	writer.Write([]byte{0x68, 0x11, 0x04, 0x33, 0x34, 0x35, 0x35})

	//cs效验位
	var one byte
	for _, v := range writer.Bytes()[2:] {
		one += v
	}
	writer.WriteByte(one)
	// 功能码
	writer.WriteByte(0x16)
	return writer.Bytes(), nil
}

func (entity *Dlt_0x33343535) Decode(data []byte) (int, error) {
	bytea := data[1:7]
	for i, j := 0, len(bytea)-1; i < j; i, j = i+1, j-1 {
		bytea[i], bytea[j] = bytea[j], bytea[i]
	}
	entity.DeviceID = hex.EncodeToString(bytea)
	//正向总电能每个字节-33,1-4,分别为,小数位,个位,百位,万位
	byteb := make([]byte, 3)
	for i := 0; i < 3; i++ {
		byteb[i] = data[14+i] - 0x33
	}
	var err error
	entity.Current_a, err = stringtoCurrent(hex.EncodeToString(byteb))
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func stringtoCurrent(s string) (float64, error) {
	a0, _ := strconv.ParseFloat(s[0:2], 64)
	a1, _ := strconv.ParseFloat(s[2:4], 64)
	a2, _ := strconv.ParseFloat(s[4:6], 64)
	res := a0*0.001 + a1*0.1 + a2*10
	return res, nil
}
