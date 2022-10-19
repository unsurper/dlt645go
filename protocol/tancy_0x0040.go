package protocol

import (
	"github.com/deatil/go-crc16/crc16"
	"time"
)

//  注册应答
type Tancy_0x0040 struct {
	//接收符号
	Receive byte
	//报文总长度
	Totallength byte
	//功能码
	FunctionCode byte
	//时间
	StartTime time.Time
	//条数
	Bars byte
}

func (entity *Tancy_0x0040) MsgID() MsgID {
	return Msgtancy_0x0040
}

func (entity *Tancy_0x0040) Encode() ([]byte, error) {
	writer := NewWriter()

	// 接收符号
	writer.WriteByte(entity.Receive)

	// 报文总长度
	writer.WriteByte(entity.Totallength)

	// 功能码
	writer.WriteByte(entity.FunctionCode)

	//// 时间
	//writer.WriteBcdTime(entity.StartTime)
	//
	//// 条数
	//writer.WriteByte(entity.Bars)

	crc16Hash := crc16.NewCRC16Hash(crc16.CRC16_MODBUS)
	crc16Hash.Write(writer.Bytes())
	crc16HashData := crc16Hash.Sum(nil)

	writer.WriteByte(crc16HashData[0])
	writer.WriteByte(crc16HashData[1])
	return writer.Bytes(), nil
}

func (entity *Tancy_0x0040) Decode(data []byte) (int, error) {
	return 0, nil
}
