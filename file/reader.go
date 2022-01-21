package file

import "encoding/binary"

type Header struct {
	FrameType uint16 //帧类型
	Channel   uint16 //通道编号，从1开始
	Timestamp uint64 //时间戳
}

func (h *Header) decodec(data []byte) {
	h.FrameType = binary.LittleEndian.Uint16(data[:2])
	h.Channel = binary.LittleEndian.Uint16(data[2:4])
	h.Timestamp = binary.LittleEndian.Uint64(data[4:])
}

type Reader interface {
	ReadFrame(int) (*Header, []byte, int)
	LastFrame() *Header
	Close()
}
