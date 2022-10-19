package protocol

import (
	"github.com/unsurper/dlt645go/errors"
	"strconv"
	"time"
)

// 封包信息
type Packet struct {
	Sum uint16
	Seq uint16
}

// 消息头
type Header struct {
	MsgID       MsgID
	DecID       uint64    //燃气表标识
	LocID       string    //远传表位置号
	Uptime      time.Time //上传时间
	Property    Property
	IccID       uint64
	MsgSerialNo uint16
	Packet      *Packet
}

// 协议编码
func (header *Header) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入消息ID
	writer.WriteUint16(uint16(header.MsgID))

	// 写入消息体属性
	if header.Packet != nil {
		header.Property.enablePacket()
	}
	writer.WriteUint16(uint16(header.Property))

	// 写入终端号码
	writer.Write(stringToBCD(strconv.FormatUint(header.IccID, 10), 6))

	// 写入消息流水号
	writer.WriteUint16(header.MsgSerialNo)

	// 写入分包信息
	if header.Property.IsEnablePacket() {
		writer.WriteUint16(header.Packet.Sum)
		writer.WriteUint16(header.Packet.Seq)
	}
	return writer.Bytes(), nil
}

// 协议解码
func (header *Header) Decode(data []byte) error {
	if len(data) < MessageHeaderSize {
		return errors.ErrInvalidHeader
	}
	reader := NewReader(data)

	// 读取消息ID
	msgID, err := reader.ReadUint16()
	if err != nil {
		return errors.ErrInvalidHeader
	}

	// 读取消息体属性
	property, err := reader.ReadUint16()
	if err != nil {
		return errors.ErrInvalidHeader
	}

	// 读取终端号码
	temp, err := reader.Read(6)
	if err != nil {
		return errors.ErrInvalidHeader
	}
	iccID, err := strconv.ParseUint(bcdToString(temp), 10, 64)
	if err != nil {
		return err
	}

	// 读取消息流水号
	serialNo, err := reader.ReadUint16()
	if err != nil {
		return errors.ErrInvalidHeader
	}

	// 读取分包信息
	if Property(property).IsEnablePacket() {
		var packet Packet

		// 读取分包总数
		packet.Sum, err = reader.ReadUint16()
		if err != nil {
			return err
		}

		// 读取分包序列号
		packet.Seq, err = reader.ReadUint16()
		if err != nil {
			return err
		}
		header.Packet = &packet
	}

	header.MsgID = MsgID(msgID)
	header.IccID = iccID
	header.Property = Property(property)
	header.MsgSerialNo = serialNo
	return nil
}
