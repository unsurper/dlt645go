package errors

import (
	"errors"
)

var (
	ErrCRC16Failed = errors.New("CRC16 verification failed")

	// 无效消息体
	ErrInvalidBody = errors.New("invalid body")
	// 消息体过长
	ErrBodyTooLong = errors.New("body too long")
	// 无效消息头
	ErrInvalidHeader = errors.New("invalid header")
	// 未找到标识符
	ErrNotFoundPrefixID = errors.New("not found prefix")
	// 无效BCD时间
	ErrInvalidBCDTime = errors.New("invalid BCD time")
	// 无效消息格式
	ErrInvalidMessage = errors.New("invalid message")
	// 无效消息校验和
	ErrInvalidCheckSum = errors.New("invalid check sum")
	// 消息类型未注册
	ErrTypeNotRegistered = errors.New("entity not registered")
	// 附加信息长度错误
	ErrInvalidExtraLength = errors.New("invalid extra length")
	// 消息解密失败
	ErrDecryptMessageFailed = errors.New("decrypt message failed")
	//
	WTF = errors.New("wtf")
)
