// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mq

type Interface interface {
	Publish(string, interface{}) error
	Subscribe(string, func([]byte)) error
	Run()
	Release()
}

type Options struct {
	Address string
	Goc     int // goroutine 数目
}

type Client struct {
	clients []Interface
	Options
	goi int // 当前使用的goroutine
}

func (o *Client) goindex() int {
	o.goi++
	if o.goi >= len(o.clients) {
		o.goi = 0
	}
	return o.goi
}

// 支持创建多goroutine发布
func New(o *Options, handler func(string) (Interface, error)) (*Client, error) {
	if o.Goc == 0 {
		o.Goc = 1
	}
	c := &Client{Options: *o}
	for i := 0; i < o.Goc; i++ {
		cli, err := handler(o.Address)
		if err != nil {
			return c, err
		}
		c.clients = append(c.clients, cli)
	}
	return c, nil
}

func NewPublish(o *Options, handler func(string) (Interface, error)) (*Client, error) {
	c, err := New(o, handler)
	if err != nil {
		c.Shutdown()
		return c, err
	}
	for _, v := range c.clients {
		v.Run()
	}
	return c, nil
}

func (o *Client) Shutdown() {
	for _, v := range o.clients {
		if v != nil {
			v.Release()
		}
	}
}

// 订阅会自动分配connection对象
// 订阅数大于连接数，出现同意连接多次订阅，报错
func (o *Client) Subscribe(subject string, goc int, handler func([]byte)) error {
	i := o.goindex()
	return o.clients[i].Subscribe(subject, handler)
}

// 动态均衡，自动适配connection发送数据
func (o *Client) Publish(topic string, v interface{}) {
	i := o.goindex()
	o.clients[i].Publish(topic, v)
}
