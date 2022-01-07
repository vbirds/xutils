package file

import (
	"os"
)

var (
	frameLength int = 1024 * 1024
)

type HwReader struct {
	data   []byte
	length int
	file   *os.File
}

func NewHwReader(filename string) *HwReader {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	return &HwReader{file: f, length: frameLength, data: make([]byte, frameLength)}
}

func (c *HwReader) ReadFrame() (*HwHeader, []byte) {
	h := hwDecode(c.file)
	if h == nil {
		return nil, nil
	}
	if h.Length > frameLength {
		c.data = make([]byte, h.Length)
	}
	datalen, err := c.file.Read(c.data)
	if datalen != int(h.Length) || err != nil {
		return nil, nil
	}
	return h, c.data
}

func (c *HwReader) Close() {
	if c.file != nil {
		c.file.Close()
	}
}
