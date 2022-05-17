// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package xutils

import (
	"fmt"
	"net"
	"strings"
)

var publicIPAddr string = ""

// HostPublicAddr public address
func HostPublicAddr() string {
	if publicIPAddr != "" {
		return publicIPAddr
	}
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer conn.Close()
	publicIPAddr = strings.Split(conn.LocalAddr().String(), ":")[0]
	return publicIPAddr
}

var localIPAddr string = ""

// HostLocalAddr local address
func HostLocalAddr() string {
	if localIPAddr != "" {
		return localIPAddr
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIPAddr = ipnet.IP.String()
				break
			}
		}
	}
	return localIPAddr
}
