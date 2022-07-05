// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mq

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type AmqpCli struct {
	Conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	routing  string
	*Options
}

func NewAmqp(opt *Options) (Interface, error) {
	cli := &AmqpCli{Options: opt}
	conn, err := amqp.Dial(cli.Address)
	if err != nil {
		return cli, err
	}
	cli.Conn = conn
	cli.channel, err = conn.Channel()
	if err != nil {
		return cli, err
	}
	return cli, nil
}

func (o *AmqpCli) Publish(topic string, v interface{}) error {
	if _, err := o.channel.QueueDeclare(topic, true, false, false, false, nil); err != nil {
		return err
	}
	data, _ := json.Marshal(v)
	o.channel.Publish(o.exchange, o.routing, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	return nil
}

func (o *AmqpCli) Subscribe(subject string, handler func([]byte) error) error {
	if _, err := o.channel.QueueDeclare(subject, true, false, false, false, nil); err != nil {
		return err
	}
	msgChan, err := o.channel.Consume(subject, "", true, false, false, true, nil)
	if err != nil {
		return err
	}
	go func() {
		for v := range msgChan {
			if err := handler(v.Body); err == nil {
				v.Ack(true)
			}
		}
	}()
	return nil
}

func (o *AmqpCli) Run() error {
	return o.channel.ExchangeDeclare(o.exchange, "topic", true, false, false, false, nil)
}

func (o *AmqpCli) Release() {
	if o.Conn != nil {
		o.Conn.Close()
		o.channel.Close()
	}
}
