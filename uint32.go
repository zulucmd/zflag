// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
	"strings"
)

// -- uint32 value
type uint32Value uint32

var _ Value = (*uint32Value)(nil)
var _ Getter = (*uint32Value)(nil)
var _ Typed = (*uint32Value)(nil)

func newUint32Value(val uint32, p *uint32) *uint32Value {
	*p = val
	return (*uint32Value)(p)
}

func (i *uint32Value) Set(val string) error {
	val = strings.TrimSpace(val)
	v, err := strconv.ParseUint(val, 0, 32)
	if err != nil {
		return errors.New("must be a non-negative integer")
	}
	*i = uint32Value(v)
	return nil
}

func (i *uint32Value) Get() interface{} {
	return uint32(*i)
}

func (i *uint32Value) Type() string {
	return "uint32"
}

func (i *uint32Value) String() string { return strconv.FormatUint(uint64(*i), 10) }

// GetUint32 return the uint32 value of a flag with the given name
func (fs *FlagSet) GetUint32(name string) (uint32, error) {
	val, err := fs.getFlagValue(name, "uint32")
	if err != nil {
		return 0, err
	}
	return val.(uint32), nil
}

// MustGetUint32 is like GetUint32, but panics on error.
func (fs *FlagSet) MustGetUint32(name string) uint32 {
	val, err := fs.GetUint32(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint32Var defines an uint32 flag with specified name, default value, and usage string.
// The argument p points to an uint32 variable in which to store the value of the flag.
func (fs *FlagSet) Uint32Var(p *uint32, name string, value uint32, usage string, opts ...Opt) {
	fs.Var(newUint32Value(value, p), name, usage, opts...)
}

// Uint32Var defines an uint32 flag with specified name, default value, and usage string.
// The argument p points to an uint32 variable in which to store the value of the flag.
func Uint32Var(p *uint32, name string, value uint32, usage string, opts ...Opt) {
	CommandLine.Uint32Var(p, name, value, usage, opts...)
}

// Uint32 defines an uint32 flag with specified name, default value, and usage string.
// The return value is the address of an uint32 variable that stores the value of the flag.
func (fs *FlagSet) Uint32(name string, value uint32, usage string, opts ...Opt) *uint32 {
	var p uint32
	fs.Uint32Var(&p, name, value, usage, opts...)
	return &p
}

// Uint32 defines an uint32 flag with specified name, default value, and usage string.
// The return value is the address of an uint32 variable that stores the value of the flag.
func Uint32(name string, value uint32, usage string, opts ...Opt) *uint32 {
	return CommandLine.Uint32(name, value, usage, opts...)
}
