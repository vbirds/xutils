package file

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	frameLength int = 1024 * 1024
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

type CacheReader struct {
	data   []byte
	length int
	file   *os.File
}

func NewCacheReader(filename string) Reader {
	f, err := os.Open(filename)
	if err != nil {
		return nil
	}
	return &CacheReader{file: f, length: frameLength, data: make([]byte, frameLength)}
}

func (c *CacheReader) ReadFrame(offset int) (*Header, []byte, int) {
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
	var h Header
	h.decodec(c.data[offset:])
	return &h, c.data[:pkglen], datalen
}

func (c *CacheReader) LastFrame() *Header {
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
	var h Header
	h.decodec(b)
	return &h
}

func (c *CacheReader) Close() {
	if c.file != nil {
		c.file.Close()
	}
}

func unix2Str(timestamp int64, layout string, offset int) string {
	var cstzone = time.FixedZone("CST", offset)
	timeStr := time.Unix(timestamp, 0).In(cstzone).Format(layout)
	return timeStr
}

type CacheWriter struct {
	root          string
	files         [13]*os.File
	deviceNo      string
	lastTimestamp uint64
	filename      string
	channel       uint16
}

func NewCacheWriter(path, deviceNo string) *CacheWriter {
	return &CacheWriter{
		root:          path,
		deviceNo:      deviceNo,
		lastTimestamp: 0,
		filename:      "",
	}
}

func NewCache(filename string) *CacheWriter {
	return &CacheWriter{
		lastTimestamp: 0,
		filename:      filename,
	}
}

func (c *CacheWriter) Ok() bool {
	return c.lastTimestamp > 0
}

func (c *CacheWriter) Close() (files []string) {
	if c.filename != "" {
		c.files[c.channel].Close()
		return
	}
	tstamp := int64(c.lastTimestamp / 1000 / 1000)
	dtStr := unix2Str(tstamp, "20060102 150405", 0)
	timeStr := dtStr[9:]
	for _, file := range c.files {
		if file == nil {
			continue
		}
		oldname := file.Name()
		file.Close()
		newname := fmt.Sprintf("%s_%s.cache", oldname, timeStr)
		if err := os.Rename(oldname, newname); err == nil {
			files = append(files, newname)
		} else {
			files = append(files, oldname)
		}
	}
	return
}

func (c *CacheWriter) createFile(channel uint16, timestamp uint64) error {
	tstamp := int64(timestamp / 1000 / 1000)
	dtStr := unix2Str(tstamp, "20060102 150405", 0)
	dateStr := dtStr[:8]
	timeStr := dtStr[9:]
	fpName := fmt.Sprintf("%s/%s/%s/ch%02d_%s_%s", c.root, dateStr, c.deviceNo, channel, dateStr, timeStr)
	if c.filename != "" {
		fpName = c.filename
		c.channel = channel
	}
	dir := filepath.Dir(fpName)
	os.MkdirAll(dir, os.ModePerm)
	fp, err := os.OpenFile(fpName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	c.files[channel] = fp
	return nil
}

func (c *CacheWriter) WriteFrame(channel, ctype uint16, timestamp uint64, data []byte) error {
	file := c.files[channel]
	if file == nil {
		if err := c.createFile(channel, timestamp); err != nil {
			return err
		}
	}
	res := cacheEncode(len(data))
	file.Write(res)
	file.Write(data)
	c.lastTimestamp = timestamp
	return nil
}
