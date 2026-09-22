// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
	"strings"
)

// -- float32 Value
type float32Value float32

var _ Value = (*float32Value)(nil)
var _ Getter = (*float32Value)(nil)
var _ Typed = (*float32Value)(nil)

func newFloat32Value(val float32, p *float32) *float32Value {
	*p = val
	return (*float32Value)(p)
}

func (f *float32Value) Set(val string) error {
	val = strings.TrimSpace(val)
	v, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return errors.New("must be a number")
	}
	*f = float32Value(v)
	return nil
}

func (f *float32Value) Get() any {
	return float32(*f)
}

func (f *float32Value) Type() string {
	return "float32"
}

func (f *float32Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 32) }

// GetFloat32 return the float32 value of a flag with the given name
func (fs *FlagSet) GetFloat32(name string) (float32, error) {
	val, err := fs.getFlagValue(name, "float32")
	if err != nil {
		return 0, err
	}
	return val.(float32), nil
}

// MustGetFloat32 is like GetFloat32, but panics on error.
func (fs *FlagSet) MustGetFloat32(name string) float32 {
	val, err := fs.GetFloat32(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float32 variable in which to store the value of the flag.
func (fs *FlagSet) Float32Var(p *float32, name string, value float32, usage string, opts ...Opt) {
	fs.Var(newFloat32Value(value, p), name, usage, opts...)
}

// Float32Var defines a float32 flag with specified name, default value, and usage string.
// The argument p points to a float32 variable in which to store the value of the flag.
func Float32Var(p *float32, name string, value float32, usage string, opts ...Opt) {
	CommandLine.Float32Var(p, name, value, usage, opts...)
}

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func (fs *FlagSet) Float32(name string, value float32, usage string, opts ...Opt) *float32 {
	var p float32
	fs.Float32Var(&p, name, value, usage, opts...)
	return &p
}

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func Float32(name string, value float32, usage string, opts ...Opt) *float32 {
	return CommandLine.Float32(name, value, usage, opts...)
}
