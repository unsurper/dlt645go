package protocol

import (
	"github.com/unsurper/dlt645go/errors"
)

// 消息体属性
type Property uint16

// 启用分包
func (property *Property) enablePacket() {
	val := uint16(*property)
	*property = Property(val | (1 << 13))
}

// 启用加密
func (property *Property) enableEncrypt() {
	val := uint16(*property)
	*property = Property(val | (1 << 10))
}

// 是否分包
func (property Property) IsEnablePacket() bool {
	val := uint16(property)
	return val&(1<<13) > 0
}

// 是否加密
func (property Property) IsEnableEncrypt() bool {
	val := uint16(property)
	return val&(1<<10) > 0
}

// 获取消息体长度
func (property *Property) GetBodySize() uint16 {
	// 前十位表示消息体长度
	// 0x3ff == ‭001111111111‬
	val := uint16(*property)
	return ((val << 6) >> 6) & 0x3ff
}

// 设置消息体长度
func (property *Property) SetBodySize(size uint16) error {
	if size > 0x3ff {
		return errors.ErrBodyTooLong
	}
	val := uint16(*property)
	*property = Property(((val >> 10) << 10) | size)
	return nil
}
