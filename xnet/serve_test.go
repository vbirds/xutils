package xnet

import (
	"fmt"
	"net"
	"net/http"
	"testing"
)

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

func TestServer(t *testing.T) {
	s := NewServe(8088)
	s.ConnHandleFunc(handler)
	s.HttpHandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "hello world")
	})
	fmt.Println(s.ListenAndServe())
}
