// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- uint16Slice Value
type uint16SliceValue struct {
	value   *[]uint16
	changed bool
}

var _ Value = (*uint16SliceValue)(nil)
var _ Getter = (*uint16SliceValue)(nil)
var _ SliceValue = (*uint16SliceValue)(nil)
var _ Typed = (*uint16SliceValue)(nil)

func newUint16SliceValue(val []uint16, p *[]uint16) *uint16SliceValue {
	isv := new(uint16SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *uint16SliceValue) Get() any {
	return *s.value
}

func (s *uint16SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	temp64, err := strconv.ParseUint(val, 0, 16)
	if err != nil {
		return errors.New("must be a non-negative integer")
	}

	if !s.changed {
		*s.value = []uint16{}
	}
	*s.value = append(*s.value, uint16(temp64))
	s.changed = true

	return nil
}

func (s *uint16SliceValue) Type() string {
	return "uint16Slice"
}

func (s *uint16SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%d", *s.value)
}

func (s *uint16SliceValue) fromString(val string) (uint16, error) {
	t64, err := strconv.ParseUint(val, 0, 16)
	if err != nil {
		return 0, errors.New("must be a non-negative integer")
	}
	return uint16(t64), nil
}

func (s *uint16SliceValue) toString(val uint16) string {
	return fmt.Sprintf("%d", val)
}

func (s *uint16SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *uint16SliceValue) Replace(val []string) error {
	out := make([]uint16, len(val))
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

func (s *uint16SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetUint16Slice return the []uint16 value of a flag with the given name
func (fs *FlagSet) GetUint16Slice(name string) ([]uint16, error) {
	val, err := fs.getFlagValue(name, "uint16Slice")
	if err != nil {
		return []uint16{}, err
	}
	return val.([]uint16), nil
}

// MustGetUint16Slice is like GetUint16Slice, but panics on error.
func (fs *FlagSet) MustGetUint16Slice(name string) []uint16 {
	val, err := fs.GetUint16Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint16SliceVar defines a []uint16 flag with specified name, default value, and usage string.
// The argument p points to a []uint16 variable in which to store the value of the flag.
func (fs *FlagSet) Uint16SliceVar(p *[]uint16, name string, value []uint16, usage string, opts ...Opt) {
	fs.Var(newUint16SliceValue(value, p), name, usage, opts...)
}

// Uint16SliceVar defines a []uint16 flag with specified name, default value, and usage string.
// The argument p points to a []uint16 variable in which to store the value of the flag.
func Uint16SliceVar(p *[]uint16, name string, value []uint16, usage string, opts ...Opt) {
	CommandLine.Uint16SliceVar(p, name, value, usage, opts...)
}

// Uint16Slice defines a []uint16 flag with specified name, default value, and usage string.
// The return value is the address of a []uint16 variable that stores the value of the flag.
func (fs *FlagSet) Uint16Slice(name string, value []uint16, usage string, opts ...Opt) *[]uint16 {
	var p []uint16
	fs.Uint16SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Uint16Slice defines a []uint16 flag with specified name, default value, and usage string.
// The return value is the address of a []uint16 variable that stores the value of the flag.
func Uint16Slice(name string, value []uint16, usage string, opts ...Opt) *[]uint16 {
	return CommandLine.Uint16Slice(name, value, usage, opts...)
}
