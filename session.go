package dlt

import (
	"errors"
	"github.com/funny/link"
	log "github.com/sirupsen/logrus"
	"github.com/unsurper/dlt645go/protocol"

	"runtime/debug"
	"sync"
	"sync/atomic"
)

var SessionClosedError = errors.New("Session Closed")
var SessionBlockedError = errors.New("Session Blocked")

var globalSessionId uint64

// 请求上下文
type requestContext struct {
	msgID    uint16
	serialNo uint16
	callback func(answer *protocol.Message)
}

// 终端会话
type Session struct {
	next    uint32
	iccID   uint64
	server  *Server
	session *link.Session

	mux      sync.Mutex
	requests []requestContext

	UserData interface{}
}

// 创建Session
func newSession(server *Server, sess *link.Session) *Session {
	return &Session{
		server:  server,
		session: sess,
	}
}

// 发送消息
func (session *Session) Send(entity protocol.Entity) (uint16, error) {
	message := protocol.Message{
		Body: entity,
		Header: protocol.Header{
			MsgID: entity.MsgID(),
			IccID: atomic.LoadUint64(&session.iccID),
			//MsgSerialNo: session.nextID(),
		},
	}
	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return message.Header.MsgSerialNo, nil
}

// 获取消息ID
func (session *Session) nextID() uint16 {
	var id uint32
	for {
		id = atomic.LoadUint32(&session.next)
		if id == 0xff {
			if atomic.CompareAndSwapUint32(&session.next, id, 1) {
				id = 1
				break
			}
		} else if atomic.CompareAndSwapUint32(&session.next, id, id+1) {
			id += 1
			break
		}
	}
	return uint16(id)
}

//JT808 平台回复设备消息
//func (session *Session) Reply(msg *protocol.Message, result protocol.Result) (uint16, error) {
//	entity := protocol.T808_0x8001{
//		ReplyMsgSerialNo: msg.Header.MsgSerialNo,
//		ReplyMsgID:       msg.Header.MsgID,
//		Result:           result,
//	}
//	return session.Send(&entity)
//}

// 获取ID
func (session *Session) ID() uint64 {
	return session.session.ID()
}

// 消息接收事件
func (session *Session) message(message *protocol.Message) {
	if message.Header.IccID > 0 {
		old := atomic.LoadUint64(&session.iccID)
		if old != 0 && old != message.Header.IccID {
			log.WithFields(log.Fields{
				"id":  session.ID(),
				"old": old,
				"new": message.Header.IccID,
			}).Warn("[dlt645] terminal IccID is inconsistent")
		}
		atomic.StoreUint64(&session.iccID, message.Header.IccID)
	}

	var msgSerialNo uint16
	switch message.Header.MsgID {
	}
	if msgSerialNo == 0 {
		return
	}

	ctx, ok := session.takeRequestContext(msgSerialNo)
	if ok {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
		ctx.callback(message)
	}
}

// 添加请求上下文
func (session *Session) addRequestContext(ctx requestContext) {
	session.mux.Lock()
	defer session.mux.Unlock()

	for idx, item := range session.requests {
		if item.msgID == ctx.msgID {
			session.requests[idx] = ctx
			return
		}
	}
	session.requests = append(session.requests, ctx)
}

// 取出请求上下文
func (session *Session) takeRequestContext(msgSerialNo uint16) (requestContext, bool) {
	session.mux.Lock()
	defer session.mux.Unlock()

	for idx, item := range session.requests {
		if item.serialNo == msgSerialNo {
			session.requests[idx] = session.requests[len(session.requests)-1]
			session.requests = session.requests[:len(session.requests)-1]
			return item, true
		}
	}
	return requestContext{}, false
}
