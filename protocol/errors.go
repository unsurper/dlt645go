package protocol

import (
	"errors"
)

var (
	// 消息体过长
	ErrBodyTooLong = errors.New("too long message body")
	// 无效消息体
	ErrInvalidBody = errors.New("invalid message body")
	// 无效消息头
	ErrInvalidHeader = errors.New("invalid message header")
	// 无效消息格式
	ErrInvalidMessage = errors.New("invalid message format")
	// 无效消息校验和
	ErrInvalidCheckSum = errors.New("invalid message check sum")
	// 方法尚未实现
	ErrMethodNotImplemented = errors.New("method not implemented")
	// 消息类型未注册
	ErrMessageNotRegistered = errors.New("message not registered")
	// 消息解码错误
	ErrEntityDecode = errors.New("entity decode error")
	// 附加信息长度错误
	ErrInvalidExtraLength = errors.New("invalid extra length")
)
