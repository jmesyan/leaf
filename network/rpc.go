package network

type RpcHandler interface {
	Marshal(p Processor, msg []interface{}) ([][]byte, error)
	Unmarshal(p Processor, data []byte) ([]interface{}, error)
	Route(a AgentServer, p Processor, msg []interface{}) error
}
