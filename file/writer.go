package file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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
