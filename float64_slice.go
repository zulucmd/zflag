// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- float64Slice Value
type float64SliceValue struct {
	value   *[]float64
	changed bool
}

var _ Value = (*float64SliceValue)(nil)
var _ Getter = (*float64SliceValue)(nil)
var _ SliceValue = (*float64SliceValue)(nil)
var _ Typed = (*float64SliceValue)(nil)

func newFloat64SliceValue(val []float64, p *[]float64) *float64SliceValue {
	isv := new(float64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *float64SliceValue) Get() interface{} {
	return *s.value
}

func (s *float64SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	out, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return errors.New("must be a number")
	}

	if !s.changed {
		*s.value = []float64{}
	}
	*s.value = append(*s.value, out)
	s.changed = true

	return nil
}

func (s *float64SliceValue) Type() string {
	return "float64Slice"
}

func (s *float64SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%f", *s.value)
}

func (s *float64SliceValue) fromString(val string) (float64, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, errors.New("must be a number")
	}
	return f, nil
}

func (s *float64SliceValue) toString(val float64) string {
	return fmt.Sprintf("%f", val)
}

func (s *float64SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *float64SliceValue) Replace(val []string) error {
	out := make([]float64, len(val))
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

func (s *float64SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetFloat64Slice return the []float64 value of a flag with the given name
func (fs *FlagSet) GetFloat64Slice(name string) ([]float64, error) {
	val, err := fs.getFlagValue(name, "float64Slice")
	if err != nil {
		return []float64{}, err
	}
	return val.([]float64), nil
}

// MustGetFloat64Slice is like GetFloat64Slice, but panics on error.
func (fs *FlagSet) MustGetFloat64Slice(name string) []float64 {
	val, err := fs.GetFloat64Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Float64SliceVar defines a float64Slice flag with specified name, default value, and usage string.
// The argument p points to a []float64 variable in which to store the value of the flag.
func (fs *FlagSet) Float64SliceVar(p *[]float64, name string, value []float64, usage string, opts ...Opt) {
	fs.Var(newFloat64SliceValue(value, p), name, usage, opts...)
}

// Float64SliceVar defines a []float64 flag with specified name, default value, and usage string.
// The argument p points to a []float64 variable in which to store the value of the flag.
func Float64SliceVar(p *[]float64, name string, value []float64, usage string, opts ...Opt) {
	CommandLine.Float64SliceVar(p, name, value, usage, opts...)
}

// Float64Slice defines a []float64 flag with specified name, default value, and usage string.
// The return value is the address of a []float64 variable that stores the value of the flag.
func (fs *FlagSet) Float64Slice(name string, value []float64, usage string, opts ...Opt) *[]float64 {
	var p []float64
	fs.Float64SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Float64Slice defines a []float64 flag with specified name, default value, and usage string.
// The return value is the address of a []float64 variable that stores the value of the flag.
func Float64Slice(name string, value []float64, usage string, opts ...Opt) *[]float64 {
	return CommandLine.Float64Slice(name, value, usage, opts...)
}
