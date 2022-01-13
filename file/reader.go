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

type CacheHeader struct {
	FrameType uint16 //帧类型
	Channel   uint16 //通道编号，从1开始
	Timestamp uint64 //时间戳
}

func (h *CacheHeader) decodec(data []byte) {
	h.FrameType = binary.LittleEndian.Uint16(data[:2])
	h.Channel = binary.LittleEndian.Uint16(data[2:4])
	h.Timestamp = binary.LittleEndian.Uint64(data[4:])
}

func NewCacheReader(filename string) *CacheReader {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	return &CacheReader{file: f, length: frameLength, data: make([]byte, frameLength)}
}

func (c *CacheReader) ReadFrame(offset int) (*CacheHeader, []byte, int) {
	datalen := cacheDecode(c.file)
	if datalen < 0 {
		return nil, nil, 0
	}
	pkglen := datalen + offset
	if pkglen > c.length {
		c.data = make([]byte, pkglen)
		c.length = pkglen
	}
	rlen, err := c.file.Read(c.data[offset:pkglen])
	if rlen != datalen || err != nil {
		return nil, nil, 0
	}
	var h CacheHeader
	h.decodec(c.data[offset:])
	return &h, c.data[:pkglen], datalen
}

func (c *CacheReader) LastFrame() *CacheHeader {
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
	var h CacheHeader
	h.decodec(b)
	return &h
}

func (c *CacheReader) Close() {
	if c.file != nil {
		c.file.Close()
	}
}
