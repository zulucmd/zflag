// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
	"strings"
)

// -- bool Value
type boolValue bool

var _ Value = (*boolValue)(nil)
var _ Getter = (*boolValue)(nil)
var _ Typed = (*boolValue)(nil)
var _ OptionalValue = (*boolValue)(nil)
var _ BoolFlag = (*boolValue)(nil)

func newBoolValue(val bool, p *bool) *boolValue {
	*p = val
	return (*boolValue)(p)
}

func (b *boolValue) Get() interface{} {
	return bool(*b)
}

func (b *boolValue) Set(val string) error {
	v := true
	if val != "" {
		val = strings.TrimSpace(val)
		var err error
		v, err = strconv.ParseBool(val)
		if err != nil {
			return errors.New("must be true or false")
		}
	}
	*b = boolValue(v)
	return nil
}

func (b *boolValue) Type() string {
	return "bool"
}

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) IsBoolFlag() bool { return true }

func (b *boolValue) IsOptional() bool { return true }

// GetBool return the bool value of a flag with the given name
func (fs *FlagSet) GetBool(name string) (bool, error) {
	val, err := fs.getFlagValue(name, "bool")
	if err != nil {
		return false, err
	}
	return val.(bool), nil
}

// MustGetBool is like GetBool, but panics on error.
func (fs *FlagSet) MustGetBool(name string) bool {
	val, err := fs.GetBool(name)
	if err != nil {
		panic(err)
	}
	return val
}

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (fs *FlagSet) BoolVar(p *bool, name string, value bool, usage string, opts ...Opt) {
	fs.Var(newBoolValue(value, p), name, usage, opts...)
}

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func BoolVar(p *bool, name string, value bool, usage string, opts ...Opt) {
	CommandLine.BoolVar(p, name, value, usage, opts...)
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func (fs *FlagSet) Bool(name string, value bool, usage string, opts ...Opt) *bool {
	var p bool
	fs.BoolVar(&p, name, value, usage, opts...)
	return &p
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func Bool(name string, value bool, usage string, opts ...Opt) *bool {
	return CommandLine.Bool(name, value, usage, opts...)
}
