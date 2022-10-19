package protocol

import "io"

// 消息实体
type Entity interface {
	MsgID() MsgID
	Encode() ([]byte, error)
	Decode([]byte) (int, error)
}

// 可分包实体
type EntityPacket interface {
	Entity
	GetTag() uint32
	GetReader() io.Reader
	SetReader(io.Reader)
	DecodePacket([]byte) error
}
