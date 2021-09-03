package rpc

import (
	"time"
)

// 和中心服务通信
// func name <Register>

// KeepAliveArgs 保活
type KeepAliveArgs struct {
	ServeId     string    `json:"serveId"`     // 服务ID
	Token       string    `json:"token"`       // 授权token
	UpdatedTime time.Time `json:"updatedTime"` // 更新时间
}

// LoginArgs 登录工作站
type LoginArgs struct {
	ServeId string `json:"serveId"` // 服务ID
	Address string `json:"address"` // 服务名称
}

// LoginReply 工作站响应
type LoginReply struct {
	Token string `json:"token"` // 授权token
}

// XLinkRegister 设备链路注册
// for external server
type XLinkRegister struct {
	Data interface{} `json:"data"`
}
