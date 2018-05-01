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


func BytesToUint8(littleEndian bool, b []byte) uint8 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint8
	if littleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	}

	return tmp
}

func Uint8ToBytes(littleEndian bool,i uint8) []byte {
	var buf = make([]byte, 1)
	if littleEndian{
		binary.LittleEndian.PutUint8(buf, i)
	} else {
		binary.BigEndian.PutUint8(buf, i)
	}
	return buf
}

func BytesToUint16(littleEndian bool, b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp uint16
	if littleEndian {
		binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	} else {
		binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	}

	return tmp
}

func Uint16ToBytes(littleEndian bool,i uint16) []byte {
	var buf = make([]byte, 2)
	if littleEndian{
		binary.LittleEndian.PutUint16(buf, i)
	} else {
		binary.BigEndian.PutUint16(buf, i)
	}
	return buf
}


