package file

type Writer interface {
	WriteFrame(channel, ctype uint16, timestamp uint64, data []byte) error
	Ok() bool
	Close() (files []string)
}
