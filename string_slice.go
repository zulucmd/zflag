// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
)

// -- stringSlice Value
type stringSliceValue struct {
	value   *[]string
	changed bool
}

var _ Value = (*stringSliceValue)(nil)
var _ Getter = (*stringSliceValue)(nil)
var _ SliceValue = (*stringSliceValue)(nil)
var _ Typed = (*stringSliceValue)(nil)

func newStringSliceValue(val []string, p *[]string) *stringSliceValue {
	ssv := new(stringSliceValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringSliceValue) Set(val string) error {
	if !s.changed {
		*s.value = []string{}
	}
	*s.value = append(*s.value, val)
	s.changed = true

	return nil
}

func (s *stringSliceValue) Get() any {
	return *s.value
}

func (s *stringSliceValue) Type() string {
	return "stringSlice"
}

func (s *stringSliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%s", *s.value)
}

func (s *stringSliceValue) Append(val string) error {
	*s.value = append(*s.value, val)
	return nil
}

func (s *stringSliceValue) Replace(val []string) error {
	*s.value = val
	return nil
}

func (s *stringSliceValue) GetSlice() []string {
	return *s.value
}

// GetStringSlice return the []string value of a flag with the given name
func (fs *FlagSet) GetStringSlice(name string) ([]string, error) {
	val, err := fs.getFlagValue(name, "stringSlice")
	if err != nil {
		return []string{}, err
	}
	return val.([]string), nil
}

// MustGetStringSlice is like GetStringSlice, but panics on error.
func (fs *FlagSet) MustGetStringSlice(name string) []string {
	val, err := fs.GetStringSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// StringSliceVar defines a []string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
func (fs *FlagSet) StringSliceVar(p *[]string, name string, value []string, usage string, opts ...Opt) {
	fs.Var(newStringSliceValue(value, p), name, usage, opts...)
}

// StringSliceVar defines a []string flag with specified name, default value, and usage string.
// The argument p points to a []string variable in which to store the value of the flag.
func StringSliceVar(p *[]string, name string, value []string, usage string, opts ...Opt) {
	CommandLine.StringSliceVar(p, name, value, usage, opts...)
}

// StringSlice defines a []string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
func (fs *FlagSet) StringSlice(name string, value []string, usage string, opts ...Opt) *[]string {
	var p []string
	fs.StringSliceVar(&p, name, value, usage, opts...)
	return &p
}

// StringSlice defines a []string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
func StringSlice(name string, value []string, usage string, opts ...Opt) *[]string {
	return CommandLine.StringSlice(name, value, usage, opts...)
}
