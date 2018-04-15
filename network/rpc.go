package network

type RpcHandler interface {
	Marshal(p *Processer, msg interface{}) ([][]byte, error)
	Unmarshal(p *Processer, data []byte) (interface{}, error)
	Route(p *Processer, msg interface{}, userData interface{}) error
}
