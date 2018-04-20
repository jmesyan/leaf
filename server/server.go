package server

import (
	"github.com/jmesyan/leaf/log"
	"github.com/jmesyan/leaf/network"
	"net"
	"reflect"
	"time"
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
	Processor       network.Processor
	Onconnected    func(*AgentConn)
	OndisConnected func(*AgentConn) error
	OnMasterConnected func(*AgentConn)

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
			a := &AgentConn{conn: conn, server: server, ticker:0}
			a.callbacks = make(map[uint16]func([]interface{}))
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
			a := &AgentConn{conn: conn, server: server, ticker:0}
			a.callbacks = make(map[uint16]func([]interface{}))
			if server.Onconnected != nil {
				server.Onconnected(a)
			}
			return a
		}
	}

	var msClient *network.TCPClient
	if server.SrvType != SERVER_MASTER && server.MasterAddr != "" {
		msClient = new(network.TCPClient)
		msClient.AutoReconnect = false
		msClient.Addr = server.MasterAddr
		msClient.ConnNum = 1
		msClient.ConnectInterval = 3 * time.Second
		msClient.PendingWriteNum = server.PendingWriteNum
		msClient.LenMsgLen = server.LenMsgLen
		msClient.MaxMsgLen = server.MaxMsgLen
		msClient.NewAgent = func(conn *network.TCPConn) network.Agent{
			a := &AgentConn{conn: conn, server: server, ticker:0}
			a.callbacks = make(map[uint16]func([]interface{}))
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
	conn     network.Conn
	server     *Server
	userData interface{}
	ticker    uint16
	callbacks map[uint16]func([]interface{})
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

func (a *AgentConn) GetTick(callback func([]interface{})) uint16{
	if (callback != nil) {
		if (a.ticker > 65535){
			a.ticker = 0
		};
		a.ticker++;
		a.callbacks[a.ticker] = callback;
		return a.ticker;
	}
	return 0;
}

func (a *AgentConn) ExecTick(tick uint16, data []interface{}) bool{
	if callback, ok := a.callbacks[tick];ok{
		callback(data)
		delete(a.callbacks, tick)
		return true;
	}
	return false;
}
