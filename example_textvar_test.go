// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"net"
	"os"

	"github.com/zulucmd/zflag/v2"
)

// Copyright 2022 The Go Authors. All rights reserved.
func ExampleTextVar() {
	fs := zflag.NewFlagSet("ExampleTextVar", zflag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var ip net.IP
	fs.TextVar(&ip, "ip", net.IPv4(192, 168, 0, 100), "`IP address` to parse")
	_ = fs.Parse([]string{"--ip", "127.0.0.1"})
	fmt.Printf("{ip: %v}\n\n", ip)

	// 256 is not a valid IPv4 component
	ip = nil
	_ = fs.Parse([]string{"--ip", "256.0.0.1"})
	fmt.Printf("{ip: %v}\n\n", ip)

	// Output:
	// {ip: 127.0.0.1}
	//
	// Usage of ExampleTextVar:
	//       --ip IP address   IP address to parse (default 192.168.0.100)
	//
	// invalid argument "256.0.0.1" for "--ip" flag: invalid IP address: 256.0.0.1
	// {ip: <nil>}
}
