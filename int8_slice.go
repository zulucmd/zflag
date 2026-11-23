// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- int8Slice Value
type int8SliceValue struct {
	value   *[]int8
	changed bool
}

var _ Value = (*int8SliceValue)(nil)
var _ Getter = (*int8SliceValue)(nil)
var _ SliceValue = (*int8SliceValue)(nil)
var _ Typed = (*int8SliceValue)(nil)

func newInt8SliceValue(val []int8, p *[]int8) *int8SliceValue {
	isv := new(int8SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *int8SliceValue) Get() any {
	return *s.value
}

func (s *int8SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	temp64, err := strconv.ParseInt(val, 0, 8)
	if err != nil {
		return errors.New("must be an integer")
	}

	if !s.changed {
		*s.value = []int8{}
	}
	*s.value = append(*s.value, int8(temp64))
	s.changed = true

	return nil
}

func (s *int8SliceValue) Type() string {
	return "int8Slice"
}

func (s *int8SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%d", *s.value)
}

func (s *int8SliceValue) fromString(val string) (int8, error) {
	t64, err := strconv.ParseInt(val, 0, 8)
	if err != nil {
		return 0, errors.New("must be an integer")
	}
	return int8(t64), nil
}

func (s *int8SliceValue) toString(val int8) string {
	return fmt.Sprintf("%d", val)
}

func (s *int8SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *int8SliceValue) Replace(val []string) error {
	out := make([]int8, len(val))
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

func (s *int8SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetInt8Slice return the []int8 value of a flag with the given name
func (fs *FlagSet) GetInt8Slice(name string) ([]int8, error) {
	val, err := fs.getFlagValue(name, "int8Slice")
	if err != nil {
		return []int8{}, err
	}
	return val.([]int8), nil
}

// MustGetInt8Slice is like GetInt8Slice, but panics on error.
func (fs *FlagSet) MustGetInt8Slice(name string) []int8 {
	val, err := fs.GetInt8Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int8SliceVar defines a []int8 flag with specified name, default value, and usage string.
// The argument p points to a []int8 variable in which to store the value of the flag.
func (fs *FlagSet) Int8SliceVar(p *[]int8, name string, value []int8, usage string, opts ...Opt) {
	fs.Var(newInt8SliceValue(value, p), name, usage, opts...)
}

// Int8SliceVar defines a []int8 flag with specified name, default value, and usage string.
// The argument p points to a []int8 variable in which to store the value of the flag.
func Int8SliceVar(p *[]int8, name string, value []int8, usage string, opts ...Opt) {
	CommandLine.Int8SliceVar(p, name, value, usage, opts...)
}

// Int8Slice defines a []int8 flag with specified name, default value, and usage string.
// The return value is the address of a []int8 variable that stores the value of the flag.
func (fs *FlagSet) Int8Slice(name string, value []int8, usage string, opts ...Opt) *[]int8 {
	var p []int8
	fs.Int8SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Int8Slice defines a []int8 flag with specified name, default value, and usage string.
// The return value is the address of a []int8 variable that stores the value of the flag.
func Int8Slice(name string, value []int8, usage string, opts ...Opt) *[]int8 {
	return CommandLine.Int8Slice(name, value, usage, opts...)
}
