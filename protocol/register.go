package protocol

// 消息ID枚举
type MsgID uint16

const (
	//  注册心跳
	Msgtancy_0x0040 MsgID = 0x0040
	//  响应间隔或小时记录打包报文
	Msgtancy_0x0008 MsgID = 0x0008
	//  响应间隔或小时记录打包报文
	Msgtancy_0x0009 MsgID = 0x0009
	//上传计量参数报文
	Msgtancy_0x000A MsgID = 0x000A
	//上传登陆信息报文
	Msgtancy_0x000B MsgID = 0x000B
	//上传实时打包报文
	Msgtancy_0x0011 MsgID = 0x0011
	//上传信号强度
	Msgtancy_0x0017 MsgID = 0x0017
)

// 消息实体映射
var entityMapper = map[uint16]func() Entity{
	uint16(Msgtancy_0x0040): func() Entity {
		return new(Tancy_0x0040)
	},
	uint16(Msgtancy_0x0008): func() Entity {
		return new(Tancy_0x0008)
	},
	uint16(Msgtancy_0x0009): func() Entity {
		return new(Tancy_0x0009)
	},
	uint16(Msgtancy_0x000A): func() Entity {
		return new(Tancy_0x000A)
	},
	uint16(Msgtancy_0x000B): func() Entity {
		return new(Tancy_0x000B)
	},
	uint16(Msgtancy_0x0011): func() Entity {
		return new(Tancy_0x0011)
	},
	uint16(Msgtancy_0x0017): func() Entity {
		return new(Tancy_0x0017)
	},
}

// 类型注册
func Register(typ uint16, creator func() Entity) {
	entityMapper[typ] = creator
}
