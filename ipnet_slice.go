// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"net"
	"strings"
)

// -- ipNetSlice Value
type ipNetSliceValue struct {
	value   *[]net.IPNet
	changed bool
}

var _ Value = (*ipNetSliceValue)(nil)
var _ Getter = (*ipNetSliceValue)(nil)
var _ SliceValue = (*ipNetSliceValue)(nil)
var _ Typed = (*ipNetSliceValue)(nil)

func newIPNetSliceValue(val []net.IPNet, p *[]net.IPNet) *ipNetSliceValue {
	ipnsv := new(ipNetSliceValue)
	ipnsv.value = p
	*ipnsv.value = val
	return ipnsv
}

func (s *ipNetSliceValue) Get() interface{} {
	return *s.value
}

// Set converts, and assigns, the IPNet argument string representation as the []net.IPNet value of this flag.
// If Set is called on a flag that already has a []net.IPNet assigned, the newly converted values will be appended.
func (s *ipNetSliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	_, n, err := net.ParseCIDR(val)
	if err != nil {
		return fmt.Errorf("invalid string being converted to CIDR: %s", val)
	}
	if n == nil {
		return fmt.Errorf("invalid string being converted to CIDR: %s", val)
	}

	if !s.changed {
		*s.value = []net.IPNet{*n}
	} else {
		*s.value = append(*s.value, *n)
	}

	s.changed = true

	return nil
}

// Type returns a string that uniquely represents this flag's type.
func (s *ipNetSliceValue) Type() string {
	return "ipNetSlice"
}

// String defines a "native" format for this net.IPNet slice flag value.
func (s *ipNetSliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%s", *s.value)
}

func (s *ipNetSliceValue) fromString(val string) (net.IPNet, error) {
	_, cidr, err := net.ParseCIDR(val)
	if err != nil {
		return net.IPNet{}, fmt.Errorf("invalid string being converted to CIDR: %s", val)
	}
	return *cidr, nil
}

func (s *ipNetSliceValue) toString(val net.IPNet) string {
	return val.String()
}

func (s *ipNetSliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *ipNetSliceValue) Replace(val []string) error {
	out := make([]net.IPNet, len(val))
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

func (s *ipNetSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetIPNetSlice returns the []net.IPNet value of a flag with the given name
func (fs *FlagSet) GetIPNetSlice(name string) ([]net.IPNet, error) {
	val, err := fs.getFlagValue(name, "ipNetSlice")
	if err != nil {
		return []net.IPNet{}, err
	}
	return val.([]net.IPNet), nil
}

// MustGetIPNetSlice is like GetIPNetSlice, but panics on error.
func (fs *FlagSet) MustGetIPNetSlice(name string) []net.IPNet {
	val, err := fs.GetIPNetSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IPNetSliceVar defines a []net.IPNet flag with specified name, default value, and usage string.
// The argument p points to a []net.IPNet variable in which to store the value of the flag.
func (fs *FlagSet) IPNetSliceVar(p *[]net.IPNet, name string, value []net.IPNet, usage string, opts ...Opt) {
	fs.Var(newIPNetSliceValue(value, p), name, usage, opts...)
}

// IPNetSliceVar defines a []net.IPNet flag with specified name, default value, and usage string.
// The argument p points to a []net.IPNet variable in which to store the value of the flag.
func IPNetSliceVar(p *[]net.IPNet, name string, value []net.IPNet, usage string, opts ...Opt) {
	CommandLine.IPNetSliceVar(p, name, value, usage, opts...)
}

// IPNetSlice defines a []net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of a []net.IPNet variable that stores the value of the flag.
func (fs *FlagSet) IPNetSlice(name string, value []net.IPNet, usage string, opts ...Opt) *[]net.IPNet {
	var p []net.IPNet
	fs.IPNetSliceVar(&p, name, value, usage, opts...)
	return &p
}

// IPNetSlice defines a []net.IPNet flag with specified name, default value, and usage string.
// The return value is the address of a []net.IPNet variable that stores the value of the flag.
func IPNetSlice(name string, value []net.IPNet, usage string, opts ...Opt) *[]net.IPNet {
	return CommandLine.IPNetSlice(name, value, usage, opts...)
}
