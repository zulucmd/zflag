// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestBoolSlice(t *testing.T) {
	tests := []struct {
		name              string
		flagDefault       []bool
		input             []string
		expectedErr       string
		expectedValues    []bool
		expectedStrValues string
		visitor           func(f *zflag.Flag)
		expectedGetSlice  []string
	}{
		{
			name:              "no value passed",
			input:             []string{},
			flagDefault:       []bool{},
			expectedErr:       "",
			expectedValues:    []bool{},
			expectedStrValues: "[]",
			expectedGetSlice:  []string{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []bool{},
			expectedErr: `invalid argument "" for "--bs" flag: must be true or false`,
		},
		{
			name:        "invalid bool",
			input:       []string{"blabla"},
			flagDefault: []bool{},
			expectedErr: `invalid argument "blabla" for "--bs" flag: must be true or false`,
		},
		{
			name:        "no csv",
			input:       []string{"true,false"},
			flagDefault: []bool{},
			expectedErr: `invalid argument "true,false" for "--bs" flag: must be true or false`,
		},
		{
			name:              "multiple values passed",
			input:             []string{"true", "false"},
			flagDefault:       []bool{},
			expectedValues:    []bool{true, false},
			expectedStrValues: "[true false]",
			expectedGetSlice:  []string{"true", "false"},
		},
		{
			name:              "with default values",
			input:             []string{"false", "true"},
			flagDefault:       []bool{true, false},
			expectedValues:    []bool{false, true},
			expectedStrValues: "[false true]",
			expectedGetSlice:  []string{"false", "true"},
		},
		{
			name:  "replace values",
			input: []string{"true", "false"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"false"})
				}
			},
			expectedValues:    []bool{false},
			expectedStrValues: "[false]",
			expectedGetSlice:  []string{"false"},
		},
		{
			name:  "replace values error",
			input: []string{"true", "false"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					err := val.Replace([]string{"notbool"})
					assertErr(t, err)
				}
			},
			expectedValues:    []bool{true, false},
			expectedStrValues: "[true false]",
			expectedGetSlice:  []string{"true", "false"},
		},
		{
			name:  "add values",
			input: []string{"true", "false"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Append("false")
				}
			},
			expectedValues:    []bool{true, false, false},
			expectedStrValues: "[true false false]",
			expectedGetSlice:  []string{"true", "false", "false"},
		},
		{
			name:  "add values error",
			input: []string{"true"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					err := val.Append("asd")
					if err == nil {
						t.Errorf("Expected an error when appending, got %s", err)
					}
				}
			},
			flagDefault:       nil,
			expectedValues:    []bool{true},
			expectedStrValues: "[true]",
			expectedGetSlice:  []string{"true"},
		},
		{
			name:              "nil default",
			input:             []string{},
			flagDefault:       nil,
			expectedValues:    nil,
			expectedStrValues: "[]",
			expectedGetSlice:  []string{},
		},
		{
			name:              "trims input",
			input:             []string{" true ", " false "},
			expectedValues:    []bool{true, false},
			expectedStrValues: "[true false]",
			expectedGetSlice:  []string{"true", "false"},
		},
		{
			name:              "all valid bool values",
			input:             []string{"true", "false", "1", "0", "t", "f", "TRUE", "FALSE", "1", "0", "T", "F", "True", "False"},
			expectedValues:    []bool{true, false, true, false, true, false, true, false, true, false, true, false, true, false},
			expectedStrValues: "[true false true false true false true false true false true false true false]",
			expectedGetSlice:  []string{"true", "false", "true", "false", "true", "false", "true", "false", "true", "false", "true", "false", "true", "false"},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var bs []bool
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.BoolSliceVar(&bs, "bs", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--bs", test.input...))
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			assertDeepEqual(t, test.expectedValues, bs)

			boolSlice, err := f.GetBoolSlice("bs")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, boolSlice)

			boolSliceGet, err := f.Get("bs")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, boolSliceGet)

			flag := f.Lookup("bs")
			assertEqual(t, test.expectedStrValues, flag.Value.String())

			sliced := flag.Value.(zflag.SliceValue)
			assertDeepEqual(t, test.expectedGetSlice, sliced.GetSlice())

			defer assertNoPanic(t)()
			mustBoolSlice := f.MustGetBoolSlice("bs")
			assertDeepEqual(t, test.expectedValues, mustBoolSlice)
		})
	}
}

func TestBoolSliceErrors(t *testing.T) {
	t.Parallel()

	var s string
	var bs []bool
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.BoolSliceVar(&bs, "bs", []bool{}, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetBoolSlice("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetBoolSlice("s")
}
