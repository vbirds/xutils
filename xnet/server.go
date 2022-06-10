// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package xnet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
)

type httpWriter struct {
	conn *net.TCPConn
	rw   *bufio.ReadWriter
}

func makeHttpWriter(conn *net.TCPConn) *httpWriter {
	w := new(httpWriter)
	w.conn = conn
	w.rw = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	return w
}

func (rw *httpWriter) Header() http.Header {
	return make(http.Header)
}

func (rw *httpWriter) WriteHeader(int) {
}

func (rw *httpWriter) Write(b []byte) (int, error) {
	return 0, nil
}

func (rw *httpWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.conn, rw.rw, nil
}

func httpRequest(buffer []byte, conn net.Conn) (*http.Request, error) {
	r := io.MultiReader(bytes.NewReader(buffer), conn)
	return http.ReadRequest(bufio.NewReader(r))
}

// Server
type Server struct {
	addr        *net.TCPAddr
	listener    *net.TCPListener
	httpMux     *http.ServeMux
	connHandler func(net.Conn, []byte) error
}

func (s *Server) newConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	defer conn.Close()
	data := make([]byte, 2048)
	recvlen, err := conn.Read(data)
	if err != nil || recvlen == 0 {
		return
	}
	if bytes.Contains(data, []byte("http")) {
		if r, err := httpRequest(data, conn); err == nil {
			w := makeHttpWriter(conn.(*net.TCPConn))
			s.httpMux.ServeHTTP(w, r)
		}
		return
	}
	if s.connHandler != nil {
		err = s.connHandler(conn, data[:recvlen])
	}
	fmt.Printf("%s connection closed. %v\n", clientAddr, err)
}

func (s *Server) ListenAndServe() (err error) {
	s.listener, err = net.ListenTCP("tcp", s.addr)
	if err != nil {
		return err
	}
	defer s.listener.Close()
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			return err
		}
		go s.newConnection(conn)
	}
}

func (s *Server) Release() {
	s.listener.Close()
}

func (s *Server) HttpHandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.httpMux.HandleFunc(pattern, handler)
}

func (s *Server) ConnHandleFunc(handler func(net.Conn, []byte) error) {
	s.connHandler = handler
}

func NewServer(port uint16) *Server {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port)) //获取一个tcpAddr
	if err != nil {
		fmt.Println("Listener create error: ", err)
		return nil
	}
	return &Server{addr: tcpAddr, httpMux: http.NewServeMux()}
}
