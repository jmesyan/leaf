package network
import (
	"net"
)
type AgentServer interface {
	WriteMsg(msg []interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	SetResultData(tick uint32, data []interface{}) error
	GetResultTicker(callback func([]interface{})) uint32
	GetServerProcessor() Processor
}