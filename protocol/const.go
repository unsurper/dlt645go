package protocol

const (
	//注册
	RegisterByte = byte(0x40)

	//发送
	SendByte = byte(0x3c)

	//接收
	ReceiveByte = byte(0x3e)

	// IP位
	IPByte = byte(0x49)

	// 标志位
	PrefixID = byte(0x7e)

	// 转义符
	EscapeByte = byte(0x7d)

	// 0x7d < ———— > 0x7d 后紧跟一个0x01
	EscapeByteSufix1 = byte(0x01)

	// 0x7e < ———— > 0x7d 后紧跟一个0x02
	EscapeByteSufix2 = byte(0x02)

	// 消息头大小
	MessageHeaderSize = 12
)
