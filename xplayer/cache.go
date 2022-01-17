package xplayer

import (
	"encoding/binary"
	"time"

	"github.com/wlgd/xutils/file"
)

type Cache struct {
	status int
	offset uint64
	reader *file.CacheReader
}

func (c *Cache) Open(filename string) bool {
	c.reader = file.NewCacheReader(filename)
	return c.reader != nil
}

func (c *Cache) Pause() {
	c.status = kStatusPause
}

func (c *Cache) Resume() {
	c.status = kStatusPlay
}

func (c *Cache) Stop() {
	c.status = kStatusStop
}

func (c *Cache) Seek(offset uint64) {
	c.offset = offset
	c.status = kStatusSeek
}

func (c *Cache) StartPlay(handler func([]byte) error) error {
	if c.reader == nil {
		return Error_Player
	}
	var (
		lstFrameStamp uint64 = 0
		lstSendTime   time.Time
	)
	defer c.reader.Close()
	var err error = nil
	for {
		if c.status == kStatusPause {
			if time.Since(lstSendTime).Seconds() > 10 {
				err = Error_TimeOut
				break
			}
			time.Sleep(1 * time.Second)
			continue
		}
		if c.status == kStatusStop {
			break
		}
		h, frame, length := c.reader.ReadFrame(8)
		if length == 0 {
			err = Error_Frame
			break
		}
		// seek
		if c.status == kStatusSeek {
			if h.FrameType != 1 || h.Timestamp-lstFrameStamp < c.offset*1000 {
				continue
			}
			lstFrameStamp = 0
			c.status = kStatusPlay
		}
		// 按时间播放
		if h.FrameType < 3 {
			sec := time.Duration((h.Timestamp - lstFrameStamp) / 1000)
			if lstFrameStamp > 0 {
				time.Sleep(sec * time.Millisecond)
			}
			lstFrameStamp = h.Timestamp
		}
		// fmt.Printf("FrameType %d Channel %d Timestamp %v length %d bufferLength %d\n", h.FrameType, h.Channel, h.Timestamp, length, len(frame))
		// 封装头
		frame[0] = 'H'
		frame[1] = 1
		binary.LittleEndian.PutUint16(frame[2:], 1000)
		binary.LittleEndian.PutUint32(frame[4:], uint32(length))
		if err := handler(frame); err != nil {
			return err
		}
		lstSendTime = time.Now()
	}
	var emptybytes [8]byte
	emptybytes[0] = 'H'
	emptybytes[1] = 1
	binary.LittleEndian.PutUint16(emptybytes[2:], 1000)
	handler(emptybytes[:])
	// log.Printf("play %s closed\n", c.reader.FileName())
	return err
}
