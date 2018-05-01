package network

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmesyan/leaf/chanrpc"
	"github.com/jmesyan/leaf/log"
	//"math"
	"reflect"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processer struct {
	littleEndian bool
	pmsgInfo      map[uint32]*MsgInfo
	rpcHandler  	RpcHandler
}

type MsgInfo struct {
	MsgType       reflect.Type
	MsgRouter     *chanrpc.Server
	MsgHandler    MsgHandler
}

type MsgHandler func([]interface{})

func NewProcesser() *Processer {
	p := new(Processer)
	p.littleEndian = false
	p.pmsgInfo = make(map[uint32]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

func (p *Processer) SetRpcHandler(rpcHandler RpcHandler) {
	p.rpcHandler = rpcHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) Register(cmd uint32, msg proto.Message) {
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	if _, ok := p.pmsgInfo[cmd]; ok {
		log.Fatal("message %v is already registered", cmd)
	}

	i := new(MsgInfo)
	i.MsgType = MsgType
	p.pmsgInfo[cmd] = i
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) SetRouter(cmd uint32, MsgRouter *chanrpc.Server) {
	i, ok := p.pmsgInfo[cmd]
	if !ok {
		log.Fatal("message %v not registered", cmd)
	}

	i.MsgRouter = MsgRouter
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) SetHandler(cmd uint32, MsgHandler MsgHandler) {
	i, ok := p.pmsgInfo[cmd]
	if !ok {
		log.Fatal("message %v not registered", cmd)
	}

	i.MsgHandler = MsgHandler
}


// goroutine safe
func (p *Processer) Route(a AgentServer, msg []interface{}) error {
	if p.rpcHandler != nil {
		return p.rpcHandler.Route(a, p, msg)
	}
	panic("bug")
	return nil
}

// goroutine safe
func (p *Processer) Unmarshal(data []byte) ([]interface{}, error) {
	if p.rpcHandler != nil {
		return p.rpcHandler.Unmarshal(p, data)
	}
	panic("bug")

}

// goroutine safe
func (p *Processer) Marshal(msg []interface{}) ([][]byte, error) {
	if p.rpcHandler != nil {
		return p.rpcHandler.Marshal(p, msg)
	}
	return nil, nil
}

func (p *Processer) GetlittleEndian() bool{
	return p.littleEndian
}

func (p *Processer) GetMsgInfo() map[uint32]*MsgInfo{
	return p.pmsgInfo
}



