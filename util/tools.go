package util
import (
	"bytes"
	"encoding/binary"
)

func BytesToUint32(littleEndian bool, b []byte) uint32 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint32
	if littleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	}

	return tmp
}

func Uint32ToBytes(littleEndian bool,i uint32) []byte {
	var buf = make([]byte, 4)
	if littleEndian{
		binary.LittleEndian.PutUint32(buf, i)
	} else {
		binary.BigEndian.PutUint32(buf, i)
	}
	return buf
}
