// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/zulucmd/zflag/v2"
)

var testDefaultTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

func assertTimeEqual(t *testing.T, expected, actual time.Time) {
	t.Helper()
	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func newTextFlagSet() (*zflag.FlagSet, *time.Time) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	ts := new(time.Time)
	f.TextVar(ts, "time", testDefaultTime, "time stamp")
	return f, ts
}

func TestText(t *testing.T) {
	tests := []struct {
		name        string
		input       []string
		expectedErr string
		expected    time.Time
	}{
		{
			name:     "rfc3339",
			input:    []string{"--time=2003-01-02T15:04:05Z"},
			expected: time.Date(2003, 1, 2, 15, 4, 5, 0, time.UTC),
		},
		{
			name:        "invalid layout",
			input:       []string{"--time=2003-01-02 15:05:01"},
			expectedErr: `invalid argument "2003-01-02 15:05:01" for "--time" flag: parsing time "2003-01-02 15:05:01" as "2006-01-02T15:04:05Z07:00": cannot parse " 15:05:01" as "T"`,
		},
		{
			name:     "numeric zone",
			input:    []string{"--time=2006-01-02T15:04:05+07:00"},
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("UTC+7", 7*60*60)),
		},
		{
			name:     "default kept when flag absent",
			input:    nil,
			expected: testDefaultTime,
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			f, ts := newTextFlagSet()

			err := f.Parse(test.input)
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)
			assertTimeEqual(t, test.expected, *ts)

			got := new(time.Time)
			assertNoErr(t, f.GetText("time", got))
			assertTimeEqual(t, test.expected, *got)

			gotGet, err := f.Get("time")
			assertNoErr(t, err)
			assertTimeEqual(t, test.expected, *gotGet.(*time.Time))

			defer assertNoPanic(t)()
			f.MustGetText("time", got)
			assertTimeEqual(t, test.expected, *got)
		})
	}
}

func TestTextConvenience(t *testing.T) {
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	when := f.Text("time", testDefaultTime, "usage").(*time.Time)
	assertTimeEqual(t, testDefaultTime, *when)

	assertNoErr(t, f.Parse([]string{"--time=2003-01-02T15:04:05Z"}))
	assertTimeEqual(t, time.Date(2003, 1, 2, 15, 4, 5, 0, time.UTC), *when)
}

func TestTextErrors(t *testing.T) {
	var s string
	var ts time.Time
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.SetOutput(ioutil.Discard)
	f.StringVar(&s, "s", "", "usage")
	f.TextVar(&ts, "time", testDefaultTime, "usage")
	assertNoErr(t, f.Parse(nil))

	assertErr(t, f.GetText("s", &ts))
	assertErr(t, f.GetText("missing", &ts))

	defer assertPanic(t)()
	f.MustGetText("s", &ts)
}
