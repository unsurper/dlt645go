package protocol

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/unsurper/dlt645go/errors"
	"reflect"
	"strconv"
)

// 消息包
type Message struct {
	Header Header
	Body   Entity
}

type DHeader struct {
	MsgID MsgID
}

// 协议编码
func (message *Message) Encode(key ...*rsa.PublicKey) ([]byte, error) {
	// 编码消息体
	var err error
	var body []byte
	if message.Body != nil && !reflect.ValueOf(message.Body).IsNil() {
		body, err = message.Body.Encode()
		if err != nil {
			return nil, err
		}

	}

	// 二进制转义
	buffer := bytes.NewBuffer(nil)

	message.write(buffer, body)
	return buffer.Bytes(), nil
}

// 协议解码
func (message *Message) Decode(data []byte, key ...*rsa.PrivateKey) error {
	// 检验标志位
	if len(data) == 0 {
		return errors.ErrInvalidMessage
	}
	if data[0] != ReceiveByte && data[0] != RegisterByte {
		return errors.ErrInvalidMessage
	}

	var header Header
	var err error
	//处理注册包
	if data[0] == RegisterByte {

		header := Header{}
		json.Unmarshal(data, &header)

		//fmt.Println(header.Imei)
		//fmt.Println()

		header.IccID, err = strconv.ParseUint(header.Imei, 10, 64)
		header.MsgID = 0x0040 //消息ID
		if err == nil {
			message.Body = nil
		} else {
			log.WithFields(log.Fields{
				"id":     fmt.Sprintf("0x%x", header.MsgID),
				"reason": err,
			}).Warn("failed to decode message")
		}
		message.Header = header
		return nil
	} else {
		header.MsgID = MsgID(data[8])                                   //消息ID
		entity, _, err := message.decode(uint16(header.MsgID), data[:]) //解析实体对象 entity     buffer : 为消息标识

		if err == nil {
			message.Body = entity
		} else {
			log.WithFields(log.Fields{
				"id":     fmt.Sprintf("0x%x", header.MsgID),
				"reason": err,
			}).Warn("failed to decode message")
		}
	}

	message.Header = header
	return nil
}

//--->
func (message *Message) decode(typ uint16, data []byte) (Entity, int, error) {
	creator, ok := entityMapper[typ]
	if !ok {
		return nil, 0, errors.ErrTypeNotRegistered
	}

	entity := creator()
	entityPacket, ok := interface{}(entity).(EntityPacket)
	if !ok {
		count, err := entity.Decode(data) //解析data数据
		if err != nil {
			return nil, 0, err
		}
		return entity, count, nil
	}
	err := entityPacket.DecodePacket(data)
	if err != nil {
		return nil, 0, err
	}
	return entityPacket, len(data), nil
}

// 写入二进制数据
func (message *Message) write(buffer *bytes.Buffer, data []byte) *Message {
	for _, b := range data {
		if b == PrefixID {
			buffer.WriteByte(EscapeByte)
			buffer.WriteByte(EscapeByteSufix2)
		} else if b == EscapeByte {
			buffer.WriteByte(EscapeByte)
			buffer.WriteByte(EscapeByteSufix1)
		} else {
			buffer.WriteByte(b)
		}
	}
	return message
}

// 校验和累加计算
func (message *Message) computeChecksum(data []byte, checkSum byte, count int) (byte, int) {
	for _, b := range data {
		checkSum = checkSum ^ b
		if b != PrefixID && b != EscapeByte {
			count++
		} else {
			count += 2
		}
	}
	return checkSum, count
}
