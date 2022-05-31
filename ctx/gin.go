// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package ctx

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	StatusFail         = -1  // 失败
	StatusOK           = 200 // 成功
	StatusError        = 500 // 错误
	StatusLoginExpired = 401 // 登录过期
	StatusForbidden    = 403 // 无权限
)

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Context 响应
type Context struct {
	response
}

// JSON 响应
func JSON(status int) *Context {
	ctx := &Context{}
	ctx.Code = status
	switch status {
	case StatusOK:
		ctx.Msg = "success"
	case StatusFail:
		ctx.Msg = "failed"
	case StatusForbidden:
		ctx.Msg = "forbidden"
	}
	return ctx
}

// SetMsg 设置消息体的内容int
func (o *Context) SetMsg(msg string) *Context {
	o.Msg = msg
	return o
}

// SetCode 设置消息体的编码
func (o *Context) SetCode(code int) *Context {
	o.Code = code
	return o
}

// WriteData 输出json到客户端， 有data字段
func (o *Context) WriteData(data interface{}, c *gin.Context) {
	o.Data = data
	c.JSON(http.StatusOK, o.response)
}

// Write 输出json到客户端, 无data字段
func (o *Context) Write(data gin.H, c *gin.Context) {
	data["code"] = o.Code
	data["msg"] = o.Msg
	c.JSON(http.StatusOK, data)
}

// JSONOk 无数据响应
func JSONOk(c *gin.Context) {
	JSON(StatusOK).WriteData(nil, c)
}

// JSONWrite
func JSONWrite(data gin.H, c *gin.Context) {
	JSON(StatusOK).Write(data, c)
}

// JSONWriteData
func JSONWriteData(v interface{}, c *gin.Context) {
	JSON(StatusOK).WriteData(v, c)
}

// JSONError
func JSONError(c *gin.Context) {
	JSON(StatusError).WriteData(nil, c)
}

// JSONWriteError 错误应答
func JSONWriteError(err error, c *gin.Context) {
	JSON(StatusError).SetMsg(err.Error()).WriteData(nil, c)
}

// ParamUInt uint参数
func ParamUInt(c *gin.Context, key string) uint {
	idstr := c.Param(key)
	id, _ := strconv.Atoi(idstr)
	return uint(id)
}

// ParamInt int参数
func ParamInt(c *gin.Context, key string) int {
	return int(ParamUInt(c, key))
}

// QueryInt int参数
func QueryInt(c *gin.Context, key string) (int, error) {
	idstr := c.Query(key)
	return strconv.Atoi(idstr)
}

// QueryUInt int参数
func QueryUInt(c *gin.Context, key string) (uint, error) {
	idstr := c.Query(key)
	id, err := strconv.Atoi(idstr)
	return uint(id), err
}

// QueryUInt64 int参数
func QueryUInt64(c *gin.Context, key string) (uint64, error) {
	idstr := c.Query(key)
	id, err := strconv.Atoi(idstr)
	return uint64(id), err
}
