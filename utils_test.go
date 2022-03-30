package xutils

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/wlgd/xutils/xnet"
)

func Test(t *testing.T) {
	fmt.Println(HostPublicAddr())
}

func handler(conn net.Conn, b []byte) error {
	data := make([]byte, 1024)
	for {
		n, err := conn.Read(data)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", data[:n])
	}
}

func TestXNet(t *testing.T) {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "hello world")
	})
	s := xnet.NewServe(8080)
	s.Handler = handler
	s.ListenAndServe()
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
