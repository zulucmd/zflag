// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
	"strings"
)

// -- float64 Value
type float64Value float64

var _ Value = (*float64Value)(nil)
var _ Getter = (*float64Value)(nil)
var _ Typed = (*float64Value)(nil)

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(val string) error {
	val = strings.TrimSpace(val)
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return errors.New("must be a number")
	}
	*f = float64Value(v)
	return nil
}

func (f *float64Value) Get() any {
	return float64(*f)
}

func (f *float64Value) Type() string {
	return "float64"
}

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

// GetFloat64 return the float64 value of a flag with the given name
func (fs *FlagSet) GetFloat64(name string) (float64, error) {
	val, err := fs.getFlagValue(name, "float64")
	if err != nil {
		return 0, err
	}
	return val.(float64), nil
}

// MustGetFloat64 is like GetFloat64, but panics on error.
func (fs *FlagSet) MustGetFloat64(name string) float64 {
	val, err := fs.GetFloat64(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func (fs *FlagSet) Float64Var(p *float64, name string, value float64, usage string, opts ...Opt) {
	fs.Var(newFloat64Value(value, p), name, usage, opts...)
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func Float64Var(p *float64, name string, value float64, usage string, opts ...Opt) {
	CommandLine.Float64Var(p, name, value, usage, opts...)
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func (fs *FlagSet) Float64(name string, value float64, usage string, opts ...Opt) *float64 {
	var p float64
	fs.Float64Var(&p, name, value, usage, opts...)
	return &p
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func Float64(name string, value float64, usage string, opts ...Opt) *float64 {
	return CommandLine.Float64(name, value, usage, opts...)
}
