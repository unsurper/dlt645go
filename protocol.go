package dlt

import (
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
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

	data := codec.bufferReceiving.Bytes()

	end := 0
	if data[0] != protocol.RegisterByte {
		for i := 0; i < len(data); i++ {
			if data[i] == protocol.CRID {
				end = i + 1
				break
			}
		}
	} else {
		end = len(data)
	}

	if data[0] != protocol.RegisterByte && data[0] != protocol.SendByte && data[0] != protocol.ReceiveByte {
		log.WithFields(log.Fields{
			"data":   fmt.Sprintf("%s", data),
			"reason": errors.ErrNotFoundPrefixID,
		}).Debug("[dlt645] failed to receive message for ErrNotFoundPrefixID")
		return protocol.Message{}, false, errors.ErrNotFoundPrefixID
	}

	var message protocol.Message
	if err := message.Decode(data[:end], codec.privateKey); err != nil {
		log.WithFields(log.Fields{
			"data":   fmt.Sprintf("0x%x", hex.EncodeToString(data)),
			"reason": err,
		}).Error("[dlt645] failed to receive message")
		return protocol.Message{}, false, err
	}

	codec.bufferReceiving.Next(end) //读取长度+len(dataa)

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
