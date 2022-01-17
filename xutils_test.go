package xutils

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/wlgd/xutils/xnet"
)


func Test(t *testing.T) {
	fmt.Println(HostPublicAddr())
}

func TestXNet(t *testing.T) {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "hello world")
	})
	s := xnet.NewServe(8080)
	s.ListenAndServe()
}