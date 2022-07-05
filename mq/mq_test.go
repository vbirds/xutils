// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mq

import (
	"log"
	"testing"
	"time"
)

func TestMqtt(t *testing.T) {
	c, err := NewPublish(&Options{Address: "127.0.0.1:35003", Goc: 5}, NewMqtt)
	if err != nil {
		c.Shutdown()
		log.Fatalln(err)
	}
	topic := "test/mqtt"
	c.Subscribe(topic, func(b []byte) error {
		log.Printf("%s: %s\n", topic, b)
		return nil
	})
	for {
		time.Sleep(2 * time.Second)
		c.Publish(topic, map[string]interface{}{
			"device": "20198002",
			"now":    time.Now().Format("2006-01-02 15:04:05"),
		})
	}
}

func TestStomp(t *testing.T) {
	c, err := NewPublish(&Options{Address: "127.0.0.1:35002", Goc: 1}, NewStomp)
	if err != nil {
		c.Shutdown()
	}
	subject := "/queue/test/stomp"
	c.Subscribe(subject, func(b []byte) error {
		log.Printf("%s:1 %s\n", subject, b)
		return nil
	})
	c.Subscribe(subject, func(b []byte) error {
		log.Printf("%s:2 %s\n", subject, b)
		return nil
	})
	for {
		time.Sleep(2 * time.Second)
		c.Publish(subject, map[string]interface{}{
			"device": "20198002",
			"now":    time.Now().Format("2006-01-02 15:04:05"),
		})
	}
}

func TestNats(t *testing.T) {
	c, err := NewPublish(&Options{Address: NatsURL, Goc: 1}, NewStomp)
	if err != nil {
		c.Shutdown()
	}
	subject := "/queue/test/nats"
	c.Subscribe(subject, func(b []byte) error {
		log.Printf("%s:1 %s\n", subject, b)
		return nil
	})
	c.Subscribe(subject, func(b []byte) error {
		log.Printf("%s:2 %s\n", subject, b)
		return nil
	})
	for {
		time.Sleep(2 * time.Second)
		c.Publish(subject, map[string]interface{}{
			"device": "20198002",
			"now":    time.Now().Format("2006-01-02 15:04:05"),
		})
	}
}
