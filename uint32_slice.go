// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- uint32Slice Value
type uint32SliceValue struct {
	value   *[]uint32
	changed bool
}

var _ Value = (*uint32SliceValue)(nil)
var _ Getter = (*uint32SliceValue)(nil)
var _ SliceValue = (*uint32SliceValue)(nil)
var _ Typed = (*uint32SliceValue)(nil)

func newUint32SliceValue(val []uint32, p *[]uint32) *uint32SliceValue {
	isv := new(uint32SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *uint32SliceValue) Get() interface{} {
	return *s.value
}

func (s *uint32SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	temp64, err := strconv.ParseUint(val, 0, 32)
	if err != nil {
		return errors.New("must be a non-negative integer")
	}

	if !s.changed {
		*s.value = []uint32{}
	}
	*s.value = append(*s.value, uint32(temp64))
	s.changed = true

	return nil
}

func (s *uint32SliceValue) Type() string {
	return "uint32Slice"
}

func (s *uint32SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%d", *s.value)
}

func (s *uint32SliceValue) fromString(val string) (uint32, error) {
	t64, err := strconv.ParseUint(val, 0, 32)
	if err != nil {
		return 0, errors.New("must be a non-negative integer")
	}
	return uint32(t64), nil
}

func (s *uint32SliceValue) toString(val uint32) string {
	return fmt.Sprintf("%d", val)
}

func (s *uint32SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *uint32SliceValue) Replace(val []string) error {
	out := make([]uint32, len(val))
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

func (s *uint32SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetUint32Slice return the []uint32 value of a flag with the given name
func (fs *FlagSet) GetUint32Slice(name string) ([]uint32, error) {
	val, err := fs.getFlagValue(name, "uint32Slice")
	if err != nil {
		return []uint32{}, err
	}
	return val.([]uint32), nil
}

// MustGetUint32Slice is like GetUint32Slice, but panics on error.
func (fs *FlagSet) MustGetUint32Slice(name string) []uint32 {
	val, err := fs.GetUint32Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Uint32SliceVar defines a []uint32 flag with specified name, default value, and usage string.
// The argument p points to a []uint32 variable in which to store the value of the flag.
func (fs *FlagSet) Uint32SliceVar(p *[]uint32, name string, value []uint32, usage string, opts ...Opt) {
	fs.Var(newUint32SliceValue(value, p), name, usage, opts...)
}

// Uint32SliceVar defines a []uint32 flag with specified name, default value, and usage string.
// The argument p points to a []uint32 variable in which to store the value of the flag.
func Uint32SliceVar(p *[]uint32, name string, value []uint32, usage string, opts ...Opt) {
	CommandLine.Uint32SliceVar(p, name, value, usage, opts...)
}

// Uint32Slice defines a []uint32 flag with specified name, default value, and usage string.
// The return value is the address of a []uint32 variable that stores the value of the flag.
func (fs *FlagSet) Uint32Slice(name string, value []uint32, usage string, opts ...Opt) *[]uint32 {
	var p []uint32
	fs.Uint32SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Uint32Slice defines a []uint32 flag with specified name, default value, and usage string.
// The return value is the address of a []uint32 variable that stores the value of the flag.
func Uint32Slice(name string, value []uint32, usage string, opts ...Opt) *[]uint32 {
	return CommandLine.Uint32Slice(name, value, usage, opts...)
}
