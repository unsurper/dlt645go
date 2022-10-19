package protocol

import "time"

type Tancy_0x000A struct {
	//修改前仪表参数
	MeterFactorBefore uint32
	//修改前仪表参数
	MeterFactorAfter uint32
	//修改时间
	Modification time.Time
}

func (entity *Tancy_0x000A) MsgID() MsgID {
	return Msgtancy_0x000A
}

func (entity *Tancy_0x000A) Encode() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (entity *Tancy_0x000A) Decode(data []byte) (int, error) {
	datalen := len(data)
	if datalen < 14 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	var err error
	entity.MeterFactorBefore, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	entity.MeterFactorAfter, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	entity.Modification, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
