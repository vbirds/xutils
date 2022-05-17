// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package xutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

const (
	contentType = "application/json;charset=utf8"
)

type response struct {
	Status int         `json:"status"`
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// httpPost http post 请求
func HttpPost(url string, requset interface{}, result interface{}) error {
	bs, err := json.Marshal(requset)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(bs)
	recv, err := http.Post(url, contentType, body)
	if err != nil {
		return err
	}
	defer recv.Body.Close()
	content, err := ioutil.ReadAll(recv.Body)
	if err != nil {
		return err
	}
	var res response
	if err := jsoniter.Unmarshal(content, &res); err != nil {
		return err
	}
	if res.Status == 10000 || res.Code == 200 {
		if result != nil && res.Data != nil {
			jsoniter.Get(content, "data").ToVal(result)
		}
		return nil
	}
	return errors.New(res.Msg)
}

// HttpGet http get 请求
func HttpGet(url string, result interface{}) error {
	recv, err := http.Get(url)
	if err != nil {
		return err
	}
	defer recv.Body.Close()
	content, err := ioutil.ReadAll(recv.Body)
	if err != nil {
		return err
	}
	var res response
	if err := jsoniter.Unmarshal(content, &res); err != nil {
		return err
	}
	if res.Status == 10000 || res.Code == 200 {
		if result != nil && res.Data != nil {
			jsoniter.Get(content, "data").ToVal(result)
		}
		return nil
	}
	return errors.New(res.Msg)
}
