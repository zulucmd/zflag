// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/zulucmd/zflag/v2"
)

func parseTime(t *testing.T, value string) time.Time {
	t.Helper()

	res, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse time %s", value)
	}

	return res
}

func TestTime(t *testing.T) {
	tests := []struct {
		name          string
		flagDefault   time.Time
		input         []string
		formats       []string
		expectedErr   string
		expectedValue time.Time
	}{
		{
			name:          "no value passed",
			input:         []string{},
			flagDefault:   time.Time{},
			expectedErr:   "",
			expectedValue: time.Time{},
		},
		{
			name:        "empty value passed",
			input:       []string{""},
			flagDefault: time.Time{},
			expectedErr: `invalid argument "" for "--st" flag: invalid time format '' must be one of: '2006-01-02T15:04:05.999999999Z07:00'`,
		},
		{
			name:        "invalid datetime",
			input:       []string{"blabla"},
			flagDefault: time.Time{},
			expectedErr: `invalid argument "blabla" for "--st" flag: invalid time format 'blabla' must be one of: '2006-01-02T15:04:05.999999999Z07:00'`,
		},
		{
			name:        "no csv",
			input:       []string{"2022-01-01T01:01:01Z,2021-01-01T01:01:01Z"},
			flagDefault: time.Time{},
			expectedErr: `invalid argument "2022-01-01T01:01:01Z,2021-01-01T01:01:01Z" for "--st" flag: invalid time format '2022-01-01T01:01:01Z,2021-01-01T01:01:01Z' must be one of: '2006-01-02T15:04:05.999999999Z07:00'`,
		},
		{
			name:          "empty defaults",
			input:         []string{"2022-01-01T01:01:01Z"},
			flagDefault:   time.Time{},
			expectedValue: parseTime(t, "2022-01-01T01:01:01Z"),
		},
		{
			name:          "with default values",
			input:         []string{"2022-01-01T01:01:01Z"},
			flagDefault:   parseTime(t, "2021-01-01T01:01:01Z"),
			expectedValue: parseTime(t, "2022-01-01T01:01:01Z"),
		},
		{
			name:          "trims input",
			input:         []string{"  2022-01-01T01:01:01Z  "},
			flagDefault:   time.Time{},
			expectedValue: parseTime(t, "2022-01-01T01:01:01Z"),
		},
	}

	t.Parallel()
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var st time.Time
			f := zflag.NewFlagSet("test", zflag.ContinueOnError)
			f.SetOutput(io.Discard)
			formats := test.formats
			if len(formats) == 0 {
				formats = []string{time.RFC3339Nano}
			}
			f.TimeVar(&st, "st", test.flagDefault, formats, "usage")
			err := f.Parse(repeatFlag("--st", test.input...))
			if test.expectedErr != "" {
				assertErrMsg(t, test.expectedErr, err)
				return
			}
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, st)

			getTime, err := f.GetTime("st")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, getTime)

			getTimeGet, err := f.Get("st")
			assertNoErr(t, err)
			assertEqual(t, test.expectedValue, getTimeGet)

			defer assertNoPanic(t)()
			mustTime := f.MustGetTime("st")
			assertEqual(t, test.expectedValue, mustTime)
		})
	}
}

func usageForTimeFlagSet(t *testing.T, def time.Time) string {
	t.Helper()
	var ts time.Time
	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.TimeVar(&ts, "time", def, []string{time.RFC3339Nano, time.RFC1123Z}, "Time")
	assertNoErr(t, f.Parse(nil))

	return f.FlagUsages()
}

func TestTimeDefaultZero(t *testing.T) {
	usage := usageForTimeFlagSet(t, time.Time{})
	if strings.Contains(usage, "default") {
		t.Errorf("expected no default value in usage, got %q", usage)
	}
}

func TestTimeDefaultNonZero(t *testing.T) {
	usage := usageForTimeFlagSet(t, time.Date(2025, 1, 1, 1, 1, 1, 0, time.UTC))
	if !strings.Contains(usage, "default") {
		t.Errorf("expected default value in usage, got %q", usage)
	}
	if !strings.Contains(usage, "2025") {
		t.Errorf("expected default value in usage, got %q", usage)
	}
}
