// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- uint8Slice Value
type uint8SliceValue struct {
	value   *[]uint8
	changed bool
}

var _ Value = (*uint8SliceValue)(nil)
var _ Getter = (*uint8SliceValue)(nil)
var _ SliceValue = (*uint8SliceValue)(nil)
var _ Typed = (*uint8SliceValue)(nil)

func newUint8SliceValue(val []uint8, p *[]uint8) *uint8SliceValue {
	isv := new(uint8SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *uint8SliceValue) Get() interface{} {
	return *s.value
}

func (s *uint8SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	temp64, err := strconv.ParseUint(val, 0, 8)
	if err != nil {
		return errors.New("must be a non-negative integer")
	}

	if !s.changed {
		*s.value = []uint8{}
	}
	*s.value = append(*s.value, uint8(temp64))
	s.changed = true

	return nil
}

func (s *uint8SliceValue) Type() string {
	return "uint8Slice"
}

func (s *uint8SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%d", *s.value)
}

func (s *uint8SliceValue) fromString(val string) (uint8, error) {
	t64, err := strconv.ParseUint(val, 0, 8)
	if err != nil {
		return 0, errors.New("must be a non-negative integer")
	}
	return uint8(t64), nil
}

func (s *uint8SliceValue) toString(val uint8) string {
	return fmt.Sprintf("%d", val)
}

func (s *uint8SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *uint8SliceValue) Replace(val []string) error {
	out := make([]uint8, len(val))
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

func (s *uint8SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetUint8Slice return the []uint8 value of a flag with the given name
func (fs *FlagSet) GetUint8Slice(name string) ([]uint8, error) {
	val, err := fs.getFlagValue(name, "uint8Slice")
	if err != nil {
		return []uint8{}, err
	}
	return val.([]uint8), nil
}

// MustGetUint8Slice is like GetUint8Slice, but panics on error.
func (fs *FlagSet) MustGetUint8Slice(name string) []uint8 {
	val, err := fs.GetUint8Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint8SliceVar defines a []uint8 flag with specified name, default value, and usage string.
// The argument p points to a []uint8 variable in which to store the value of the flag.
func (fs *FlagSet) Uint8SliceVar(p *[]uint8, name string, value []uint8, usage string, opts ...Opt) {
	fs.Var(newUint8SliceValue(value, p), name, usage, opts...)
}

// Uint8SliceVar defines a []uint8 flag with specified name, default value, and usage string.
// The argument p points to a []uint8 variable in which to store the value of the flag.
func Uint8SliceVar(p *[]uint8, name string, value []uint8, usage string, opts ...Opt) {
	CommandLine.Uint8SliceVar(p, name, value, usage, opts...)
}

// Uint8Slice defines a []uint8 flag with specified name, default value, and usage string.
// The return value is the address of a []uint8 variable that stores the value of the flag.
func (fs *FlagSet) Uint8Slice(name string, value []uint8, usage string, opts ...Opt) *[]uint8 {
	var p []uint8
	fs.Uint8SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Uint8Slice defines a []uint8 flag with specified name, default value, and usage string.
// The return value is the address of a []uint8 variable that stores the value of the flag.
func Uint8Slice(name string, value []uint8, usage string, opts ...Opt) *[]uint8 {
	return CommandLine.Uint8Slice(name, value, usage, opts...)
}
