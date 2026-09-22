// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strings"
	"time"
)

// TimeValue adapts time.Time for use as a flag.
type TimeValue struct {
	*time.Time
	formats []string
}

var _ Value = (*TimeValue)(nil)
var _ Getter = (*TimeValue)(nil)
var _ Typed = (*TimeValue)(nil)

func newTimeValue(val time.Time, p *time.Time, formats []string) *TimeValue {
	*p = val
	return &TimeValue{
		Time:    p,
		formats: formats,
	}
}

// Set time.Time value from string based on accepted formats.
func (d *TimeValue) Set(s string) error {
	s = strings.TrimSpace(s)
	for _, f := range d.formats {
		v, err := time.Parse(f, s)
		if err != nil {
			continue
		}
		*d.Time = v
		return nil
	}

	formatsString := "'" + strings.Join(d.formats, "', '") + "'"
	return fmt.Errorf("invalid time format '%s' must be one of: %s", s, formatsString)
}

// Get returns the current time value.
func (d *TimeValue) Get() any {
	return *d.Time
}

// Type name for time.Time flags.
func (d *TimeValue) Type() string {
	return "time"
}

// String returns the time formatted as RFC3339Nano, or an empty string if the
// time is zero.
func (d *TimeValue) String() string {
	if d.Time.IsZero() {
		return ""
	}
	return d.Time.Format(time.RFC3339Nano)
}

// GetTime return the time value of a flag with the given name
func (fs *FlagSet) GetTime(name string) (time.Time, error) {
	val, err := fs.getFlagValue(name, "time")
	if err != nil {
		return time.Time{}, err
	}
	return val.(time.Time), nil
}

// MustGetTime is like GetTime, but panics on error.
func (fs *FlagSet) MustGetTime(name string) time.Time {
	val, err := fs.GetTime(name)
	if err != nil {
		panic(err)
	}
	return val
}

// TimeVar defines a time.Time flag with specified name, default value, and usage string.
// The argument p points to a time.Time variable in which to store the value of the flag.
func (fs *FlagSet) TimeVar(p *time.Time, name string, value time.Time, formats []string, usage string, opts ...Opt) {
	fs.Var(newTimeValue(value, p, formats), name, usage, opts...)
}

// TimeVar defines a time.Time flag with specified name, default value, and usage string.
// The argument p points to a time.Time variable in which to store the value of the flag.
func TimeVar(p *time.Time, name string, value time.Time, formats []string, usage string, opts ...Opt) {
	CommandLine.Var(newTimeValue(value, p, formats), name, usage, opts...)
}

// Time defines a time.Time flag with specified name, default value, and usage string.
// The return value is the address of a time.Time variable that stores the value of the flag.
func (fs *FlagSet) Time(name string, value time.Time, formats []string, usage string, opts ...Opt) *time.Time {
	p := new(time.Time)
	fs.TimeVar(p, name, value, formats, usage, opts...)
	return p
}

// Time is like Time, but accepts a shorthand letter that can be used after a single dash.
func Time(name string, value time.Time, formats []string, usage string, opts ...Opt) *time.Time {
	return CommandLine.Time(name, value, formats, usage, opts...)
}
