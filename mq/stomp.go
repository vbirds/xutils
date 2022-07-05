// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mq

import (
	"encoding/json"
	"time"

	"github.com/go-stomp/stomp/v3"
)

type StompCli struct {
	Conn    *stomp.Conn
	ConnOpt []func(*stomp.Conn) error
	*Options
}

func NewStomp(opt *Options) (Interface, error) {
	cli := &StompCli{Options: opt}
	cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.HeartBeat(60*time.Second, 60*time.Second))
	if opt.User != "" {
		cli.ConnOpt = append(cli.ConnOpt, stomp.ConnOpt.Login(opt.User, opt.Pswd))
	}
	conn, err := stomp.Dial("tcp", opt.Address, cli.ConnOpt...)
	if err == nil {
		cli.Conn = conn
	}
	return cli, nil
}

func (o *StompCli) Publish(topic string, v interface{}) error {
	data, _ := json.Marshal(v)
	err := o.Conn.Send(topic, "text/plain", data)
	if err == nil {
		return err
	}
	o.Conn, err = stomp.Dial("tcp", o.Address, o.ConnOpt...)
	if err != nil {
		return err
	}
	return o.Conn.Send(topic, "text/plain", data)
}

func (o *StompCli) Subscribe(subject string, handler func([]byte) error) error {
	s, err := o.Conn.Subscribe(subject, stomp.AckClient)
	if err != nil {
		return err
	}
	go func() {
		for v := range s.C {
			if err := handler(v.Body); err == nil {
				o.Conn.Ack(v)
			}
		}
	}()
	return nil
}

func (o *StompCli) Run() error {
	return nil
}

func (o *StompCli) Release() {
	if o.Conn != nil {
		o.Conn.Disconnect()
	}
}
