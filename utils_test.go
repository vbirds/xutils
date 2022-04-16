package xutils

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	fmt.Println(HostPublicAddr())
}

func TestBitmap(t *testing.T) {
	bmp := DefaultBitMap
	for i := 0; i < 500; i++ {
		bmp.Set(1000000 + i)
	}
	fmt.Println(bmp.Include(63))
	fmt.Println(bmp.Include(67))
	fmt.Println(bmp.All())
	fmt.Println(len(bmp.bits))
}
