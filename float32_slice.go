// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- float32Slice Value
type float32SliceValue struct {
	value   *[]float32
	changed bool
}

var _ Value = (*float32SliceValue)(nil)
var _ Getter = (*float32SliceValue)(nil)
var _ SliceValue = (*float32SliceValue)(nil)
var _ Typed = (*float32SliceValue)(nil)

func newFloat32SliceValue(val []float32, p *[]float32) *float32SliceValue {
	isv := new(float32SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *float32SliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	temp64, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return errors.New("must be a number")
	}

	if !s.changed {
		*s.value = []float32{}
	}
	*s.value = append(*s.value, float32(temp64))
	s.changed = true

	return nil
}

func (s *float32SliceValue) Get() any {
	return *s.value
}

func (s *float32SliceValue) Type() string {
	return "float32Slice"
}

func (s *float32SliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%f", *s.value)
}

func (s *float32SliceValue) fromString(val string) (float32, error) {
	t64, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, errors.New("must be a number")
	}
	return float32(t64), nil
}

func (s *float32SliceValue) toString(val float32) string {
	return fmt.Sprintf("%f", val)
}

func (s *float32SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *float32SliceValue) Replace(val []string) error {
	out := make([]float32, len(val))
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

func (s *float32SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetFloat32Slice return the []float32 value of a flag with the given name
func (fs *FlagSet) GetFloat32Slice(name string) ([]float32, error) {
	val, err := fs.getFlagValue(name, "float32Slice")
	if err != nil {
		return []float32{}, err
	}
	return val.([]float32), nil
}

// MustGetFloat32Slice is like GetFloat32Slice, but panics on error.
func (fs *FlagSet) MustGetFloat32Slice(name string) []float32 {
	val, err := fs.GetFloat32Slice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// Float32SliceVar defines a float32Slice flag with specified name, default value, and usage string.
// The argument p points to a []float32 variable in which to store the value of the flag.
func (fs *FlagSet) Float32SliceVar(p *[]float32, name string, value []float32, usage string, opts ...Opt) {
	fs.Var(newFloat32SliceValue(value, p), name, usage, opts...)
}

// Float32SliceVar defines a []float32 flag with specified name, default value, and usage string.
// The argument p points to a []float32 variable in which to store the value of the flag.
func Float32SliceVar(p *[]float32, name string, value []float32, usage string, opts ...Opt) {
	CommandLine.Float32SliceVar(p, name, value, usage, opts...)
}

// Float32Slice defines a []float32 flag with specified name, default value, and usage string.
// The return value is the address of a []float32 variable that stores the value of the flag.
func (fs *FlagSet) Float32Slice(name string, value []float32, usage string, opts ...Opt) *[]float32 {
	var p []float32
	fs.Float32SliceVar(&p, name, value, usage, opts...)
	return &p
}

// Float32Slice defines a []float32 flag with specified name, default value, and usage string.
// The return value is the address of a []float32 variable that stores the value of the flag.
func Float32Slice(name string, value []float32, usage string, opts ...Opt) *[]float32 {
	return CommandLine.Float32Slice(name, value, usage, opts...)
}
