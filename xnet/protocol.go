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

const _headerLen = 8

type Header struct {
	Flag    uint8  //标识，固定为'H'
	Version uint8  //版本，当前协议版本为1
	Code    uint16 //消息代码
	Length  uint32 //负载长度(不包含当前消息头长度)
}

func MsgHeader(code uint16, b []byte) []byte {
	if b == nil {
		b = make([]byte, _headerLen)
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
	return h, b[_headerLen:]
}

type Msg struct {
	sdata  []byte
	rdata  bytes.Buffer
	pkglen uint32
}

func (o *Msg) Pack(code uint16, b []byte) []byte {
	datalen := len(b) + _headerLen + 1
	if cap(o.sdata) < datalen {
		o.sdata = make([]byte, datalen)
	}
	if b != nil {
		copy(o.sdata[_headerLen:], b)
	}
	if code == MsgJSON {
		o.sdata[datalen-1] = '\x00'
	} else {
		datalen--
	}
	return MsgHeader(code, o.sdata[:datalen])
}

func (o *Msg) UnPack(b []byte) (*Header, []byte) {
	if o.pkglen > 0 {
		o.rdata.Next(int(o.pkglen))
		o.pkglen = 0
	}
	o.rdata.Write(b)
	rlen := o.rdata.Len()
	if rlen < _headerLen {
		return nil, nil
	}
	h, data := MsgUnPack(o.rdata.Bytes())
	datalen := h.Length + _headerLen
	if rlen < int(datalen) {
		return nil, nil
	}
	o.pkglen = datalen
	return h, data
}

func ShouldBindJSON(b []byte, v interface{}) error {
	if b == nil {
		return errors.New("recv no data")
	}
	data := bytes.TrimRight(b, "\x00")
	return json.Unmarshal(data, v)
}
