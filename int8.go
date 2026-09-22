// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
	"strings"
)

// -- int8 Value
type int8Value int8

var _ Value = (*int8Value)(nil)
var _ Getter = (*int8Value)(nil)
var _ Typed = (*int8Value)(nil)

func newInt8Value(val int8, p *int8) *int8Value {
	*p = val
	return (*int8Value)(p)
}

func (i *int8Value) Set(val string) error {
	val = strings.TrimSpace(val)
	v, err := strconv.ParseInt(val, 0, 8)
	if err != nil {
		return errors.New("must be an integer")
	}
	*i = int8Value(v)
	return nil
}

func (i *int8Value) Get() interface{} {
	return int8(*i)
}

func (i *int8Value) Type() string {
	return "int8"
}

func (i *int8Value) String() string { return strconv.FormatInt(int64(*i), 10) }

// GetInt8 return the int8 value of a flag with the given name
func (fs *FlagSet) GetInt8(name string) (int8, error) {
	val, err := fs.getFlagValue(name, "int8")
	if err != nil {
		return 0, err
	}
	return val.(int8), nil
}

// MustGetInt8 is like GetInt8, but panics on error.
func (fs *FlagSet) MustGetInt8(name string) int8 {
	val, err := fs.GetInt8(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int8Var defines an int8 flag with specified name, default value, and usage string.
// The argument p points to an int8 variable in which to store the value of the flag.
func (fs *FlagSet) Int8Var(p *int8, name string, value int8, usage string, opts ...Opt) {
	fs.Var(newInt8Value(value, p), name, usage, opts...)
}

// Int8Var defines an int8 flag with specified name, default value, and usage string.
// The argument p points to an int8 variable in which to store the value of the flag.
func Int8Var(p *int8, name string, value int8, usage string, opts ...Opt) {
	CommandLine.Int8Var(p, name, value, usage, opts...)
}

// Int8 defines an int8 flag with specified name, default value, and usage string.
// The return value is the address of an int8 variable that stores the value of the flag.
func (fs *FlagSet) Int8(name string, value int8, usage string, opts ...Opt) *int8 {
	var p int8
	fs.Int8Var(&p, name, value, usage, opts...)
	return &p
}

// Int8 defines an int8 flag with specified name, default value, and usage string.
// The return value is the address of an int8 variable that stores the value of the flag.
func Int8(name string, value int8, usage string, opts ...Opt) *int8 {
	return CommandLine.Int8(name, value, usage, opts...)
}
