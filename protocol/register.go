package protocol

// 消息ID枚举
type MsgID uint16

const (
	//  注册心跳
	Msgdlt_0x0040 MsgID = 0x0040
	//上传计量参数报文
	Msgdlt_0x0091 MsgID = 0x0091
)

// 消息实体映射
var entityMapper = map[uint16]func() Entity{
	uint16(Msgdlt_0x0040): func() Entity {
		return new(Dlt_0x0040)
	},
	uint16(Msgdlt_0x0091): func() Entity {
		return new(Dlt_0x0091)
	},
}

// 类型注册
func Register(typ uint16, creator func() Entity) {
	entityMapper[typ] = creator
}
