package protocol

// 消息ID枚举
type MsgID uint32

const (
	//  注册心跳
	Msgdlt_0x0040 MsgID = 0x0040
	//上传正向有功总电能
	Msgdlt_0x33333433 MsgID = 0x33333433
	//上传反向有功总电能
	Msgdlt_0x33333533 MsgID = 0x33333533
	//上传A电压
	Msgdlt_0x33343435 MsgID = 0x33343435
	//上传A电流
	Msgdlt_0x33343535 MsgID = 0x33343535
	//上传B电压
	Msgdlt_0x33353435 MsgID = 0x33353435
	//上传B电流
	Msgdlt_0x33353535 MsgID = 0x33353535
	//上传C电压
	Msgdlt_0x33363435 MsgID = 0x33363435
	//上传C电流
	Msgdlt_0x33363535 MsgID = 0x33363535
	//当前总有功功率
	Msgdlt_0x33333635 MsgID = 0x33333635
	//当前总无功功率
	Msgdlt_0x33333735 MsgID = 0x33333735
	//总功率因数
	Msgdlt_0x33333935 MsgID = 0x33333935
)

// 消息实体映射
var entityMapper = map[uint32]func() Entity{
	uint32(Msgdlt_0x0040): func() Entity {
		return new(Dlt_0x0040)
	},
	uint32(Msgdlt_0x33333433): func() Entity {
		return new(Dlt_0x33333433)
	},
	uint32(Msgdlt_0x33333533): func() Entity {
		return new(Dlt_0x33333533)
	},
	uint32(Msgdlt_0x33343435): func() Entity {
		return new(Dlt_0x33343435)
	},
	uint32(Msgdlt_0x33343535): func() Entity {
		return new(Dlt_0x33343535)
	},

	uint32(Msgdlt_0x33353435): func() Entity {
		return new(Dlt_0x33353435)
	},
	uint32(Msgdlt_0x33353535): func() Entity {
		return new(Dlt_0x33353535)
	},
	uint32(Msgdlt_0x33363435): func() Entity {
		return new(Dlt_0x33363435)
	},
	uint32(Msgdlt_0x33363535): func() Entity {
		return new(Dlt_0x33363535)
	},
	uint32(Msgdlt_0x33333635): func() Entity {
		return new(Dlt_0x33333635)
	},
	uint32(Msgdlt_0x33333735): func() Entity {
		return new(Dlt_0x33333735)
	},
	uint32(Msgdlt_0x33333935): func() Entity {
		return new(Dlt_0x33333935)
	},
}

// 类型注册
func Register(typ uint32, creator func() Entity) {
	entityMapper[typ] = creator
}
