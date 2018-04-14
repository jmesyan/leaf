package network

type RpcHandler interface {
	Marshal(msg interface{}) ([][]byte, error)
	Unmarshal(data []byte) (interface{}, error)
	Route(msg interface{}, userData interface{}) error
}
