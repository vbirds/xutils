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

type Serve struct {
	addr    *net.TCPAddr
	Handler func(net.Conn, []byte) error
}

func (s *Serve) newConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	defer conn.Close()
	data := make([]byte, 2048)
	recvlen, err := conn.Read(data)
	if err != nil || recvlen == 0 {
		return
	}
	r, err := httpRequest(data, conn)
	if err == nil {
		w := makeHttpWriter(conn.(*net.TCPConn))
		http.DefaultServeMux.ServeHTTP(w, r)
		return
	}
	fmt.Printf("%s connection success.", clientAddr)
	err = s.Handler(conn, data[:recvlen])
	fmt.Printf("%s connection closed. %v\n", clientAddr, err)
}

func (s *Serve) ListenAndServe() error {
	listener, err := net.ListenTCP("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}
		go s.newConnection(conn)
	}
}

func NewServe(port uint16) *Serve {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port)) //获取一个tcpAddr
	if err != nil {
		fmt.Println("Listener create error: ", err)
		return nil
	}
	return &Serve{addr: tcpAddr}
}
