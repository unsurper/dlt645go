package protocol

import "time"

type Tancy_0x0008 struct {
	// 数据
	Dates []Date
	// 当前价格
	Nowprice uint32
	//总帧数
	TotalFrames byte
	//帧序号
	FramesNO byte
}

type Date struct {
	//上传时间
	Uptime time.Time
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
}

func (entity *Tancy_0x0008) MsgID() MsgID {
	return Msgtancy_0x0008
}

func (entity *Tancy_0x0008) Encode() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (entity *Tancy_0x0008) Decode(data []byte) (int, error) {
	datalen := len(data)
	if datalen < 50 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)
	entity.Dates = make([]Date, datalen/50)
	var err error
	for i := 0; i < datalen/50; i++ {
		entity.Dates[i].Uptime, err = reader.ReadBcdTime()
		if err != nil {
			return 0, err
		}
		// 标况累积流量
		entity.Dates[i].OpcumFlow, err = reader.ReadUint64()
		if err != nil {
			return 0, err
		}
		// 工况累积流量
		entity.Dates[i].WocumFlow, err = reader.ReadUint64()
		if err != nil {
			return 0, err
		}
		// 标况瞬时流量
		entity.Dates[i].OpmomFlow, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		// 工况瞬时流量
		entity.Dates[i].WomomFlow, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		// 燃气温度
		entity.Dates[i].TGT, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		// 燃气压力
		entity.Dates[i].TGP, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		// 状态字
		entity.Dates[i].State, err = reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		//剩余量
		entity.Dates[i].Remain, err = reader.ReadUint64()
		if err != nil {
			return 0, err
		}
	}
	//当前价格
	entity.Nowprice, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	//总帧数
	entity.TotalFrames, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	//帧序号
	entity.FramesNO, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
