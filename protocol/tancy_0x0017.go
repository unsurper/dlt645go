package protocol

type Tancy_0x0017 struct {
	//标况累积流量
	OpcumFlow uint64
	//工况累积流量
	WocumFlow uint64
	//标况瞬时流量
	WomomFlow uint32
	//工况瞬时流量
	OpmomFlow uint32
	//燃气温度
	TGT uint32
	//燃气压力
	TGP uint32
	//状态码
	State uint32
	//剩余量
	Remain uint64
	// 当前价格
	Nowprice uint32
}

func (entity *Tancy_0x0017) MsgID() MsgID {
	return Msgtancy_0x0017
}

func (entity *Tancy_0x0017) Encode() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (entity *Tancy_0x0017) Decode(data []byte) (int, error) {
	datalen := len(data)
	if datalen < 44 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)
	var err error

	// 标况累积流量
	entity.OpcumFlow, err = reader.ReadUint64()
	if err != nil {
		return 0, err
	}
	// 工况累积流量
	entity.WocumFlow, err = reader.ReadUint64()
	if err != nil {
		return 0, err
	}
	// 标况瞬时流量
	entity.OpmomFlow, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	// 工况瞬时流量
	entity.WomomFlow, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	// 燃气温度
	entity.TGT, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	// 燃气压力
	entity.TGP, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	// 状态字
	entity.State, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	//剩余量
	entity.Remain, err = reader.ReadUint64()
	if err != nil {
		return 0, err
	}

	//当前价格
	entity.Nowprice, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
