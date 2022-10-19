package protocol

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
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
	if len(data) < 2 || (data[0] != ReceiveByte && data[0] != RegisterByte) {
		return errors.ErrInvalidMessage
	}
	if len(data) == 0 {
		return errors.ErrInvalidMessage
	}

	var header Header
	var err error

	//处理注册包
	if data[0] == RegisterByte {
		i := 2
		for ; i < len(data); i++ {
			if data[i] == IPByte {
				break
			}
		}
		IccID, err := strconv.Atoi(string(data[2:i]))
		if err != nil {
			return err
		}
		header.MsgID = MsgID(data[1]) //消息ID
		header.IccID = uint64(IccID)  //用户名唯一标识码
		log.WithFields(log.Fields{
			"DTU": fmt.Sprintf("user: %s online", data[2:i]),
		}).Info("Register DTU")
		entity, _, err := message.decode(uint16(header.MsgID), data) //解析实体对象 entity     buffer : 为消息标识
		if err == nil {
			message.Body = entity
		} else {
			log.WithFields(log.Fields{
				"id":     fmt.Sprintf("0x%x", header.MsgID),
				"reason": err,
			}).Warn("failed to decode message")
		}
		message.Header = header
		return nil
	}

	//处理响应信号强度报文
	if data[2] == 0x17 {

		header.MsgID = MsgID(data[2]) //消息ID
		header.IccID = uint64(0)      //用户名唯一标识码
		DecID, _ := strconv.Atoi(bcdToString(data[3:5]))
		header.DecID = uint64(DecID)                                     //燃气表唯一标识码
		entity, _, err := message.decode(uint16(header.MsgID), data[3:]) //解析实体对象 entity     buffer : 为消息标识
		if err == nil {
			message.Body = entity
		} else {
			log.WithFields(log.Fields{
				"id":     fmt.Sprintf("0x%x", header.MsgID),
				"reason": err,
			}).Warn("failed to decode message")
		}
		message.Header = header
		return nil
	}

	header.MsgID = MsgID(data[2]) //消息ID

	dec := bcdToString(data[3:11])
	if dec != "" {
		DecID, err := strconv.Atoi(dec)
		if err != nil {
			return err
		}
		header.DecID = uint64(DecID) //燃气表唯一标识码
	} else {
		header.DecID = uint64(0) //燃气表唯一标识码
	}

	header.LocID = hex.EncodeToString(data[11:19]) //远传位置号

	iic := bcdToString(data[19:25])
	if iic != "" {
		IccID, err := strconv.Atoi(iic)
		if err != nil {
			return err
		}
		header.IccID = uint64(IccID) //用户唯一标识码
	}

	header.Uptime, err = fromBCDTime(data[25:31]) //打包上传时间
	if err != nil {
		return err
	}

	entity, _, err := message.decode(uint16(header.MsgID), data[31:]) //解析实体对象 entity     buffer : 为消息标识

	if err == nil {
		message.Body = entity
	} else {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", header.MsgID),
			"reason": err,
		}).Warn("failed to decode message")
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
