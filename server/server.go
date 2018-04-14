package server

import (
	"github.com/jmesyan/leaf/log"
	"github.com/jmesyan/leaf/network"
	"net"
	"reflect"
	"time"
)

type Server struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	Onconnected    func(*AgentConn)
	OndisConnected func(*AgentConn) error

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
}

func (server *Server) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if server.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = server.WSAddr
		wsServer.MaxConnNum = server.MaxConnNum
		wsServer.PendingWriteNum = server.PendingWriteNum
		wsServer.MaxMsgLen = server.MaxMsgLen
		wsServer.HTTPTimeout = server.HTTPTimeout
		wsServer.CertFile = server.CertFile
		wsServer.KeyFile = server.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &AgentConn{conn: conn, server: server}
			if server.Onconnected != nil {
				server.Onconnected(a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if server.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = server.TCPAddr
		tcpServer.MaxConnNum = server.MaxConnNum
		tcpServer.PendingWriteNum = server.PendingWriteNum
		tcpServer.LenMsgLen = server.LenMsgLen
		tcpServer.MaxMsgLen = server.MaxMsgLen
		tcpServer.LittleEndian = server.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &AgentConn{conn: conn, server: server}
			if server.Onconnected != nil {
				server.Onconnected(a)
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (server *Server) OnDestroy() {}

type AgentConn struct {
	conn     network.Conn
	server     *Server
	userData interface{}
}

func (a *AgentConn) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.server.Processor != nil {
			msg, err := a.server.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.server.Processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *AgentConn) OnClose() {
	if a.server.OndisConnected != nil{
		err := a.server.OndisConnected(a)
		if err != nil {
			log.Error("rpc error: %v", err)
		}
	}
}

func (a *AgentConn) WriteMsg(msg interface{}) {
	if a.server.Processor != nil {
		data, err := a.server.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *AgentConn) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *AgentConn) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *AgentConn) Close() {
	a.conn.Close()
}

func (a *AgentConn) Destroy() {
	a.conn.Destroy()
}

func (a *AgentConn) UserData() interface{} {
	return a.userData
}

func (a *AgentConn) SetUserData(data interface{}) {
	a.userData = data
}
