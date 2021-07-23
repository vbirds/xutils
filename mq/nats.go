package mq

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	DefaultURL = nats.DefaultURL
)

type NatsClient struct {
	url     string
	Conn    *nats.Conn
	pubs    chan interface{}
	onlysub bool
}

// natsPublish pubilsh data struct
type natsPublish struct {
	Topic string
	Data  interface{}
}

type NatsMsgHandler func([]byte)
type NatsMsgRspHandler func([]byte) []byte

func (n *NatsClient) checkNetStatus() error {
	if n.Conn == nil || n.Conn.IsClosed() {
		nc, err := nats.Connect(n.url)
		if err != nil {
			return err
		}
		n.Conn = nc
	}
	return nil
}

// NewNatsClient nats connect
func NewNatsClient(url string, onlysub bool) (*NatsClient, error) {
	cli := &NatsClient{}
	cli.url = url
	nc, err := nats.Connect(url)
	if err != nil {
		return cli, err
	}
	cli.Conn = nc
	cli.onlysub = onlysub
	if onlysub {
		return cli, nil
	}
	cli.pubs = make(chan interface{}, 1)
	go func() {
		for v := range cli.pubs {
			if v == nil {
				return
			}
			pub, _ := v.(natsPublish)
			bs, _ := json.Marshal(pub.Data)
			if err := cli.checkNetStatus(); err == nil {
				cli.Conn.Publish(pub.Topic, bs)
			}
		}
	}()
	return cli, nil
}

// Publish publish
func (n *NatsClient) Publish(topic string, data interface{}) error {
	if n.onlysub {
		return errors.New("unsupport pubilsh")
	}
	np := natsPublish{
		Topic: topic,
		Data:  data,
	}
	n.pubs <- np
	return nil
}

// Subscribe subscribe
func (n *NatsClient) Subscribe(topic string, handler NatsMsgHandler) error {
	_, err := n.Conn.Subscribe(topic, func(m *nats.Msg) {
		handler(m.Data)
	})
	return err
}

// SubscribeRsp response the request
func (n *NatsClient) SubscribeRsp(topic string, handler NatsMsgRspHandler) error {
	_, err := n.Conn.Subscribe(topic, func(m *nats.Msg) {
		rsp := handler(m.Data)
		n.Conn.Publish(m.Reply, rsp)
	})
	return err
}

// Request return reponse data, error
func (n *NatsClient) Request(topic string, data []byte, msec time.Duration) ([]byte, error) {
	msg, err := n.Conn.Request(topic, data, msec*time.Microsecond)
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}

func (n *NatsClient) Release() {
	if n.Conn == nil {
		return
	}
	if !n.onlysub {
		n.pubs <- nil
	}
	n.Conn.Drain()
}
