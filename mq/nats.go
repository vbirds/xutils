// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mq

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	NatsURL = nats.DefaultURL
)

type pubMsg struct {
	topic string
	data  interface{}
}

type NatsCli struct {
	Conn *nats.Conn
	msgs chan *pubMsg
	*Options
}

func (o *NatsCli) isConnected() error {
	if o.Conn == nil || o.Conn.IsClosed() {
		c, err := nats.Connect(o.Address)
		if err != nil {
			return err
		}
		o.Conn = c
	}
	return nil
}

// NewNats nats connect
func NewNats(opt *Options) (Interface, error) {
	cli := &NatsCli{Options: opt}
	conn, err := nats.Connect(opt.Address)
	if err != nil {
		return cli, err
	}
	cli.Conn = conn
	return cli, nil
}

func (o *NatsCli) Run() error {
	o.msgs = make(chan *pubMsg, 10)
	go func() {
		for v := range o.msgs {
			if v == nil {
				break
			}
			if err := o.isConnected(); err == nil {
				data, _ := json.Marshal(v.data)
				o.Conn.Publish(v.topic, data)
			}
		}
	}()
	return nil
}

// Publish publish
func (o *NatsCli) Publish(topic string, v interface{}) error {
	o.msgs <- &pubMsg{topic: topic, data: v}
	return nil
}

// Subscribe subscribe
func (o *NatsCli) Subscribe(topic string, handler func([]byte) error) error {
	_, err := o.Conn.Subscribe(topic, func(m *nats.Msg) {
		if err := handler(m.Data); err == nil {
			m.Ack()
		}
	})
	return err
}

// SubscribeRsp response the request
func (o *NatsCli) SubscribeRsp(topic string, handler func([]byte) []byte) error {
	_, err := o.Conn.Subscribe(topic, func(m *nats.Msg) {
		rsp := handler(m.Data)
		o.Conn.Publish(m.Reply, rsp)
	})
	return err
}

// Request return reponse data, error
func (o *NatsCli) Request(topic string, data []byte, msec time.Duration) ([]byte, error) {
	msg, err := o.Conn.Request(topic, data, msec*time.Microsecond)
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}

func (o *NatsCli) Release() {
	if o.Conn == nil {
		return
	}
	if o.msgs != nil {
		o.msgs <- nil
	}
	o.Conn.Drain()
}
