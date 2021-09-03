package xutils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func YMLConf(fpname string, obj interface{}) error {
	yfile, err := ioutil.ReadFile(fpname)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yfile, obj)
}
