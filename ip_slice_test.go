// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"net"
	"reflect"
	"strings"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestIPSlice(t *testing.T) {
	tests := []struct {
		name           string
		flagDefault    []net.IP
		input          []string
		expectedErr    string
		expectedValues []net.IP
		visitor        func(f *zflag.Flag)
	}{
		{
			name:           "no value passed",
			input:          []string{},
			flagDefault:    []net.IP{},
			expectedErr:    "",
			expectedValues: []net.IP{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []net.IP{},
			expectedErr: `invalid argument "" for "--ips" flag: invalid string being converted to IP address`,
		},
		{
			name:        "invalid ip",
			input:       []string{"blabla"},
			flagDefault: []net.IP{},
			expectedErr: `invalid argument "blabla" for "--ips" flag: invalid string being converted to IP address`,
		},
		{
			name:        "no csv",
			input:       []string{"192.168.1.1,172.16.1.1"},
			flagDefault: []net.IP{},
			expectedErr: `invalid argument "192.168.1.1,172.16.1.1" for "--ips" flag: invalid string being converted to IP address`,
		},
		{
			name:           "empty value passed",
			input:          []string{"192.168.1.1", "10.0.0.1", "0:0:0:0:0:0:0:2"},
			flagDefault:    []net.IP{},
			expectedValues: []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("10.0.0.1"), net.ParseIP("0:0:0:0:0:0:0:2")},
		},
		{
			name:           "with default values",
			input:          []string{"192.168.1.1", "0:0:0:0:0:0:0:2"},
			flagDefault:    []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("0:0:0:0:0:0:0:1")},
			expectedValues: []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("0:0:0:0:0:0:0:2")},
		},
		{
			name:  "as slice values",
			input: []string{"192.168.1.1", "0:0:0:0:0:0:0:2"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"192.168.1.2"})
				}
			},
			expectedValues: []net.IP{net.ParseIP("192.168.1.2")},
		},
		{
			name:           "trims input ipv4",
			input:          []string{"204.228.73.195", "86.141.15.94"},
			expectedValues: []net.IP{net.ParseIP("204.228.73.195"), net.ParseIP("86.141.15.94")},
		},
		{
			name:           "trims input ipv6",
			input:          []string{"2e5e:66b2:6441:848:5b74:76ea:574c:3a7b", "        2e5e:66b2:6441:848:5b74:76ea:574c:3a7b", "2e5e:66b2:6441:848:5b74:76ea:574c:3a7b     ", "   2e5e:66b2:6441:848:5b74:76ea:574c:3a7b  "},
			expectedValues: []net.IP{net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"), net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"), net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b"), net.ParseIP("2e5e:66b2:6441:848:5b74:76ea:574c:3a7b")},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ips []net.IP
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.IPSliceVar(&ips, "ips", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ips", test.input...))
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

			if !reflect.DeepEqual(test.expectedValues, ips) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, ips)
			}

			getIPS, err := f.GetIPSlice("ips")
			if err != nil {
				t.Fatal("got an error from GetIPSlice():", err)
			}
			if !reflect.DeepEqual(test.expectedValues, getIPS) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, getIPS)
			}

			getIPSGet, err := f.Get("ips")
			if err != nil {
				t.Fatal("got an error from Get():", err)
			}
			if !reflect.DeepEqual(getIPSGet, getIPS) {
				t.Fatalf("expected %[1]v with type %[1]T but got %[2]v with type %[2]T", test.expectedValues, getIPS)
			}
		})
	}
}

func setUpIPSliceValue(t *testing.T, def []net.IP) (*[]net.IP, zflag.SliceValue) {
	t.Helper()
	var ips []net.IP
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.IPSliceVar(&ips, "ips", def, "usage")
	assertNoErr(t, f.Parse(nil))

	value, ok := f.Lookup("ips").Value.(zflag.SliceValue)
	if !ok {
		t.Fatal("flag value does not implement SliceValue")
	}
	return &ips, value
}

func TestIPSliceValueRejectsInvalidIPs(t *testing.T) {
	tests := []struct {
		name           string
		operation      func(zflag.SliceValue) error
		expectedErrMsg string
	}{
		{
			name:           "append",
			operation:      func(v zflag.SliceValue) error { return v.Append("not-an-ip") },
			expectedErrMsg: "invalid string being converted to IP address: not-an-ip",
		},
		{
			name:           "replace",
			operation:      func(v zflag.SliceValue) error { return v.Replace([]string{"2001:db8::1", "not-an-ip"}) },
			expectedErrMsg: "invalid string being converted to IP address: not-an-ip",
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			t.Run("returns error", func(t *testing.T) {
				t.Parallel()
				_, value := setUpIPSliceValue(t, []net.IP{net.ParseIP("192.0.2.1")})
				assertErrMsg(t, test.expectedErrMsg, test.operation(value))
			})

			t.Run("leaves value unchanged", func(t *testing.T) {
				t.Parallel()
				ips, value := setUpIPSliceValue(t, []net.IP{net.ParseIP("192.0.2.1")})
				_ = test.operation(value)
				assertDeepEqual(t, []net.IP{net.ParseIP("192.0.2.1")}, *ips)
			})
		})
	}
}

func TestIPSliceValueAppendAndReplace(t *testing.T) {
	ips, value := setUpIPSliceValue(t, nil)

	assertNoErr(t, value.Append(" 192.0.2.1 "))
	assertDeepEqual(t, []net.IP{net.ParseIP("192.0.2.1")}, *ips)

	assertNoErr(t, value.Replace([]string{" 2001:db8::1 "}))
	assertDeepEqual(t, []net.IP{net.ParseIP("2001:db8::1")}, *ips)
}
