// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestIPNet(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    net.IPNet
		input          []string
		expectedErr    string
		expectedValues net.IPNet
	}{
		{
			name:           "valid cidr",
			input:          []string{"0.0.0.0/0"},
			expectedValues: getCIDR("0.0.0.0/0"),
		},
		{
			name:           "valid cidr",
			input:          []string{"1.2.3.4/8"},
			expectedValues: getCIDR("1.0.0.0/8"),
		},
		{
			name:           "valid cidr",
			input:          []string{"127.0.0.1/16"},
			expectedValues: getCIDR("127.0.0.0/16"),
		},
		{
			name:           "valid mask",
			input:          []string{"255.255.255.255/19"},
			expectedValues: getCIDR("255.255.224.0/19"),
		},
		{
			name:           "valid mask",
			input:          []string{"255.255.255.255/32"},
			expectedValues: getCIDR("255.255.255.255/32"),
		},
		{
			name:        "invalid cidr",
			input:       []string{"/0"},
			expectedErr: "invalid CIDR address: /0",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0"},
			expectedErr: "invalid CIDR address: 0",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0/0"},
			expectedErr: "invalid CIDR address: 0/0",
		},
		{
			name:        "invalid cidr",
			input:       []string{"localhost/0"},
			expectedErr: "invalid CIDR address: localhost/0",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0.0.0/4"},
			expectedErr: "invalid CIDR address: 0.0.0/4",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0.0.0./8"},
			expectedErr: "invalid CIDR address: 0.0.0./8",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0.0.0.0./12"},
			expectedErr: "invalid CIDR address: 0.0.0.0./12",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0.0.0.256/16"},
			expectedErr: "invalid CIDR address: 0.0.0.256/16",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0.0.0.0 /20"},
			expectedErr: "invalid CIDR address: 0.0.0.0 /20",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0.0.0.0/ 24"},
			expectedErr: "invalid CIDR address: 0.0.0.0/ 24",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0 . 0 . 0 . 0 / 28"},
			expectedErr: "invalid CIDR address: 0 . 0 . 0 . 0 / 28",
		},
		{
			name:        "invalid cidr",
			input:       []string{"0.0.0.0/33"},
			expectedErr: `invalid argument "0.0.0.0/33" for "--ip" flag: invalid CIDR address: 0.0.0.0/33`,
		},

		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    net.IPNet{},
			expectedValues: net.IPNet{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: net.IPNet{},
			expectedErr: `invalid argument "" for "--ip" flag: invalid CIDR address:`,
		},
		{
			name:        "no csv",
			input:       []string{"192.168.1.1/32,172.16.1.1/32"},
			flagDefault: net.IPNet{},
			expectedErr: `invalid argument "192.168.1.1/32,172.16.1.1/32" for "--ip" flag: invalid CIDR address: 192.168.1.1/32,172.16.1.1/32`,
		},
		{
			name:           "multiple value passed",
			input:          []string{"192.168.1.1/32", "10.0.0.1/32"},
			flagDefault:    net.IPNet{},
			expectedValues: getCIDR("10.0.0.1/32"),
		},
		{
			name:           "overrides default values",
			input:          []string{"0:0:0:0:0:0:0:2/64"},
			flagDefault:    getCIDR("192.168.1.1/32"),
			expectedValues: getCIDR("0:0:0:0:0:0:0:2/64"),
		},
		{
			name:           "with default values",
			input:          []string{},
			flagDefault:    getCIDR("192.168.1.1/24"),
			expectedValues: getCIDR("192.168.1.1/24"),
		},
		{
			name:           "trims input",
			input:          []string{"    192.168.1.1/24    "},
			flagDefault:    net.IPNet{},
			expectedValues: getCIDR("192.168.1.1/24"),
		},
		{
			name:           "trims input",
			input:          []string{"    0:0:0:0:0:0:0:2/64    "},
			flagDefault:    net.IPNet{},
			expectedValues: getCIDR("0:0:0:0:0:0:0:2/64"),
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ip net.IPNet
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.IPNetVar(&ip, "ip", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ip", test.input...))
			if test.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected an error; got none")
				}
				if test.expectedErr != "" && !strings.Contains(err.Error(), test.expectedErr) {
					t.Fatalf("expected error to contain %q, but was: %s", test.expectedErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error; got %q", err)
			}

			if !reflect.DeepEqual(test.expectedValues, ip) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, ip)
			}

			getIPS, err := f.GetIPNet("ip")
			if err != nil {
				t.Fatal("got an error from GetIP():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, getIPS) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", getIPS, test.expectedValues)
			}

			getIPSGet, err := f.Get("ip")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(getIPSGet, getIPS) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", getIPS, getIPSGet)
			}
		})
	}
}
