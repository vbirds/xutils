package xutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

const (
	contentType = "application/json;charset=utf8"
)

type response struct {
	Status int         `json:"status"`
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
	if res.Status != 10000 {
		return errors.New(res.Msg)
	}
	if result != nil {
		jsoniter.Get(content, "data").ToVal(result)
	}
	return nil
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
	if res.Status != 10000 {
		return errors.New(res.Msg)
	}
	if result != nil {
		jsoniter.Get(content, "data").ToVal(result)
	}
	return nil
}

// ServerOpt 服务配置信息
type ServerOpt struct {
	Name       string `json:"name"`
	HttpPort   uint16 `json:"httpPort"`
	AccessPort uint16 `json:"accessPort"`
	RpcPort    uint16 `json:"rpcPort"`
	Status     int    `json:"status"`
	Address    string // 服务IP
}

// Services 服务信息
type Server struct {
	Local   ServerOpt `json:"local"`
	Station ServerOpt `json:"station"`
}

func HttpApplyAuth(url, serverId string) (*Server, error) {
	if url == "" {
		return nil, errors.New("please set authority address firstly")
	}
	address := fmt.Sprintf("%s/%s", url, serverId)
	s := &Server{}
	if err := HttpGet(address, s); err != nil {
		return nil, ErrorHttp(address)
	}
	if s.Local.Status != SERVE_StatusOk {
		return nil, ErrDisabled
	}
	return s, nil
}

type LogAlarmLink struct {
	ServerID         string `json:"serverID"`
	DeviceNo         string `json:"deviceNo"`
	AlarmType        int    `json:"alarmType"`
	AlarmGuid        string `json:"alarmGuid"`
	ResStartTime     string `json:"resStartTime"`
	ResEndTime       string `json:"resEndTime"`
	ResRealStartTime string `json:"resRealStartTime"`
	ResRealEndTime   string `json:"resRealEndTime"`
	ExeStartTime     string `json:"exeStartTime"`
	ExeEndTime       string `json:"exeEndTime"`
	Action           int    `json:"action"`
}

func HttpPostLogAlarmLink(url string, v *LogAlarmLink) error {
	return HttpPost(url, v, nil)
}
