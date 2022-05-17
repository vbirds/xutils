// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package xutils

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// JSONFile load json file
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

// YAMLFile load yaml file
func YAMLFile(filename string, obj interface{}) error {
	yfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yfile, obj)
}
