// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
	"strings"
)

// -- int64 Value
type int64Value int64

var _ Value = (*int64Value)(nil)
var _ Getter = (*int64Value)(nil)
var _ Typed = (*int64Value)(nil)

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(val string) error {
	val = strings.TrimSpace(val)
	v, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		return errors.New("must be an integer")
	}
	*i = int64Value(v)
	return nil
}

func (i *int64Value) Get() any {
	return int64(*i)
}

func (i *int64Value) Type() string {
	return "int64"
}

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

// GetInt64 return the int64 value of a flag with the given name
func (fs *FlagSet) GetInt64(name string) (int64, error) {
	val, err := fs.getFlagValue(name, "int64")
	if err != nil {
		return 0, err
	}
	return val.(int64), nil
}

// MustGetInt64 is like GetInt64, but panics on error.
func (fs *FlagSet) MustGetInt64(name string) int64 {
	val, err := fs.GetInt64(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func (fs *FlagSet) Int64Var(p *int64, name string, value int64, usage string, opts ...Opt) {
	fs.Var(newInt64Value(value, p), name, usage, opts...)
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func Int64Var(p *int64, name string, value int64, usage string, opts ...Opt) {
	CommandLine.Int64Var(p, name, value, usage, opts...)
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func (fs *FlagSet) Int64(name string, value int64, usage string, opts ...Opt) *int64 {
	var p int64
	fs.Int64Var(&p, name, value, usage, opts...)
	return &p
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func Int64(name string, value int64, usage string, opts ...Opt) *int64 {
	return CommandLine.Int64(name, value, usage, opts...)
}
