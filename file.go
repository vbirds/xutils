package xutils

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// JSONFile 初始化配置参数
func JSONFile(filename string, obj interface{}) error {
	jsonFp, err := os.Open(filename)
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

func YAMLFile(filename string, obj interface{}) error {
	yfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yfile, obj)
}
