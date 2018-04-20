package network

type RpcHandler interface {
	Marshal(p *Processer, msg []interface{}) ([][]byte, error)
	Unmarshal(p *Processer, data []byte) ([]interface{}, error)
	Route(a AgentServer, p *Processer, msg []interface{}) error
}
