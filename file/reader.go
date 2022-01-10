package file

import (
	"bytes"
	"encoding/binary"
	"os"
)

var (
	frameLength int = 1024 * 1024
)

type CacheReader struct {
	data   []byte
	length int
	file   *os.File
}

type HwHeader struct {
	FrameType uint16 //帧类型
	Channel   uint16 //通道编号，从1开始
	Timestamp uint64 //时间戳
}

func NewCacheReader(filename string) *CacheReader {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	return &CacheReader{file: f, length: frameLength, data: make([]byte, frameLength)}
}

func (c *CacheReader) ReadFrame() []byte {
	datalen := cacheDecode(c.file)
	if datalen > frameLength {
		c.data = make([]byte, datalen)
	}
	rlen, err := c.file.Read(c.data)
	if rlen != datalen || err != nil {
		return nil
	}
	return c.data
}

func (c *CacheReader) LastFrame() *HwHeader {
	c.file.Seek(int64(frameLength), os.SEEK_END)
	rlen, err := c.file.Read(c.data)
	if err != nil || rlen < 23 {
		return nil
	}
	pos := bytes.Index(c.data, cacheFlag)
	if pos < 0 {
		return nil
	}
	b := c.data[pos+11:]
	var h HwHeader
	h.FrameType = binary.LittleEndian.Uint16(b[:2])
	h.Channel = binary.LittleEndian.Uint16(b[2:4])
	h.Timestamp = binary.LittleEndian.Uint64(b[4:])
	return &h
}

func (c *CacheReader) Close() {
	if c.file != nil {
		c.file.Close()
	}
}
