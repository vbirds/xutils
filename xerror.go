package xutils

import (
	"errors"
)

var (
	ErrServeIsRunning = errors.New("serve is running")              //  服务已运行
	ErrHttpURL        = errors.New("http url error")                //  请求URL错误
	ErrPostRequest    = errors.New("request error")                 //  post请求错误
	ErrParameter      = errors.New("parameter error")               //  参数错误
	ErrLogin          = errors.New("login error")                   //  登录错误
	ErrServeNoExist   = errors.New("serve no exist")                //  服务不存在
	ErrServeStoped    = errors.New("the service has been disabled") //  服务已停止
	ErrAuthority      = errors.New("authority error")               //  授权失败
	ErrIntnernal      = errors.New("internal error")                //  内部错误(更新)
)
