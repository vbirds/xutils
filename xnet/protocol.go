// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package xnet

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
)

type Header struct {
	Flag    uint8  //标识，固定为'H'
	Version uint8  //版本，当前协议版本为1
	Code    uint16 //消息代码
	Length  uint32 //负载长度(不包含当前消息头长度)
}

func MsgHeader(code uint16, b []byte) []byte {
	if b == nil {
		b = make([]byte, 8)
	}
	b[0] = 'H'
	b[1] = 1
	binary.LittleEndian.PutUint16(b[2:], code)
	binary.LittleEndian.PutUint32(b[4:], uint32(len(b)-8))
	return b
}

func MsgPack(code uint16, v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var m Msg
	return m.Pack(code, data)
}
func MsgUnPack(b []byte) (*Header, []byte) {
	h := &Header{}
	h.Flag = b[0]
	h.Version = b[1]
	h.Code = binary.LittleEndian.Uint16(b[2:])
	h.Length = binary.LittleEndian.Uint32(b[4:])
	return h, b[8:]
}

type Msg struct {
	data []byte
}

func (o *Msg) Pack(code uint16, b []byte) []byte {
	datalen := len(b) + 8
	if cap(o.data) < datalen {
		o.data = make([]byte, datalen+1)
	}
	if b != nil {
		copy(o.data[8:], b)
	}
	if code == MsgJSON {
		o.data[datalen] = '\x00'
		datalen += 1
	}
	return MsgHeader(code, o.data[:datalen])
}

func (m *Msg) UnPack(b []byte) (*Header, []byte) {
	return MsgUnPack(b)
}

func ShouldBindJSON(b []byte, v interface{}) error {
	if b == nil {
		return errors.New("recv no data")
	}
	data := bytes.TrimRight(b, "\x00")
	return json.Unmarshal(data, v)
}
