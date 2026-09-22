// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/zulucmd/zflag/v2"
)

func TestDurationSlice(t *testing.T) {
	tests := []struct {
		name              string
		flagDefault       []time.Duration
		input             []string
		expectedErr       string
		expectedValues    []time.Duration
		expectedStrValues string
		expectedGetSlice  []string
		visitor           func(f *zflag.Flag)
	}{
		{
			name:              "no value passed",
			input:             []string{},
			flagDefault:       []time.Duration{},
			expectedErr:       "",
			expectedValues:    []time.Duration{},
			expectedStrValues: "[]",
			expectedGetSlice:  []string{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: []time.Duration{},
			expectedErr: `invalid argument "" for "--ds" flag: must be a duration like "30s" or "5m"`,
		},
		{
			name:        "invalid unit",
			input:       []string{"2q"},
			flagDefault: []time.Duration{},
			expectedErr: `invalid argument "2q" for "--ds" flag: must be a duration like "30s" or "5m"`,
		},
		{
			name:        "invalid duration",
			input:       []string{"blabla"},
			flagDefault: []time.Duration{},
			expectedErr: `invalid argument "blabla" for "--ds" flag: must be a duration like "30s" or "5m"`,
		},
		{
			name:        "no csv",
			input:       []string{"2m,2h"},
			flagDefault: []time.Duration{},
			expectedErr: `invalid argument "2m,2h" for "--ds" flag: must be a duration like "30s" or "5m"`,
		},
		{
			name:              "defaults returned",
			input:             []string{},
			flagDefault:       []time.Duration{time.Second, time.Minute},
			expectedValues:    []time.Duration{time.Second, time.Minute},
			expectedStrValues: "[1s 1m0s]",
			expectedGetSlice:  []string{"1s", "1m0s"},
		},
		{
			name:              "defaults overwritten",
			input:             []string{"1m", "1s"},
			flagDefault:       []time.Duration{time.Second, time.Minute},
			expectedValues:    []time.Duration{time.Minute, time.Second},
			expectedStrValues: "[1m0s 1s]",
			expectedGetSlice:  []string{"1m0s", "1s"},
		},
		{
			name:  "replace values",
			input: []string{"1m", "1h"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Replace([]string{"1s"})
				}
			},
			expectedValues:    []time.Duration{time.Second},
			expectedStrValues: "[1s]",
			expectedGetSlice:  []string{"1s"},
		},
		{
			name:  "replace values error",
			input: []string{"5m"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					err := val.Replace([]string{"notduration"})
					assertErr(t, err)
				}
			},
			expectedValues:    []time.Duration{time.Minute * 5},
			expectedStrValues: "[5m0s]",
			expectedGetSlice:  []string{"5m0s"},
		},
		{
			name:  "add values",
			input: []string{"1s"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					_ = val.Append("1m")
				}
			},
			expectedValues:    []time.Duration{time.Second, time.Minute},
			expectedStrValues: "[1s 1m0s]",
			expectedGetSlice:  []string{"1s", "1m0s"},
		},
		{
			name:  "add values error",
			input: []string{"1m"},
			visitor: func(f *zflag.Flag) {
				if val, ok := f.Value.(zflag.SliceValue); ok {
					err := val.Append("asd")
					assertErr(t, err)
				}
			},
			flagDefault:       nil,
			expectedValues:    []time.Duration{time.Minute},
			expectedStrValues: "[1m0s]",
			expectedGetSlice:  []string{"1m0s"},
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
			input:             []string{" 1ns", "2ms  ", "  3m    ", "    4h"},
			expectedValues:    []time.Duration{1 * time.Nanosecond, 2 * time.Millisecond, 3 * time.Minute, 4 * time.Hour},
			expectedStrValues: "[1ns 2ms 3m0s 4h0m0s]",
			expectedGetSlice:  []string{"1ns", "2ms", "3m0s", "4h0m0s"},
		},
		{
			name:              "valid values",
			input:             []string{"1ns", "2ms", "3m", "4h"},
			expectedValues:    []time.Duration{1 * time.Nanosecond, 2 * time.Millisecond, 3 * time.Minute, 4 * time.Hour},
			expectedStrValues: "[1ns 2ms 3m0s 4h0m0s]",
			expectedGetSlice:  []string{"1ns", "2ms", "3m0s", "4h0m0s"},
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var ds []time.Duration
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(ioutil.Discard)
			f.DurationSliceVar(&ds, "ds", test.flagDefault, "usage")
			err := f.Parse(repeatFlag("--ds", test.input...))
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)

			if test.visitor != nil {
				f.VisitAll(test.visitor)
			}

			assertDeepEqual(t, test.expectedValues, ds)

			durSlice, err := f.GetDurationSlice("ds")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, durSlice)

			durSliceGet, err := f.Get("ds")
			assertNoErr(t, err)
			assertDeepEqual(t, test.expectedValues, durSliceGet)

			flag := f.Lookup("ds")
			assertEqual(t, test.expectedStrValues, flag.Value.String())

			sliced := flag.Value.(zflag.SliceValue)
			assertDeepEqual(t, test.expectedGetSlice, sliced.GetSlice())

			defer assertNoPanic(t)()
			mustDurSlice := f.MustGetDurationSlice("ds")
			assertDeepEqual(t, test.expectedValues, mustDurSlice)
		})
	}
}

func TestDurationSliceErrors(t *testing.T) {
	t.Parallel()

	var s string
	var bs []time.Duration
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.DurationSliceVar(&bs, "bs", []time.Duration{}, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetDurationSlice("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetDurationSlice("s")
}
