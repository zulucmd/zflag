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

// Helper function to set static slices
func getCIDR(val string) net.IPNet {
	_, cidr, _ := net.ParseCIDR(val)
	return *cidr
}

func TestIPNetSlice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []net.IPNet
		input          []string
		expectedErr    string
		expectedValues []net.IPNet
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []net.IPNet{},
			expectedErr:    "",
			expectedValues: []net.IPNet{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []net.IPNet{},
			expectedErr: `invalid argument "" for "--cidr" flag: invalid string being converted to CIDR: `,
		},
		{
			name:        "invalid ip",
			input:       []string{"blabla"},
			flagDefault: []net.IPNet{},
			expectedErr: `invalid argument "blabla" for "--cidr" flag: invalid string being converted to CIDR: blabla`,
		},
		{
			name:        "no csv",
			input:       []string{"192.168.1.1/16,172.16.1.1/16"},
			flagDefault: []net.IPNet{},
			expectedErr: `invalid argument "192.168.1.1/16,172.16.1.1/16" for "--cidr" flag: invalid string being converted to CIDR: 192.168.1.1/16,172.16.1.1/16`,
		},
		{
			name:           "empty defaults",
			input:          []string{"192.168.1.1/16", "10.0.0.1/16", "fd00::/64"},
			flagDefault:    []net.IPNet{},
			expectedValues: []net.IPNet{getCIDR("192.168.1.1/16"), getCIDR("10.0.0.1/16"), getCIDR("fd00::/64")},
		},
		{
			name:           "overrides default values",
			input:          []string{"192.168.1.1/16", "fd00::/64"},
			flagDefault:    []net.IPNet{getCIDR("fd00::/64"), getCIDR("192.168.1.1/16")},
			expectedValues: []net.IPNet{getCIDR("192.168.1.1/16"), getCIDR("fd00::/64")},
		},
		{
			name:           "with default values",
			input:          []string{},
			flagDefault:    []net.IPNet{getCIDR("192.168.1.1/16"), getCIDR("fd00::/64")},
			expectedValues: []net.IPNet{getCIDR("192.168.1.1/16"), getCIDR("fd00::/64")},
		},
		{
			name:  "sets values",
			input: []string{"192.168.1.1/16", "fd00::/64"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"192.168.1.2/24"})
				}
			},
			expectedValues: []net.IPNet{getCIDR("192.168.1.2/24")},
		},
		{
			name:           "trims input ipv4",
			input:          []string{"204.228.73.195/32", "86.141.15.94/32"},
			expectedValues: []net.IPNet{getCIDR("204.228.73.195/32"), getCIDR("86.141.15.94/32")},
		},
		{
			name:           "trims input ipv6",
			input:          []string{"fd00::/64", "        fd00::/64"},
			expectedValues: []net.IPNet{getCIDR("fd00::/64"), getCIDR("fd00::/64")},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var cidrs []net.IPNet
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.IPNetSliceVar(&cidrs, "cidr", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--cidr", test.input...))
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

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			if !reflect.DeepEqual(test.expectedValues, cidrs) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, cidrs)
			}

			getIPS, err := f.GetIPNetSlice("cidr")
			if err != nil {
				t.Fatal("got an error from GetIPNetSlice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, getIPS) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", getIPS, test.expectedValues)
			}

			getIPSGet, err := f.Get("cidr")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(getIPSGet, getIPS) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", getIPS, getIPSGet)
			}
		})
	}
}
