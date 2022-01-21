package xplayer

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/wlgd/xutils/file"
)

const (
	kStatusPlay = iota
	kStatusStop
	kStatusPause
	kStatusSeek
)

var (
	Error_Player  = errors.New("player open")
	Error_TimeOut = errors.New("timeout")
	Error_Frame   = errors.New("frame read")
)

type Context interface {
	Open(string) bool
	Pause()
	Resume()
	Stop()
	Seek(uint64)
	StartPlay(func([]byte) error) error
}

type xPlayer struct {
	status int
	offset uint64
	reader file.Reader
}

func (c *xPlayer) Pause() {
	c.status = kStatusPause
}

func (c *xPlayer) Resume() {
	c.status = kStatusPlay
}

func (c *xPlayer) Stop() {
	c.status = kStatusStop
}

func (c *xPlayer) Seek(offset uint64) {
	c.offset = offset
	c.status = kStatusSeek
}

func (c *xPlayer) StartPlay(handler func([]byte) error) error {
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
