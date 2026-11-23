// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"strconv"
)

// -- count Value
type countValue int

var _ Value = (*countValue)(nil)
var _ Getter = (*countValue)(nil)
var _ Typed = (*countValue)(nil)
var _ OptionalValue = (*countValue)(nil)

func newCountValue(val int, p *int) *countValue {
	*p = val
	return (*countValue)(p)
}

func (i *countValue) Set(val string) error {
	if val == "" {
		*i++
		return nil
	}

	v, err := strconv.ParseInt(val, 0, 0)
	if err != nil {
		return errors.New("must be an integer")
	}
	*i = countValue(v)

	return nil
}

func (i *countValue) Get() any {
	return int(*i)
}

func (i *countValue) Type() string {
	return "count"
}

func (i *countValue) String() string { return strconv.Itoa(int(*i)) }

func (i *countValue) IsOptional() bool { return true }

// GetCount return the int value of a flag with the given name
func (fs *FlagSet) GetCount(name string) (int, error) {
	val, err := fs.getFlagValue(name, "count")
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

// MustGetCount is like GetCount, but panics on error.
func (fs *FlagSet) MustGetCount(name string) int {
	val, err := fs.GetCount(name)
	if err != nil {
		panic(err)
	}
	return val
}

// CountVar defines a count flag with specified name, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
// A count flag will add 1 to its value every time it is found on the command line
func (fs *FlagSet) CountVar(p *int, name string, usage string, opts ...Opt) {
	fs.Var(newCountValue(0, p), name, usage, opts...)
}

// CountVar like CountVar only the flag is placed on the CommandLine instead of a given flag set
func CountVar(p *int, name string, usage string, opts ...Opt) {
	CommandLine.CountVar(p, name, usage, opts...)
}

// Count defines a count flag with specified name, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
// A count flag will add 1 to its value every time it is found on the command line
func (fs *FlagSet) Count(name string, usage string, opts ...Opt) *int {
	var p int
	fs.CountVar(&p, name, usage, opts...)
	return &p
}

// Count defines a count flag with specified name, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
// A count flag will add 1 to its value every time it is found on the command line
func Count(name string, usage string, opts ...Opt) *int {
	return CommandLine.Count(name, usage, opts...)
}
