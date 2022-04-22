package xutils

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// JSONConf 初始化配置参数
func JSONConf(jsonFile string, obj interface{}) error {
	jsonFp, err := os.Open(jsonFile)
	if err != nil {
		return err
	}
	defer jsonFp.Close()
	var jsString string
	iReader := bufio.NewReader(jsonFp)
	for {
		tString, err := iReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		jsString = jsString + tString
	}
	return json.Unmarshal([]byte(jsString), obj)
}

func YMLConf(fpname string, obj interface{}) error {
	yfile, err := ioutil.ReadFile(fpname)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yfile, obj)
}
