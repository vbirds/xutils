package xplayer

import "errors"

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