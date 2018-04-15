package network

import (

	"github.com/golang/protobuf/proto"
	"github.com/jmesyan/leaf/chanrpc"
	"github.com/jmesyan/leaf/log"
	//"math"
	"reflect"
	"errors"
	"fmt"
	"encoding/json"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processer struct {
	littleEndian bool
	PmsgInfo      map[string]*MsgInfo
	rpcHandler  	RpcHandler
}

type MsgInfo struct {
	MsgType       reflect.Type
	MsgRouter     *chanrpc.Server
	MsgHandler    MsgHandler
	MsgRawHandler MsgHandler
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	MsgID      string
	MsgRawData []byte
}

func NewProcesser() *Processer {
	p := new(Processer)
	p.littleEndian = false
	p.PmsgInfo = make(map[string]*MsgInfo)
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
func (p *Processer) Register(msg proto.Message) string {
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	MsgID := MsgType.Elem().Name()
	if MsgID == "" {
		log.Fatal("unnamed json message")
	}
	if _, ok := p.PmsgInfo[MsgID]; ok {
		log.Fatal("message %v is already registered", MsgID)
	}

	i := new(MsgInfo)
	i.MsgType = MsgType
	p.PmsgInfo[MsgID] = i
	return MsgID
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) SetRouter(msg proto.Message, MsgRouter *chanrpc.Server) {
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	MsgID := MsgType.Elem().Name()
	i, ok := p.PmsgInfo[MsgID]
	if !ok {
		log.Fatal("message %v not registered", MsgID)
	}

	i.MsgRouter = MsgRouter

}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) SetHandler(msg proto.Message, MsgHandler MsgHandler) {
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	MsgID := MsgType.Elem().Name()
	i, ok := p.PmsgInfo[MsgID]
	if !ok {
		log.Fatal("message %v not registered", MsgID)
	}

	i.MsgHandler = MsgHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processer) SetRawHandler(MsgID string, MsgRawHandler MsgHandler) {
	i, ok := p.PmsgInfo[MsgID]
	if !ok {
		log.Fatal("message %v not registered", MsgID)
	}

	i.MsgRawHandler = MsgRawHandler
}

// goroutine safe
func (p *Processer) Route(msg interface{}, userData interface{}) error {
	if p.rpcHandler != nil {
		return p.rpcHandler.Route(p, msg, userData)
	}
	panic("bug")
	return nil
}

// goroutine safe
func (p *Processer) Unmarshal(data []byte) (interface{}, error) {
	if p.rpcHandler != nil {
		return p.rpcHandler.Unmarshal(p, data)
	}
	panic("bug")

}

// goroutine safe
func (p *Processer) Marshal(msg interface{}) ([][]byte, error) {
	if p.rpcHandler != nil {
		return p.rpcHandler.Marshal(p, msg)
	}
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		return nil, errors.New("json message pointer required")
	}
	msgID := MsgType.Elem().Name()
	if _, ok := p.PmsgInfo[msgID]; !ok {
		return nil, fmt.Errorf("message %v not registered", msgID)
	}

	pmsg, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		return nil, fmt.Errorf("message %v marshal failed", msg)
	}
	// data
	m := map[string][]byte{msgID: pmsg}
	data, err := json.Marshal(m)
	return [][]byte{data}, err
}
