// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package xnet

const (
	ActUnknow     = 0
	ActHeartbeat  = 1
	ActCommonResp = 2
)

const (
	ErrNo                 = "0"  //成功
	ErrInvalidParameter   = "1"  //无效参数
	ErrInvalidCommand     = "2"  //无效命令
	ErrInvalidService     = "3"  //无效服务
	ErrServiceUnstart     = "4"  //服务未启动
	ErrObjectNotExist     = "5"  //对象不存在
	ErrWaitingTimeout     = "6"  //等待超时
	ErrConnectedFailed    = "7"  //连接失败
	ErrFileOpenFailed     = "8"  //文件打开失败
	ErrDeviceOffline      = "9"  //设备不在线
	ErrDeviceBusing       = "10" //设备忙
	ErrNetworkException   = "11" //网络异常
	ErrDiskInitFailed     = "12" //获取磁盘媒体块失败
	ErrDiskReadFailed     = "13" //读取媒体块失败
	ErrDiskOpenFileFailed = "14" //打开媒体文件失败
	ErrUnsupport          = "25" //不支持
)

//消息代码
const (
	MsgUnknow    = 0    //无效
	MsgHeartbeat = 1    //心跳
	MsgJSON      = 2    //json over tcp
	MsgHTTP      = 3    //http request+flv muxer
	MsgWebsocket = 5    //websocket
	MsgMedia     = 1000 //媒体数据
)

const (
	LogUnknow           = 0
	LogRecvTask         = 1
	LogStartDownload    = 2
	LogDownloadSuccess  = 3
	LogErrorDevRegister = 4
	LogErrorRetries     = 5
	LogErrorRcvData     = 6
	LogQueryNofiles     = 7
)
