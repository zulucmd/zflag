// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- int64Slice Value
type int64SliceValue struct {
	value   *[]int64
	changed bool
}

var _ Value = (*int64SliceValue)(nil)
var _ Getter = (*int64SliceValue)(nil)
var _ SliceValue = (*int64SliceValue)(nil)
var _ Typed = (*int64SliceValue)(nil)

func newInt64SliceValue(val []int64, p *[]int64) *int64SliceValue {
	isv := new(int64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *int64SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	out, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		return errors.New("must be an integer")
	}

	if !s.changed {
		*s.value = []int64{}
	}
	*s.value = append(*s.value, out)
	s.changed = true

	return nil
}

func (s *int64SliceValue) Get() any {
	return *s.value
}

func (s *int64SliceValue) Type() string {
	return "int64Slice"
}

func (s *int64SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%d", *s.value)
}

func (s *int64SliceValue) fromString(val string) (int64, error) {
	i, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		return 0, errors.New("must be an integer")
	}
	return i, nil
}

func (s *int64SliceValue) toString(val int64) string {
	return fmt.Sprintf("%d", val)
}

func (s *int64SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *int64SliceValue) Replace(val []string) error {
	out := make([]int64, len(val))
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

func (s *int64SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetInt64Slice return the []int64 value of a flag with the given name
func (fs *FlagSet) GetInt64Slice(name string) ([]int64, error) {
	val, err := fs.getFlagValue(name, "int64Slice")
	if err != nil {
		return []int64{}, err
	}
	return val.([]int64), nil
}

// MustGetInt64Slice is like GetInt64Slice, but panics on error.
func (fs *FlagSet) MustGetInt64Slice(name string) []int64 {
	val, err := fs.GetInt64Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Int64SliceVar defines a []int64 flag with specified name, default value, and usage string.
// The argument p points to a []int64 variable in which to store the value of the flag.
func (fs *FlagSet) Int64SliceVar(p *[]int64, name string, value []int64, usage string, opts ...Opt) {
	fs.Var(newInt64SliceValue(value, p), name, usage, opts...)
}

// Int64SliceVar defines a []int64 flag with specified name, default value, and usage string.
// The argument p points to a []int64 variable in which to store the value of the flag.
func Int64SliceVar(p *[]int64, name string, value []int64, usage string, opts ...Opt) {
	CommandLine.Int64SliceVar(p, name, value, usage, opts...)
}

// Int64Slice defines a []int64 flag with specified name, default value, and usage string.
// The return value is the address of a []int64 variable that stores the value of the flag.
func (fs *FlagSet) Int64Slice(name string, value []int64, usage string, opts ...Opt) *[]int64 {
	var p []int64
	fs.Int64SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Int64Slice defines a []int64 flag with specified name, default value, and usage string.
// The return value is the address of a []int64 variable that stores the value of the flag.
func Int64Slice(name string, value []int64, usage string, opts ...Opt) *[]int64 {
	return CommandLine.Int64Slice(name, value, usage, opts...)
}
