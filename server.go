package dlt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/unsurper/dlt645go/protocol"
	"net"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/funny/link"
	log "github.com/sirupsen/logrus"
)

// 服务器选项
type Options struct {
	Keepalive    int64
	CloseHandler func(*Session)
	PrivateKey   *rsa.PrivateKey
}

// 协议服务器
type Server struct {
	server          *link.Server
	handler         sessionHandler
	timer           *CountdownTimer
	mutex           sync.Mutex
	sessions        map[uint64]*Session
	privateKey      *rsa.PrivateKey
	closeHandler    func(*Session)
	messageHandlers sync.Map
}

type MessageHandler func(*Session, *protocol.Message)

type Handler interface {
	HandleSession(*Session)
}

var _ Handler = HandlerFunc(nil)

type HandlerFunc func(*Session)

func (f HandlerFunc) HandleSession(session *Session) {
	f(session)
}

// 创建服务
func NewServer(options Options) (*Server, error) {
	if options.Keepalive <= 0 {
		options.Keepalive = 60
	}

	server := Server{
		closeHandler: options.CloseHandler,
		sessions:     make(map[uint64]*Session),
	}
	server.handler.server = &server
	server.timer = NewCountdownTimer(options.Keepalive, server.handleReadTimeout)
	return &server, nil
}

// 运行服务
func (server *Server) Run(network string, port int) error {
	if server.server != nil {
		return errors.New("server already running")
	}

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listen, err := net.Listen(network, address)
	if err != nil {
		return err
	}

	p := Protocol{
		privateKey: server.privateKey,
	}
	server.server = link.NewServer(listen, &p, 24, server.handler)
	log.Infof("[dlt645] protocol server started on %s", address)
	return server.server.Serve()
}

// 停止服务
func (server *Server) Stop() {
	if server.server != nil {
		server.server.Stop()
		server.server = nil
	}
}

// 关闭连接
func (session *Session) Close() error {
	return session.session.Close()
}

// 处理读超时
func (server *Server) handleReadTimeout(key string) {
	sessionID, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return
	}

	session, ok := server.GetSession(sessionID)
	if !ok {
		return
	}
	session.Close()

	log.WithFields(log.Fields{
		"id":        sessionID,
		"device_id": session.iccID,
	}).Debug("session read timeout")
}

// 获取Session
func (server *Server) GetSession(id uint64) (*Session, bool) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	session, ok := server.sessions[id]
	if !ok {
		return nil, false
	}
	return session, true
}

// 处理关闭
func (server *Server) handleClose(session *Session) {
	server.mutex.Lock()
	delete(server.sessions, session.ID())
	server.mutex.Unlock()

	server.timer.Remove(strconv.FormatUint(session.ID(), 10))
	if server.closeHandler != nil {
		func() {
			defer func() {
				if err := recover(); err != nil {
					debug.PrintStack()
				}
			}()
			server.closeHandler(session)
		}()
	}

	log.WithFields(log.Fields{
		"id":        session.ID(),
		"device_id": session.iccID,
	}).Debug("session closed")
}

// 分派消息
func (server *Server) dispatchMessage(session *Session, message *protocol.Message) {
	log.WithFields(log.Fields{
		"id": fmt.Sprintf("0x%x", message.Header.MsgID),
	}).Trace("[dlt645] dispatch message")

	handler, ok := server.messageHandlers.Load(message.Header.MsgID)
	if !ok {
		log.WithFields(log.Fields{
			"id": fmt.Sprintf("0x%x", message.Header.MsgID),
		}).Info("[dlt645] dispatch message canceled, handler not found")
		return
	}

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	server.timer.Update(strconv.FormatUint(session.ID(), 10))
	handler.(MessageHandler)(session, message)
}

// 添加消息处理
func (server *Server) AddHandler(msgID protocol.MsgID, handler MessageHandler) {
	if handler != nil {
		server.messageHandlers.Store(msgID, handler)
	}
}
