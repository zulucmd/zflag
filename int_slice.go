// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strconv"
	"strings"
)

// -- intSlice Value
type intSliceValue struct {
	value   *[]int
	changed bool
}

var _ Value = (*intSliceValue)(nil)
var _ Getter = (*intSliceValue)(nil)
var _ Typed = (*intSliceValue)(nil)
var _ SliceValue = (*intSliceValue)(nil)

func newIntSliceValue(val []int, p *[]int) *intSliceValue {
	isv := new(intSliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *intSliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	out, err := strconv.Atoi(val)
	if err != nil {
		return err
	}

	if !s.changed {
		*s.value = []int{}
	}
	*s.value = append(*s.value, out)
	s.changed = true

	return nil
}

func (s *intSliceValue) Get() interface{} {
	return *s.value
}

func (s *intSliceValue) Type() string {
	return "intSlice"
}

func (s *intSliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%d", *s.value)
}

func (s *intSliceValue) Append(val string) error {
	i, err := strconv.Atoi(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *intSliceValue) Replace(val []string) error {
	out := make([]int, len(val))
	for i, d := range val {
		var err error
		out[i], err = strconv.Atoi(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *intSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = strconv.Itoa(d)
	}
	return out
}

// GetIntSlice return the []int value of a flag with the given name
func (fs *FlagSet) GetIntSlice(name string) ([]int, error) {
	val, err := fs.getFlagValue(name, "intSlice")
	if err != nil {
		return []int{}, err
	}
	return val.([]int), nil
}

// MustGetIntSlice is like GetIntSlice, but panics on error.
func (fs *FlagSet) MustGetIntSlice(name string) []int {
	val, err := fs.GetIntSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IntSliceVar defines a []int flag with specified name, default value, and usage string.
// The argument p points to a []int variable in which to store the value of the flag.
func (fs *FlagSet) IntSliceVar(p *[]int, name string, value []int, usage string, opts ...Opt) {
	fs.Var(newIntSliceValue(value, p), name, usage, opts...)
}

// IntSliceVar defines a []int flag with specified name, default value, and usage string.
// The argument p points to a []int variable in which to store the value of the flag.
func IntSliceVar(p *[]int, name string, value []int, usage string, opts ...Opt) {
	CommandLine.IntSliceVar(p, name, value, usage, opts...)
}

// IntSlice defines a []int flag with specified name, default value, and usage string.
// The return value is the address of a []int variable that stores the value of the flag.
func (fs *FlagSet) IntSlice(name string, value []int, usage string, opts ...Opt) *[]int {
	var p []int
	fs.IntSliceVar(&p, name, value, usage, opts...)
	return &p
}

// IntSlice defines a []int flag with specified name, default value, and usage string.
// The return value is the address of a []int variable that stores the value of the flag.
func IntSlice(name string, value []int, usage string, opts ...Opt) *[]int {
	return CommandLine.IntSlice(name, value, usage, opts...)
}
