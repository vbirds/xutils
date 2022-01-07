package file

import (
	"encoding/binary"
	"io"
)

var (
	hwFlag string = "hwmedia"
)

type HwHeader struct {
	FrameType    uint16
	FrameChannel uint16
	Timestamp    uint64
	Length       int
}

func hwEncode(channel, ctype uint16, timestamp uint64, length int) []byte {
	var data [23]byte
	copy(data[:7], hwFlag)
	binary.LittleEndian.PutUint32(data[7:], uint32(length)+12)
	binary.LittleEndian.PutUint16(data[11:], ctype)
	binary.LittleEndian.PutUint16(data[13:], channel)
	binary.LittleEndian.PutUint64(data[15:], timestamp)
	return data[:23]
}

func hwDecode(r io.Reader) (h *HwHeader) {
	var buf [23]byte
	reclen, err := r.Read(buf[:])
	if reclen != 23 || err != nil {
		return nil
	}
	if string(buf[:]) != hwFlag {
		return nil
	}
	length := binary.LittleEndian.Uint32(buf[7:])
	h.Length = int(length) - 12
	h.FrameType = binary.LittleEndian.Uint16(buf[11:])
	h.FrameChannel = binary.LittleEndian.Uint16(buf[13:])
	h.Timestamp = binary.LittleEndian.Uint64(buf[15:])
	return
}
