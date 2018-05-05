package network

type Processor interface {
	// must goroutine safe
	Route(a SrvAgent, msg []interface{}) error
	// must goroutine safe
	Unmarshal(data []byte) ([]interface{}, error)
	// must goroutine safe
	Marshal(msg []interface{}) ([][]byte, error)

	GetLittleEndian() bool
}
