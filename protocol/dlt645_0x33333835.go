package protocol

import (
	"encoding/hex"
	"strconv"
)

type Dlt_0x33333835 struct {
	//接收表号
	Devicename []byte
	//表号
	DeviceID string
	//当前视在功率
	Apparent_power float64
}

func (entity *Dlt_0x33333835) MsgID() MsgID {
	return Msgdlt_0x33333835
}

func (entity *Dlt_0x33333835) Encode() ([]byte, error) {
	writer := NewWriter()

	// 接收符号
	writer.Write([]byte{0xFE, 0xFE, 0x68})
	writer.Write(entity.Devicename)
	writer.Write([]byte{0x68, 0x11, 0x04, 0x33, 0x33, 0x38, 0x35})

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

func (entity *Dlt_0x33333835) Decode(data []byte) (int, error) {
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
	entity.Apparent_power, err = stringtoApparent_power(hex.EncodeToString(byteb))
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func stringtoApparent_power(s string) (float64, error) {
	a0, _ := strconv.ParseFloat(s[0:2], 64)
	a1, _ := strconv.ParseFloat(s[2:4], 64)
	a2, _ := strconv.ParseFloat(s[4:6], 64)
	res := a0*0.0001 + a1*0.01 + a2
	return res, nil
}
