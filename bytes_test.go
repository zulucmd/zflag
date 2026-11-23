// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestBytesHex(t *testing.T) {
	tests := []struct {
		name          string
		flagDefault   []byte
		input         []string
		expectedErr   string
		expectedValue string
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   []byte{},
			expectedErr:   "",
			expectedValue: "",
		},
		{
			name:          "empty value passed",
			input:         []string{""},
			flagDefault:   []byte{},
			expectedErr:   "",
			expectedValue: "",
		},
		{
			name:        "invalid byte hex short string",
			input:       []string{"0"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "0" for "--bytes" flag: encoding/hex: odd length hex string`,
		},
		{
			name:        "invalid byte hex odd-length string",
			input:       []string{"000"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "000" for "--bytes" flag: encoding/hex: odd length hex string`,
		},
		{
			name:        "invalid byte hex non-hex char",
			input:       []string{"qq"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "qq" for "--bytes" flag: encoding/hex: invalid byte: U+0071 'q'`,
		},
		{
			name:        "no csv",
			input:       []string{"0101,0101"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "0101,0101" for "--bytes" flag: encoding/hex: invalid byte: U+002C ','`,
		},
		{
			name:          "repeated value gets the last call",
			input:         []string{"01", "0101"},
			flagDefault:   []byte{},
			expectedValue: "0101",
		},
		{
			name:          "default values get overwritten",
			input:         []string{"0101"},
			flagDefault:   []byte("01"),
			expectedValue: "0101",
		},
		{
			name:          "trims input",
			input:         []string{" 01 "},
			expectedValue: "01",
		},
		{
			name:          "test valid",
			input:         []string{"01"},
			expectedValue: "01",
		},
		{
			name:          "test valid",
			input:         []string{"0101"},
			expectedValue: "0101",
		},
		{
			name:          "test valid",
			input:         []string{"1234567890abcdef"},
			expectedValue: "1234567890ABCDEF",
		},
		{
			name:          "test valid",
			input:         []string{"1234567890ABCDEF"},
			expectedValue: "1234567890ABCDEF",
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var bytes []byte
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.BytesHexVar(&bytes, "bytes", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--bytes", test.input...))
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, fmt.Sprintf("%X", bytes))

			bytesHex, err := f.GetBytesHex("bytes")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, fmt.Sprintf("%X", bytesHex))

			bytesHexGet, err := f.Get("bytes")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, fmt.Sprintf("%X", bytesHexGet))

			flag := f.Lookup("bytes")
			assertEqual(t, test.expectedValue, flag.Value.String())

			defer assertNoPanic(t)()
			mustBytesHex := f.MustGetBytesHex("bytes")
			assertDeepEqual(t, test.expectedValue, fmt.Sprintf("%X", mustBytesHex))
		})
	}
}

func TestBytesHexErrors(t *testing.T) {
	var s string
	var b []byte
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.BytesHexVar(&b, "b", []byte{}, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetBytesHex("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetBytesHex("s")
}

func TestBytesB64(t *testing.T) {
	tests := []struct {
		name             string
		flagDefault      []byte
		input            []string
		expectedErr      string
		expectedValue    []byte
		expectedStrValue string
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   []byte{},
			expectedValue: []byte{},
		},
		{
			name:          "empty value passed",
			input:         []string{""},
			flagDefault:   []byte{},
			expectedValue: []byte{},
		},
		{
			name:        "invalid byte base64 short string",
			input:       []string{"A"},
			flagDefault: []byte{},
			expectedErr: `invalid argument "A" for "--bytes" flag: illegal base64 data at input byte 0`,
		},
		{
			name:        "invalid byte hex non-hex char",
			input:       []string{"Aï=="},
			flagDefault: []byte{},
			expectedErr: `invalid argument "Aï==" for "--bytes" flag: illegal base64 data at input byte 1`,
		},
		{
			name:        "no csv",
			input:       []string{"AQ==,AQ=="},
			flagDefault: []byte{},
			expectedErr: `invalid argument "AQ==,AQ==" for "--bytes" flag: illegal base64 data at input byte 4`,
		},
		{
			name:             "default value gets returned when none passed",
			input:            []string{},
			flagDefault:      []byte("bye"),
			expectedValue:    []byte("bye"),
			expectedStrValue: "Ynll",
		},
		{
			name:             "repeated value gets the last call",
			input:            []string{"aGk=", "Ynll"},
			flagDefault:      []byte{},
			expectedValue:    []byte("bye"),
			expectedStrValue: "Ynll",
		},
		{
			name:             "default values get overwritten",
			input:            []string{"Ynll"},
			flagDefault:      []byte("aGk="),
			expectedValue:    []byte("bye"),
			expectedStrValue: "Ynll",
		},
		{
			name:             "trims input",
			input:            []string{" Ynll "},
			expectedValue:    []byte("bye"),
			expectedStrValue: "Ynll",
		},
		{
			name:             "test valid",
			input:            []string{"Ynll"},
			expectedValue:    []byte("bye"),
			expectedStrValue: "Ynll",
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var bytes []byte
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.BytesBase64Var(&bytes, "bytes", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--bytes", test.input...))
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValue, bytes)

			bytesB64, err := f.GetBytesBase64("bytes")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValue, bytesB64)

			bytesB64Get, err := f.Get("bytes")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValue, bytesB64Get)

			flag := f.Lookup("bytes")
			assertEqual(t, test.expectedStrValue, flag.Value.String())

			defer assertNoPanic(t)()
			mustBytesBase64 := f.MustGetBytesBase64("bytes")
			assertDeepEqual(t, test.expectedValue, mustBytesBase64)
		})
	}
}

func TestBytesBase64Errors(t *testing.T) {
	t.Parallel()

	var s bool
	var b []byte
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.BoolVar(&s, "s", false, "usage")
	f.BytesBase64Var(&b, "b", []byte{}, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetBytesBase64("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetBytesBase64("s")
}
