package file

import (
	"bytes"
	"encoding/binary"
	"io"
)

var (
	cacheFlag []byte = []byte("hwmedia")
)

func cacheEncode(length int) []byte {
	var data [11]byte
	copy(data[:7], cacheFlag)
	binary.LittleEndian.PutUint32(data[7:], uint32(length))
	return data[:]
}

func cacheDecode(r io.Reader) int {
	var buf [11]byte
	reclen, err := r.Read(buf[:])
	if err != nil || reclen != 11 || !bytes.Equal(buf[:7], cacheFlag) {
		return -1
	}
	length := binary.LittleEndian.Uint32(buf[7:])
	return int(length)
}
