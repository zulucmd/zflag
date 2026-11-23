// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -- boolSlice Value
type boolSliceValue struct {
	value   *[]bool
	changed bool
}

var _ Value = (*boolSliceValue)(nil)
var _ Getter = (*boolSliceValue)(nil)
var _ SliceValue = (*boolSliceValue)(nil)
var _ Typed = (*boolSliceValue)(nil)

func newBoolSliceValue(val []bool, p *[]bool) *boolSliceValue {
	bsv := new(boolSliceValue)
	bsv.value = p
	*bsv.value = val
	return bsv
}

// Set converts, and assigns, the boolean argument string representation as the []bool value of this flag.
// If Set is called on a flag that already has a []bool assigned, the newly converted values will be appended.
func (s *boolSliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	b, err := strconv.ParseBool(val)
	if err != nil {
		return errors.New("must be true or false")
	}

	if !s.changed {
		*s.value = []bool{}
	}
	*s.value = append(*s.value, b)
	s.changed = true

	return nil
}

func (s *boolSliceValue) Get() any {
	return *s.value
}

// Type returns a string that uniquely represents this flag's type.
func (s *boolSliceValue) Type() string {
	return "boolSlice"
}

// String defines a "native" format for this boolean slice flag value.
func (s *boolSliceValue) String() string {
	if s.value == nil || *s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%t", *s.value)
}

func (s *boolSliceValue) fromString(val string) (bool, error) {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false, errors.New("must be true or false")
	}
	return b, nil
}

func (s *boolSliceValue) toString(val bool) string {
	return strconv.FormatBool(val)
}

func (s *boolSliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *boolSliceValue) Replace(val []string) error {
	out := make([]bool, len(val))
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

func (s *boolSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetBoolSlice returns the []bool value of a flag with the given name.
func (fs *FlagSet) GetBoolSlice(name string) ([]bool, error) {
	val, err := fs.getFlagValue(name, "boolSlice")
	if err != nil {
		return []bool{}, err
	}
	return val.([]bool), nil
}

// MustGetBoolSlice is like GetBoolSlice, but panics on error.
func (fs *FlagSet) MustGetBoolSlice(name string) []bool {
	val, err := fs.GetBoolSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// BoolSliceVar defines a boolSlice flag with specified name, default value, and usage string.
// The argument p points to a []bool variable in which to store the value of the flag.
func (fs *FlagSet) BoolSliceVar(p *[]bool, name string, value []bool, usage string, opts ...Opt) {
	fs.Var(newBoolSliceValue(value, p), name, usage, opts...)
}

// BoolSliceVar defines a []bool flag with specified name, default value, and usage string.
// The argument p points to a []bool variable in which to store the value of the flag.
func BoolSliceVar(p *[]bool, name string, value []bool, usage string, opts ...Opt) {
	CommandLine.BoolSliceVar(p, name, value, usage, opts...)
}

// BoolSlice defines a []bool flag with specified name, default value, and usage string.
// The return value is the address of a []bool variable that stores the value of the flag.
func (fs *FlagSet) BoolSlice(name string, value []bool, usage string, opts ...Opt) *[]bool {
	var p []bool
	fs.BoolSliceVar(&p, name, value, usage, opts...)
	return &p
}

// BoolSlice defines a []bool flag with specified name, default value, and usage string.
// The return value is the address of a []bool variable that stores the value of the flag.
func BoolSlice(name string, value []bool, usage string, opts ...Opt) *[]bool {
	return CommandLine.BoolSlice(name, value, usage, opts...)
}
