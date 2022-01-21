package xplayer

import (
	"github.com/wlgd/xutils/file"
)

type Cache struct {
	xPlayer
}

func (c *Cache) Open(filename string) bool {
	c.reader = file.NewCacheReader(filename)
	return c.reader != nil
}
