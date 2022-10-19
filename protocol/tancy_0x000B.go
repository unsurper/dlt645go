package protocol

type Tancy_0x000B struct {
	//间隔记录周期
	SamplesPerCycle byte
	//间隔记录容量
	SamplePerVolume uint16
	//小时记录容量
	SampleHourlyVolume uint16
	//定时或间隔模式
	SetDriverIntervalMessage byte
	//是否长期在线
	Online byte
	//主动或被动上传
	UploadMode byte
	//日次数
	DayTimes byte
	//定时时间 1-10
	Timing string
	//首次时间
	FirstTime uint16
	//上传周期
	UploadCycle uint16
	//电池上传周期
	BatteryUploadCycle uint16
	//小时打包或实时数据
	PackingMode byte
	//是否重复发送
	RepeatSending byte
}

func (entity *Tancy_0x000B) MsgID() MsgID {
	return Msgtancy_0x000B
}

func (entity *Tancy_0x000B) Encode() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (entity *Tancy_0x000B) Decode(data []byte) (int, error) {
	datalen := len(data)
	if datalen < 67 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	var err error
	entity.SamplesPerCycle, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.SamplePerVolume, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	entity.SampleHourlyVolume, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	entity.SetDriverIntervalMessage, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Online, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.UploadMode, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.DayTimes, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.Timing, err = reader.ReadString()
	if err != nil {
		return 0, err
	}
	entity.FirstTime, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	entity.UploadCycle, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	entity.BatteryUploadCycle, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	entity.PackingMode, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	entity.RepeatSending, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}
	return len(data) - reader.Len(), nil
}
