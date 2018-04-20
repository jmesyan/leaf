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
	PmsgInfo      map[uint16]*MsgInfo
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
	p.PmsgInfo = make(map[uint16]*MsgInfo)
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
func (p *Processer) Register(cmd uint16, msg proto.Message) {
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	if _, ok := p.PmsgInfo[cmd]; ok {
		log.Fatal("message %v is already registered", cmd)
	}

	i := new(MsgInfo)
	i.MsgType = MsgType
	p.PmsgInfo[cmd] = i
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) SetRouter(cmd uint16, MsgRouter *chanrpc.Server) {
	i, ok := p.PmsgInfo[cmd]
	if !ok {
		log.Fatal("message %v not registered", cmd)
	}

	i.MsgRouter = MsgRouter
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) SetHandler(cmd uint16, MsgHandler MsgHandler) {
	i, ok := p.PmsgInfo[cmd]
	if !ok {
		log.Fatal("message %v not registered", cmd)
	}

	i.MsgHandler = MsgHandler
}


// goroutine safe
func (p *Processer) Route(agent AgentServer, msg []interface{}) error {
	if p.rpcHandler != nil {
		return p.rpcHandler.Route(agent, p, msg)
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

	//MsgType := reflect.TypeOf(msg)
	//if MsgType == nil || MsgType.Kind() != reflect.Ptr {
	//	return nil, errors.New("json message pointer required")
	//}
	//msgID := MsgType.Elem().Name()
	//if _, ok := p.PmsgInfo[msgID]; !ok {
	//	return nil, fmt.Errorf("message %v not registered", msgID)
	//}
	//
	//pmsg, err := proto.Marshal(msg.(proto.Message))
	//if err != nil {
	//	return nil, fmt.Errorf("message %v marshal failed", msg)
	//}
	//// data
	//m := map[string][]byte{msgID: pmsg}
	//data, err := json.Marshal(m)
	return nil, nil
}

