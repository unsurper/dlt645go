package protocol

//  注册应答
type Dlt_0x0040 struct {
	//接收符号
	Devicename []byte
}

func (entity *Dlt_0x0040) MsgID() MsgID {
	return Msgdlt_0x0040
}

func (entity *Dlt_0x0040) Encode() ([]byte, error) {
	writer := NewWriter()

	// 接收符号
	writer.Write([]byte{0xFE, 0xFE, 0x68})
	writer.Write(entity.Devicename)
	writer.Write([]byte{0x68, 0x11, 0x04, 0x33, 0x33, 0x34, 0x33})

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

func (entity *Dlt_0x0040) Decode(data []byte) (int, error) {
	return 0, nil
}
