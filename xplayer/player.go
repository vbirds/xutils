package xplayer

const (
	kStatusPlay = iota
	kStatusStop
	kStatusPause
	kStatusSeek
)

type Context interface {
	Open(string) bool
	Pause()
	Resume()
	Stop()
	Seek(uint64)
	StartPlay(func([]byte) error)
}