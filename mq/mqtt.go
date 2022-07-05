// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mq

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttCli struct {
	Conn    mqtt.Client
	connOpt *mqtt.ClientOptions
	*Options
}

func (o *MqttCli) connectServe() error {
	o.Conn = mqtt.NewClient(o.connOpt)
	token := o.Conn.Connect()
	token.Wait()
	return token.Error()
}

func NewMqtt(opt *Options) (Interface, error) {
	cli := &MqttCli{Options: opt}
	cli.connOpt = mqtt.NewClientOptions()
	cli.connOpt.AddBroker(opt.Address)
	if opt.User != "" {
		cli.connOpt.SetUsername(opt.User)
		cli.connOpt.SetPassword(opt.Pswd)
	}
	if err := cli.connectServe(); err != nil {
		return cli, err
	}
	return cli, nil
}

func (o *MqttCli) Publish(topic string, v interface{}) error {
	if !o.Conn.IsConnected() {
		if err := o.connectServe(); err != nil {
			return err
		}
	}
	data, _ := json.Marshal(v)
	return o.Conn.Publish("topic/"+topic, 0, false, data).Error()
}

func (o *MqttCli) Subscribe(subject string, handler func([]byte) error) error {
	o.Conn.Subscribe("topic/"+subject, 1, func(c mqtt.Client, m mqtt.Message) {
		if err := handler(m.Payload()); err == nil {
			m.Ack()
		}
	})
	return nil
}

func (o *MqttCli) Run() error {
	return nil
}

func (o *MqttCli) Release() {
	if o.Conn != nil {
		o.Conn.Disconnect(1)
	}
}
