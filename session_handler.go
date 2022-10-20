package dlt

import (
	"github.com/funny/link"
	log "github.com/sirupsen/logrus"
	"github.com/unsurper/dlt645go/protocol"
	"strconv"
)

// Session处理
type sessionHandler struct {
	server          *Server
	autoMergePacket bool
}

func (handler sessionHandler) HandleSession(sess *link.Session) {
	log.WithFields(log.Fields{
		"id": sess.ID(),
	}).Debug("[tancy-flow] new session created")

	// 创建Session
	session := newSession(handler.server, sess)
	handler.server.mutex.Lock()
	handler.server.sessions[sess.ID()] = session
	handler.server.mutex.Unlock()
	handler.server.timer.Update(strconv.FormatUint(session.ID(), 10))
	sess.AddCloseCallback(nil, nil, func() {
		handler.server.handleClose(session)
	})

	for {
		// 接收消息
		msg, err := sess.Receive()
		if err != nil {
			sess.Close()
			break
		}
		// 分发消息
		message := msg.(protocol.Message)
		if message.Header.MsgID == protocol.MsgID(protocol.RegisterByte) {
			session.iccID = message.Header.IccID
		} else if message.Header.IccID == 0 {
			message.Header.IccID = session.iccID
		}
		session.message(&message)
		handler.server.dispatchMessage(session, &message)
	}
}
