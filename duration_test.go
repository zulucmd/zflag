// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"testing"
	"time"

	"github.com/zulucmd/zflag/v2"
)

func TestDuration(t *testing.T) {
	tests := []struct {
		name          string
		flagDefault   time.Duration
		input         []string
		expectedErr   string
		expectedValue time.Duration
		extraOpts     []zflag.Opt
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   time.Second,
			expectedErr:   "",
			expectedValue: time.Second,
		},
		{
			name:        "empty value passed",
			input:       repeatFlag("--dur", ""),
			expectedErr: `invalid argument "" for "--dur" flag: must be a duration like "30s" or "5m"`,
		},
		{
			name:        "invalid time.Duration",
			input:       repeatFlag("--dur", "blabla"),
			expectedErr: `invalid argument "blabla" for "--dur" flag: must be a duration like "30s" or "5m"`,
		},
		{
			name:        "no csv",
			input:       repeatFlag("--dur", "1s,2m"),
			expectedErr: `invalid argument "1s,2m" for "--dur" flag: must be a duration like "30s" or "5m"`,
		},
		{
			name:          "repeated value",
			input:         repeatFlag("--dur", "1s", "1m"),
			expectedValue: time.Minute,
		},
		{
			name:          "with default values",
			input:         repeatFlag("--dur", "4m"),
			expectedValue: time.Minute * 4,
		},
		{
			name:          "trims input true",
			input:         repeatFlag("--dur", " 1m "),
			expectedValue: time.Minute,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var dur time.Duration
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			f.DurationVar(&dur, "dur", test.flagDefault, "usage", test.extraOpts...)
			err := f.Parse(test.input)
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, dur)

			getBS, err := f.GetDuration("dur")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, getBS)

			getBSGet, err := f.Get("dur")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, getBSGet)

			defer assertNoPanic(t)()
			mustDuration := f.MustGetDuration("dur")
			assertEqual(t, test.expectedValue, mustDuration)
		})
	}
}

func TestDurationErrors(t *testing.T) {
	t.Parallel()

	var s string
	var dur time.Duration
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.DurationVar(&dur, "dur", time.Minute, "usage")
	err := f.Parse([]string{})
	assertNoErr(t, err)

	_, err = f.GetDuration("s")
	assertErr(t, err)

	defer assertPanic(t)()
	_ = f.MustGetDuration("s")
}
