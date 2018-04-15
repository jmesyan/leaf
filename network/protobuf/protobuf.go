package protobuf

import (
	//"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jmesyan/leaf/chanrpc"
	"github.com/jmesyan/leaf/log"
	"github.com/jmesyan/leaf/network"
	//"math"
	"reflect"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool
	PmsgInfo      map[string]*MsgInfo
	rpcHandler   network.RpcHandler
}

type MsgInfo struct {
	MsgType       reflect.Type
	msgRouter     *chanrpc.Server
	msgHandler    MsgHandler
	msgRawHandler MsgHandler
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgID      string
	msgRawData []byte
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.littleEndian = false
	p.PmsgInfo = make(map[string]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

func (p *Processor) SetRpcHandler(rpcHandler network.RpcHandler) {
	p.rpcHandler = rpcHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msg proto.Message) string {
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	msgID := MsgType.Elem().Name()
	if msgID == "" {
		log.Fatal("unnamed json message")
	}
	if _, ok := p.PmsgInfo[msgID]; ok {
		log.Fatal("message %v is already registered", msgID)
	}

	i := new(MsgInfo)
	i.MsgType = MsgType
	p.PmsgInfo[msgID] = i
	return msgID
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msg proto.Message, msgRouter *chanrpc.Server) {
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := MsgType.Elem().Name()
	i, ok := p.PmsgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRouter = msgRouter

}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msg proto.Message, msgHandler MsgHandler) {
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := MsgType.Elem().Name()
	i, ok := p.PmsgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgHandler = msgHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRawHandler(msgID string, msgRawHandler MsgHandler) {
	i, ok := p.PmsgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRawHandler = msgRawHandler
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	if p.rpcHandler != nil {
		return p.rpcHandler.Route(p, msg, userData)
	}
	// raw
	if msgRaw, ok := msg.(MsgRaw); ok {
		i, ok := p.PmsgInfo[msgRaw.msgID]
		if !ok {
			return fmt.Errorf("message %v not registered", msgRaw.msgID)
		}
		if i.msgRawHandler != nil {
			i.msgRawHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, userData})
		}
		return nil
	}

	// json
	MsgType := reflect.TypeOf(msg)
	if MsgType == nil || MsgType.Kind() != reflect.Ptr {
		return errors.New("json message pointer required")
	}
	msgID := MsgType.Elem().Name()
	i, ok := p.PmsgInfo[msgID]
	if !ok {
		return fmt.Errorf("message %v not registered", msgID)
	}
	if i.msgHandler != nil {
		i.msgHandler([]interface{}{msg, userData})
	}
	if i.msgRouter != nil {
		i.msgRouter.Go(MsgType, msg, userData)
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
	if p.rpcHandler != nil {
		return p.rpcHandler.Unmarshal(p, data)
	}

	var m map[string][]byte
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	if len(m) != 1 {
		return nil, errors.New("invalid json data")
	}

	for msgID, data := range m {
		i, ok := p.PmsgInfo[msgID]
		if !ok {
			return nil, fmt.Errorf("message %v not registered", msgID)
		}

		// msg
		if i.msgRawHandler != nil {
			return MsgRaw{msgID, data}, nil
		} else {
			msg := reflect.New(i.MsgType.Elem()).Interface()
			return msg, proto.Unmarshal(data, msg.(proto.Message))
		}
	}

	panic("bug")

}

// goroutine safe
func (p *Processor) Marshal(msg interface{}) ([][]byte, error) {
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
