// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
	"strings"
)

// -- uint Value
type uintValue uint

var _ Value = (*uintValue)(nil)
var _ Getter = (*uintValue)(nil)
var _ Typed = (*uintValue)(nil)

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (i *uintValue) Set(val string) error {
	val = strings.TrimSpace(val)
	v, err := strconv.ParseUint(val, 0, 0)
	if err != nil {
		return errors.New("must be a non-negative integer")
	}
	*i = uintValue(v)
	return nil
}

func (i *uintValue) Get() interface{} {
	return uint(*i)
}

func (i *uintValue) Type() string {
	return "uint"
}

func (i *uintValue) String() string { return strconv.FormatUint(uint64(*i), 10) }

// GetUint return the uint value of a flag with the given name
func (fs *FlagSet) GetUint(name string) (uint, error) {
	val, err := fs.getFlagValue(name, "uint")
	if err != nil {
		return 0, err
	}
	return val.(uint), nil
}

// MustGetUint is like GetUint, but panics on error.
func (fs *FlagSet) MustGetUint(name string) uint {
	val, err := fs.GetUint(name)
	if err != nil {
		panic(err)
	}
	return val
}

// UintVar defines an uint flag with specified name, default value, and usage string.
// The argument p points to an uint variable in which to store the value of the flag.
func (fs *FlagSet) UintVar(p *uint, name string, value uint, usage string, opts ...Opt) {
	fs.Var(newUintValue(value, p), name, usage, opts...)
}

// UintVar defines an uint flag with specified name, default value, and usage string.
// The argument p points to an uint variable in which to store the value of the flag.
func UintVar(p *uint, name string, value uint, usage string, opts ...Opt) {
	CommandLine.UintVar(p, name, value, usage, opts...)
}

// Uint defines an uint flag with specified name, default value, and usage string.
// The return value is the address of an uint variable that stores the value of the flag.
func (fs *FlagSet) Uint(name string, value uint, usage string, opts ...Opt) *uint {
	var p uint
	fs.UintVar(&p, name, value, usage, opts...)
	return &p
}

// Uint defines an uint flag with specified name, default value, and usage string.
// The return value is the address of an uint variable that stores the value of the flag.
func Uint(name string, value uint, usage string, opts ...Opt) *uint {
	return CommandLine.Uint(name, value, usage, opts...)
}
