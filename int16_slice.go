// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- int16Slice Value
type int16SliceValue struct {
	value   *[]int16
	changed bool
}

var _ Value = (*int16SliceValue)(nil)
var _ Getter = (*int16SliceValue)(nil)
var _ SliceValue = (*int16SliceValue)(nil)
var _ Typed = (*int16SliceValue)(nil)

func newInt16SliceValue(val []int16, p *[]int16) *int16SliceValue {
	isv := new(int16SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *int16SliceValue) Get() interface{} {
	return *s.value
}

func (s *int16SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	temp64, err := strconv.ParseInt(val, 0, 16)
	if err != nil {
		return errors.New("must be an integer")
	}

	if !s.changed {
		*s.value = []int16{}
	}
	*s.value = append(*s.value, int16(temp64))
	s.changed = true

	return nil
}

func (s *int16SliceValue) Type() string {
	return "int16Slice"
}

func (s *int16SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%d", s.value)
}

func (s *int16SliceValue) fromString(val string) (int16, error) {
	t64, err := strconv.ParseInt(val, 0, 16)
	if err != nil {
		return 0, errors.New("must be an integer")
	}
	return int16(t64), nil
}

func (s *int16SliceValue) toString(val int16) string {
	return fmt.Sprintf("%d", val)
}

func (s *int16SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *int16SliceValue) Replace(val []string) error {
	out := make([]int16, len(val))
	for i, d := range val {
		var err error
		out[i], err = s.fromString(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *int16SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetInt16Slice return the []int16 value of a flag with the given name
func (fs *FlagSet) GetInt16Slice(name string) ([]int16, error) {
	val, err := fs.getFlagValue(name, "int16Slice")
	if err != nil {
		return []int16{}, err
	}
	return val.([]int16), nil
}

// MustGetInt16Slice is like GetInt16Slice, but panics on error.
func (fs *FlagSet) MustGetInt16Slice(name string) []int16 {
	val, err := fs.GetInt16Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int16SliceVar defines a []int16 flag with specified name, default value, and usage string.
// The argument p points to a []int16 variable in which to store the value of the flag.
func (fs *FlagSet) Int16SliceVar(p *[]int16, name string, value []int16, usage string, opts ...Opt) {
	fs.Var(newInt16SliceValue(value, p), name, usage, opts...)
}

// Int16SliceVar defines a []int16 flag with specified name, default value, and usage string.
// The argument p points to a []int16 variable in which to store the value of the flag.
func Int16SliceVar(p *[]int16, name string, value []int16, usage string, opts ...Opt) {
	CommandLine.Int16SliceVar(p, name, value, usage, opts...)
}

// Int16Slice defines a []int16 flag with specified name, default value, and usage string.
// The return value is the address of a []int16 variable that stores the value of the flag.
func (fs *FlagSet) Int16Slice(name string, value []int16, usage string, opts ...Opt) *[]int16 {
	var p []int16
	fs.Int16SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Int16Slice defines a []int16 flag with specified name, default value, and usage string.
// The return value is the address of a []int16 variable that stores the value of the flag.
func Int16Slice(name string, value []int16, usage string, opts ...Opt) *[]int16 {
	return CommandLine.Int16Slice(name, value, usage, opts...)
}
