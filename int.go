// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
	"strings"
)

// -- int Value
type intValue int

var _ Value = (*intValue)(nil)
var _ Getter = (*intValue)(nil)
var _ Typed = (*intValue)(nil)

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(val string) error {
	val = strings.TrimSpace(val)
	v, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		return errors.New("must be an integer")
	}
	*i = intValue(v)
	return nil
}

func (i *intValue) Get() interface{} {
	return int(*i)
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

// GetInt return the int value of a flag with the given name
func (fs *FlagSet) GetInt(name string) (int, error) {
	val, err := fs.getFlagValue(name, "int")
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

// MustGetInt is like GetInt, but panics on error.
func (fs *FlagSet) MustGetInt(name string) int {
	val, err := fs.GetInt(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func (fs *FlagSet) IntVar(p *int, name string, value int, usage string, opts ...Opt) {
	fs.Var(newIntValue(value, p), name, usage, opts...)
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string, opts ...Opt) {
	CommandLine.IntVar(p, name, value, usage, opts...)
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (fs *FlagSet) Int(name string, value int, usage string, opts ...Opt) *int {
	var p int
	fs.IntVar(&p, name, value, usage, opts...)
	return &p
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, value int, usage string, opts ...Opt) *int {
	return CommandLine.Int(name, value, usage, opts...)
}
