package protocol

const (
	//注册
	RegisterByte = byte(0x7B)

	//发送
	SendByte = byte(0xFE)

	//接收
	ReceiveByte = byte(0x68)

	// IP位

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
