package network

import (
	"github.com/jmesyan/leaf/log"
	"net"
	"reflect"
	"time"
	"fmt"
)
type ServerType int32

const (
	SERVER_MASTER ServerType = iota+1
	SERVER_GAME
	SERVER_GATE
)

type Server struct {
	SrvType ServerType
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       Processor
	Onconnected    func(SrvAgent)
	OndisConnected func(SrvAgent) error
	OnMasterConnected func(SrvAgent)

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	//master
	MasterAddr  string
}

func (server *Server) Run(closeSig chan bool) {
	var wsServer *WSServer
	if server.WSAddr != "" {
		wsServer = new(WSServer)
		wsServer.Addr = server.WSAddr
		wsServer.MaxConnNum = server.MaxConnNum
		wsServer.PendingWriteNum = server.PendingWriteNum
		wsServer.MaxMsgLen = server.MaxMsgLen
		wsServer.HTTPTimeout = server.HTTPTimeout
		wsServer.CertFile = server.CertFile
		wsServer.KeyFile = server.KeyFile
		wsServer.NewAgent = func(conn *WSConn) Agent {
			a := &AgentConn{conn: conn, server: server}
			a.ResultMgr = NewAsyncResultMgr()
			if server.Onconnected != nil {
				server.Onconnected(a)
			}
			return a
		}
	}

	var tcpServer *TCPServer
	if server.TCPAddr != "" {
		tcpServer = new(TCPServer)
		tcpServer.Addr = server.TCPAddr
		tcpServer.MaxConnNum = server.MaxConnNum
		tcpServer.PendingWriteNum = server.PendingWriteNum
		tcpServer.LenMsgLen = server.LenMsgLen
		tcpServer.MaxMsgLen = server.MaxMsgLen
		tcpServer.LittleEndian = server.LittleEndian
		tcpServer.NewAgent = func(conn *TCPConn) Agent {
			a := &AgentConn{conn: conn, server: server}
			a.ResultMgr = NewAsyncResultMgr()
			if server.Onconnected != nil {
				server.Onconnected(a)
			}
			return a
		}
	}

	var msClient *TCPClient
	if server.SrvType != SERVER_MASTER && server.MasterAddr != "" {
		msClient = new(TCPClient)
		msClient.AutoReconnect = false
		msClient.Addr = server.MasterAddr
		msClient.ConnNum = 1
		msClient.ConnectInterval = 3 * time.Second
		msClient.PendingWriteNum = server.PendingWriteNum
		msClient.LenMsgLen = server.LenMsgLen
		msClient.MaxMsgLen = server.MaxMsgLen
		msClient.NewAgent = func(conn *TCPConn) Agent{
			a := &AgentConn{conn: conn, server: server}
			a.ResultMgr = NewAsyncResultMgr()
			if server.OnMasterConnected != nil {
				server.OnMasterConnected(a)
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
	if (msClient != nil){
		msClient.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}

	if msClient != nil{
		msClient.Close()
	}
}

func (server *Server) OnDestroy() {}

type AgentConn struct {
	conn     Conn
	server     *Server
	userData interface{}
	ResultMgr  *AsyncResultMgr
}

func (a *AgentConn) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}
		log.Debug("read date=>%v", data)
		if a.server.Processor != nil {
			msg, err := a.server.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.server.Processor.Route(a, msg)
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

func (a *AgentConn) WriteMsg(msg []interface{}) {
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

func (a *AgentConn) SetResultData(tick uint32, data []interface{}) error{
	return a.ResultMgr.FillAsyncResult(tick, data)
}

func (a *AgentConn) GetResultTicker(callback func([]interface{})) uint32{
	if callback == nil{
		return 0
	}
	asyncR, err := a.ResultMgr.Add(false,callback)
	if (err != nil){
		fmt.Println(err)
		return 0;
	}
	return asyncR.GetKey()
}

func (a *AgentConn) GetServer() *Server{
	return a.server
}