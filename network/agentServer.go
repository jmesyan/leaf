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
	GetTick(callback func([]interface{})) uint16
	ExecTick(tick uint16, data []interface{}) bool
}