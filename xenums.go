package xutils

import (
	"errors"
	"fmt"
)

var (
	ErrIsRunning = errors.New("is running")       //  服务已运行
	ErrParameter = errors.New("parameter error")  //  参数错误
	ErrLogin     = errors.New("login error")      //  登录错误
	ErrInvalid   = errors.New("this is invalid")  //  服务不存在
	ErrDisabled  = errors.New("this is disabled") //  服务已停止
	ErrAuthority = errors.New("authority error")  //  授权失败
	ErrIntnernal = errors.New("internal error")   //  内部错误(更新)
)

func ErrorHttp(url string) error {
	return fmt.Errorf("%s request error", url)
}

const (
	SERVE_StatusIdel     = iota // 空闲
	SERVE_StatusOk              // 正常
	SERVE_StatusDisabled        // 禁止
)
