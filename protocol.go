package tancy

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/deatil/go-crc16/crc16"
	"github.com/funny/link"
	log "github.com/sirupsen/logrus"
	"github.com/unsurper/dlt645go/errors"
	"github.com/unsurper/dlt645go/protocol"
	"io"
)

type Protocol struct {
	privateKey *rsa.PrivateKey
}

// 创建编解码器
func (p Protocol) NewCodec(rw io.ReadWriter) (link.Codec, error) {
	codec := &ProtocolCodec{
		w:               rw,
		r:               rw,
		privateKey:      p.privateKey,
		bufferReceiving: bytes.NewBuffer(nil),
	}
	codec.closer, _ = rw.(io.Closer)
	return codec, nil
}

// 编解码器
type ProtocolCodec struct {
	w               io.Writer
	r               io.Reader
	closer          io.Closer
	publicKey       *rsa.PublicKey
	privateKey      *rsa.PrivateKey
	bufferReceiving *bytes.Buffer
}

// 获取RSA公钥
func (codec *ProtocolCodec) GetPublicKey() *rsa.PublicKey {
	return codec.publicKey
}

// 设置RSA公钥
func (codec *ProtocolCodec) SetPublicKey(publicKey *rsa.PublicKey) {
	codec.publicKey = publicKey
}

// 关闭读写
func (codec *ProtocolCodec) Close() error {
	if codec.closer != nil {
		return codec.closer.Close()
	}
	return nil
}

// 发送消息
func (codec *ProtocolCodec) Send(msg interface{}) error {
	message, ok := msg.(protocol.Message)
	if !ok {
		log.WithFields(log.Fields{
			"reason": errors.ErrInvalidMessage,
		}).Error("[JT/T 808] failed to write message")
		return errors.ErrInvalidMessage
	}

	var err error
	var data []byte
	if codec.publicKey == nil || !message.Header.Property.IsEnableEncrypt() {
		data, err = message.Encode()
	} else {
		data, err = message.Encode(codec.publicKey)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", message.Header.MsgID),
			"reason": err,
		}).Error("[JT/T 808] failed to write message")
		return err
	}

	count, err := codec.w.Write(data)
	if err != nil {
		log.WithFields(log.Fields{
			"id":     fmt.Sprintf("0x%x", message.Header.MsgID),
			"reason": err,
		}).Error("[JT/T 808] failed to write message")
		return err
	}

	log.WithFields(log.Fields{
		"device_id":   message.Header.IccID,
		"msg_type_id": fmt.Sprintf("0x%x", message.Header.MsgID),
		"bytes":       count,
	}).Debug("TX:")
	return nil
}

// 接收消息
func (codec *ProtocolCodec) Receive() (interface{}, error) {
	message, ok, err := codec.readFromBuffer()
	if ok {
		return message, nil
	}
	if err != nil {
		return nil, err
	}

	var buffer [512]byte
	for {
		count, err := io.ReadAtLeast(codec.r, buffer[:], 1)
		if err != nil {
			return nil, err
		}
		codec.bufferReceiving.Write(buffer[:count])

		if codec.bufferReceiving.Len() == 0 {
			continue
		}
		if codec.bufferReceiving.Len() > 0xffff {
			return nil, errors.ErrBodyTooLong
		}

		message, ok, err := codec.readFromBuffer()
		if ok {
			return message, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

// 从缓冲区读取
func (codec *ProtocolCodec) readFromBuffer() (protocol.Message, bool, error) {

	if codec.bufferReceiving.Len() == 0 {
		return protocol.Message{}, false, nil
	}

	dataa := codec.bufferReceiving.Bytes()

	var data []byte
	if dataa[0] == 51 && dataa[1] == 101 {
		// to hex
		var err error
		data, err = hex.DecodeString(string(dataa))
		if err != nil {
			log.WithFields(log.Fields{
				"data":   fmt.Sprintf("%s", dataa),
				"reason": err,
			}).Error("[tancy-flow] failed to hex.DecodeString")
			return protocol.Message{}, false, errors.ErrNotFoundPrefixID
		}
	} else {
		data = dataa
	}
	end := 0

	if data[0] != protocol.RegisterByte && data[0] != protocol.SendByte && data[0] != protocol.ReceiveByte {

		log.WithFields(log.Fields{
			"data":   fmt.Sprintf("%s", data),
			"reason": errors.ErrNotFoundPrefixID,
		}).Debug("[tancy-flow] failed to receive message")
		return protocol.Message{}, false, errors.ErrNotFoundPrefixID
	}

	//CRC16验证
	if data[0] == protocol.SendByte || data[0] == protocol.ReceiveByte {

		var datalen int
		datalen = int(data[1])
		if datalen != len(data) {
			log.WithFields(log.Fields{
				"data":   hex.EncodeToString(data),
				"reason": errors.ErrNotFoundPrefixID,
			}).Error("[tancy-flow] datalength is wrong")
			return protocol.Message{}, false, errors.ErrNotFoundPrefixID
		}
		crc16Hash := crc16.NewCRC16Hash(crc16.CRC16_MODBUS)
		crc16Hash.Write(data[:datalen-2])
		crc16HashData := crc16Hash.Sum(nil)
		crc16HashData2 := hex.EncodeToString(crc16HashData)
		dataHash := hex.EncodeToString(data[datalen-2:])
		data[datalen-2], data[datalen-1] = data[datalen-1], data[datalen-2]
		dataHash2 := hex.EncodeToString(data[datalen-2:])
		//fmt.Println(dataHash, crc16HashData2)
		if dataHash != crc16HashData2 && dataHash2 != crc16HashData2 {
			log.WithFields(log.Fields{
				"data":   hex.EncodeToString(data),
				"reason": errors.ErrCRC16Failed,
			}).Error("[tancy-flow] CRC16 is Wrong")
			return protocol.Message{}, false, errors.ErrCRC16Failed
		}
	}

	var message protocol.Message
	if err := message.Decode(data, codec.privateKey); err != nil {
		log.WithFields(log.Fields{
			"data":   fmt.Sprintf("0x%x", hex.EncodeToString(data)),
			"reason": err,
		}).Error("[tancy-flow] failed to receive message")
		return protocol.Message{}, false, err
	}

	codec.bufferReceiving.Next(end + len(dataa)) //读取长度+len(dataa)

	log.WithFields(log.Fields{
		"device_id":   message.Header.IccID,
		"msg_type_id": fmt.Sprintf("%X", message.Header.MsgID),
	}).Debug("RX:")

	log.WithFields(log.Fields{
		"device_id": message.Header.IccID,
		"hex":       fmt.Sprintf("%0X", data[:]),
		//"Hex": fmt.Sprintf("0x%x", hex.EncodeToString(data[:end+1])),
	}).Trace("RX Raw:")

	return message, true, nil
}
