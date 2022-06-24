// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package xnet

import (
	"fmt"
	"log"
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
	s, err := NewServer(8088)
	if err != nil {
		log.Fatalln(err)
	}
	s.ConnHandleFunc(handler)
	s.HttpHandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(rw, "hello world")
	})
	fmt.Println(s.ListenAndServe())
}
