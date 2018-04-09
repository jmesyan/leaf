package protobuf

import (
	//"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/log"
	//"math"
	"reflect"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool
	//msgInfo      []*MsgInfo
	//msgID        map[reflect.Type]uint16
	//marshalHandler func(msg interface{}) ([][]byte, error)
	//unmarshalHandler func(data []byte) (interface{}, error)
	msgInfo map[string]*MsgInfo

}

type MsgInfo struct {
	msgType       reflect.Type
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
	//p.msgID = make(map[reflect.Type]uint16)
	p.msgInfo = make(map[string]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msg proto.Message) string {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	//fmt.Println(msgType.Elem().Name())
	//if _, ok := p.msgID[msgType]; ok {
	//	log.Fatal("message %s is already registered", msgType)
	//}
	//if len(p.msgInfo) >= math.MaxUint16 {
	//	log.Fatal("too many protobuf messages (max = %v)", math.MaxUint16)
	//}
	//
	//i := new(MsgInfo)
	//i.msgType = msgType
	//p.msgInfo = append(p.msgInfo, i)
	//id := uint16(len(p.msgInfo) - 1)
	//p.msgID[msgType] = id
	//return id
	msgID := msgType.Elem().Name()
	if msgID == "" {
		log.Fatal("unnamed json message")
	}
	if _, ok := p.msgInfo[msgID]; ok {
		log.Fatal("message %v is already registered", msgID)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[msgID] = i
	fmt.Println(8888, p.msgInfo)
	return msgID
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msg proto.Message, msgRouter *chanrpc.Server) {
	//msgType := reflect.TypeOf(msg)
	//if msgType == nil || msgType.Kind() != reflect.Ptr {
	//	log.Fatal("json message pointer required")
	//}
	//msgID := msgType.Elem().Name()
	//i, ok := p.msgInfo[msgID]
	//if !ok {
	//	log.Fatal("message %v not registered", msgID)
	//}
	//
	//i.msgRouter = msgRouter
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRouter = msgRouter

}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msg proto.Message, msgHandler MsgHandler) {
	//msgType := reflect.TypeOf(msg)
	//if msgType == nil || msgType.Kind() != reflect.Ptr {
	//	log.Fatal("json message pointer required")
	//}
	//msgID := msgType.Elem().Name()
	//i, ok := p.msgInfo[msgID]
	//if !ok {
	//	log.Fatal("message %v not registered", msgID)
	//}
	//
	//i.msgHandler = msgHandler
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgHandler = msgHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRawHandler(msgID string, msgRawHandler MsgHandler) {
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRawHandler = msgRawHandler
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	// raw
	if msgRaw, ok := msg.(MsgRaw); ok {
		i, ok := p.msgInfo[msgRaw.msgID]
		if !ok {
			return fmt.Errorf("message %v not registered", msgRaw.msgID)
		}
		if i.msgRawHandler != nil {
			i.msgRawHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, userData})
		}
		return nil
	}

	// json
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return errors.New("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		return fmt.Errorf("message %v not registered", msgID)
	}
	if i.msgHandler != nil {
		i.msgHandler([]interface{}{msg, userData})
	}
	if i.msgRouter != nil {
		i.msgRouter.Go(msgType, msg, userData)
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
	//if len(data) < 2 {
	//	return nil, errors.New("protobuf data too short")
	//}
	//
	//// id
	//var id uint16
	//if p.littleEndian {
	//	id = binary.LittleEndian.Uint16(data)
	//} else {
	//	id = binary.BigEndian.Uint16(data)
	//}
	//if id >= uint16(len(p.msgInfo)) {
	//	return nil, fmt.Errorf("message id %v not registered", id)
	//}
	//
	//// msg
	//i := p.msgInfo[id]
	//if i.msgRawHandler != nil {
	//	return MsgRaw{id, data[2:]}, nil
	//} else {
	//	msg := reflect.New(i.msgType.Elem()).Interface()
	//	return msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
	//}

	var m map[string][]byte
	err := json.Unmarshal(data, &m)
	if err != nil {
		fmt.Println(1111)
		return nil, err
	}
	if len(m) != 1 {
		return nil, errors.New("invalid json data")
	}

	for msgID, data := range m {
		i, ok := p.msgInfo[msgID]
		if !ok {
			return nil, fmt.Errorf("message %v not registered", msgID)
		}

		// msg
		if i.msgRawHandler != nil {
			return MsgRaw{msgID, data}, nil
		} else {
			msg := reflect.New(i.msgType.Elem()).Interface()
			return msg, proto.Unmarshal(data, msg.(proto.Message))
		}
	}

	panic("bug")

}

// goroutine safe
func (p *Processor) Marshal(msg interface{}) ([][]byte, error) {
	//msgType := reflect.TypeOf(msg)
	//
	//// id
	//_id, ok := p.msgID[msgType]
	//if !ok {
	//	err := fmt.Errorf("message %s not registered", msgType)
	//	return nil, err
	//}
	//
	//id := make([]byte, 2)
	//if p.littleEndian {
	//	binary.LittleEndian.PutUint16(id, _id)
	//} else {
	//	binary.BigEndian.PutUint16(id, _id)
	//}
	//
	//// data
	//data, err := proto.Marshal(msg.(proto.Message))
	//return [][]byte{id, data}, err

	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return nil, errors.New("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	if _, ok := p.msgInfo[msgID]; !ok {
		return nil, fmt.Errorf("message %v not registered", msgID)
	}

	pmsg, err := proto.Marshal(msg.(proto.Message))
	if err !=nil{
		return nil, fmt.Errorf("message %v marshal failed", msg)
	}
	// data
	m := map[string][]byte{msgID:pmsg}
	data, err := json.Marshal(m)
	return [][]byte{data}, err
}

// goroutine safe
//func (p *Processor) Range(f func(id uint16, t reflect.Type)) {
//	for id, i := range p.msgInfo {
//		f(uint16(id), i.msgType)
//	}
//
//	//
//	msgType := reflect.TypeOf(msg)
//	if msgType == nil || msgType.Kind() != reflect.Ptr {
//		return nil, errors.New("json message pointer required")
//	}
//	msgID := msgType.Elem().Name()
//	if _, ok := p.msgInfo[msgID]; !ok {
//		return nil, fmt.Errorf("message %v not registered", msgID)
//	}
//
//	// data
//	m := map[string]interface{}{msgID: msg}
//	data, err := json.Marshal(m)
//	return [][]byte{data}, err
//}
