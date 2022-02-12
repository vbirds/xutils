package xutils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// JSONConf 初始化配置参数
func JSONConf(jsonFile string, obj interface{}) {
	jsonFp, err := os.Open(jsonFile)
	if err != nil {
		fmt.Println("load error" + jsonFile)
		os.Exit(0)
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
	if err := json.Unmarshal([]byte(jsString), obj); err != nil {
		fmt.Println("json error " + jsonFile)
		os.Exit(0)
	}
}

func YMLConf(fpname string, obj interface{}) error {
	yfile, err := ioutil.ReadFile(fpname)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yfile, obj)
}
