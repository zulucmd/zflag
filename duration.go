// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strings"
	"time"
)

// -- time.Duration Value
type durationValue time.Duration

var _ Value = (*durationValue)(nil)
var _ Getter = (*durationValue)(nil)
var _ Typed = (*durationValue)(nil)

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(val string) error {
	val = strings.TrimSpace(val)
	v, err := time.ParseDuration(val)
	if err != nil {
		return errors.New(`must be a duration like "30s" or "5m"`)
	}
	*d = durationValue(v)
	return nil
}

func (d *durationValue) Get() interface{} {
	return time.Duration(*d)
}

func (d *durationValue) Type() string {
	return "duration"
}

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

// GetDuration return the duration value of a flag with the given name
func (fs *FlagSet) GetDuration(name string) (time.Duration, error) {
	val, err := fs.getFlagValue(name, "duration")
	if err != nil {
		return 0, err
	}
	return val.(time.Duration), nil
}

// MustGetDuration is like GetDuration, but panics on error.
func (fs *FlagSet) MustGetDuration(name string) time.Duration {
	val, err := fs.GetDuration(name)
	if err != nil {
		panic(err)
	}
	return val
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func (fs *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string, opts ...Opt) {
	fs.Var(newDurationValue(value, p), name, usage, opts...)
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func DurationVar(p *time.Duration, name string, value time.Duration, usage string, opts ...Opt) {
	CommandLine.DurationVar(p, name, value, usage, opts...)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func (fs *FlagSet) Duration(name string, value time.Duration, usage string, opts ...Opt) *time.Duration {
	var p time.Duration
	fs.DurationVar(&p, name, value, usage, opts...)
	return &p
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func Duration(name string, value time.Duration, usage string, opts ...Opt) *time.Duration {
	return CommandLine.Duration(name, value, usage, opts...)
}
