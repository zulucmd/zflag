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

func TestIP(t *testing.T) {
	tests := []struct {
		name        string
		flagDefault net.IP
		input       []string
		expectedErr string
		expectedIPs net.IP
	}{
		{
			name:        "no value passed",
			input:       []string{},
			flagDefault: net.IP{},
			expectedErr: "",
			expectedIPs: net.IP{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: net.IP{},
			expectedErr: `invalid argument "" for "--ip" flag: failed to parse IP: ""`,
		},
		{
			name:        "invalid ip",
			input:       []string{"blabla"},
			flagDefault: net.IP{},
			expectedErr: `invalid argument "blabla" for "--ip" flag: failed to parse IP: "blabla"`,
		},
		{
			name:        "no csv",
			input:       []string{"192.168.1.1,172.16.1.1"},
			flagDefault: net.IP{},
			expectedErr: `invalid argument "192.168.1.1,172.16.1.1" for "--ip" flag: failed to parse IP: "192.168.1.1,172.16.1.1"`,
		},
		{
			name:        "multiple value passed",
			input:       []string{"192.168.1.1", "10.0.0.1"},
			flagDefault: net.IP{},
			expectedIPs: net.ParseIP("10.0.0.1"),
		},
		{
			name:        "overrides default values",
			input:       []string{"0:0:0:0:0:0:0:2"},
			flagDefault: net.ParseIP("192.168.1.1"),
			expectedIPs: net.ParseIP("0:0:0:0:0:0:0:2"),
		},
		{
			name:        "with default values",
			input:       []string{},
			flagDefault: net.ParseIP("192.168.1.1"),
			expectedIPs: net.ParseIP("192.168.1.1"),
		},
		{
			name:        "trims input",
			input:       []string{"    192.168.1.1    "},
			flagDefault: net.IP{},
			expectedIPs: net.ParseIP("192.168.1.1"),
		},
		{
			name:        "trims input",
			input:       []string{"    0:0:0:0:0:0:0:2    "},
			flagDefault: net.IP{},
			expectedIPs: net.ParseIP("0:0:0:0:0:0:0:2"),
		},
		{
			// https://github.com/spf13/pflag/issues/351
			name:        "allows nil value",
			input:       []string{},
			flagDefault: nil,
			expectedIPs: nil,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ip net.IP
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.IPVar(&ip, "ip", test.flagDefault, "usage")
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

			if !reflect.DeepEqual(test.expectedIPs, ip) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedIPs, ip)
			}

			getIPS, err := f.GetIP("ip")
			if err != nil {
				t.Fatal("got an error from GetIP():", err)
			}
			if !reflect.DeepEqual(test.expectedIPs, getIPS) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", getIPS, test.expectedIPs)
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
