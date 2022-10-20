package protocol

//  注册应答
type Dlt_0x0040 struct {
}

func (entity *Dlt_0x0040) MsgID() MsgID {
	return Msgdlt_0x0040
}

func (entity *Dlt_0x0040) Encode() ([]byte, error) {
	panic("emmm")
}

func (entity *Dlt_0x0040) Decode(data []byte) (int, error) {
	return 0, nil
}
