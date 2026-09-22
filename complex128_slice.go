// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.15
// +build go1.15

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- complex128Slice Value
type complex128SliceValue struct {
	value   *[]complex128
	changed bool
}

var _ Value = (*complex128SliceValue)(nil)
var _ Getter = (*complex128SliceValue)(nil)
var _ SliceValue = (*complex128SliceValue)(nil)
var _ Typed = (*complex128SliceValue)(nil)

func newComplex128SliceValue(val []complex128, p *[]complex128) *complex128SliceValue {
	isv := new(complex128SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *complex128SliceValue) Get() interface{} {
	return *s.value
}

func (s *complex128SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	out, err := strconv.ParseComplex(val, 128)
	if err != nil {
		return errors.New("must be a complex number")
	}

	if !s.changed {
		*s.value = []complex128{}
	}
	*s.value = append(*s.value, out)
	s.changed = true

	return nil
}

func (s *complex128SliceValue) Type() string {
	return "complex128Slice"
}

func (s *complex128SliceValue) String() string {
	if s.value == nil || *s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%f", *s.value)
}

func (s *complex128SliceValue) fromString(val string) (complex128, error) {
	c, err := strconv.ParseComplex(val, 128)
	if err != nil {
		return 0, errors.New("must be a complex number")
	}
	return c, nil
}

func (s *complex128SliceValue) toString(val complex128) string {
	return fmt.Sprintf("%f", val)
}

func (s *complex128SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *complex128SliceValue) Replace(val []string) error {
	out := make([]complex128, len(val))
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

func (s *complex128SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetComplex128Slice return the []complex128 value of a flag with the given name
func (fs *FlagSet) GetComplex128Slice(name string) ([]complex128, error) {
	val, err := fs.getFlagValue(name, "complex128Slice")
	if err != nil {
		return []complex128{}, err
	}
	return val.([]complex128), nil
}

// MustGetComplex128Slice is like GetComplex128Slice, but panics on error.
func (fs *FlagSet) MustGetComplex128Slice(name string) []complex128 {
	val, err := fs.GetComplex128Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Complex128SliceVar defines a complex128Slice flag with specified name, default value, and usage string.
// The argument p points to a []complex128 variable in which to store the value of the flag.
func (fs *FlagSet) Complex128SliceVar(p *[]complex128, name string, value []complex128, usage string, opts ...Opt) {
	fs.Var(newComplex128SliceValue(value, p), name, usage, opts...)
}

// Complex128SliceVar defines a complex128[] flag with specified name, default value, and usage string.
// The argument p points to a complex128[] variable in which to store the value of the flag.
func Complex128SliceVar(p *[]complex128, name string, value []complex128, usage string, opts ...Opt) {
	CommandLine.Complex128SliceVar(p, name, value, usage, opts...)
}

// Complex128Slice defines a []complex128 flag with specified name, default value, and usage string.
// The return value is the address of a []complex128 variable that stores the value of the flag.
func (fs *FlagSet) Complex128Slice(name string, value []complex128, usage string, opts ...Opt) *[]complex128 {
	var p []complex128
	fs.Complex128SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Complex128Slice defines a []complex128 flag with specified name, default value, and usage string.
// The return value is the address of a []complex128 variable that stores the value of the flag.
func Complex128Slice(name string, value []complex128, usage string, opts ...Opt) *[]complex128 {
	return CommandLine.Complex128Slice(name, value, usage, opts...)
}
