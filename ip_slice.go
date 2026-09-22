// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// -- ipSlice Value
type ipSliceValue struct {
	value   *[]net.IP
	changed bool
}

var _ Value = (*ipSliceValue)(nil)
var _ Getter = (*ipSliceValue)(nil)
var _ SliceValue = (*ipSliceValue)(nil)
var _ Typed = (*ipSliceValue)(nil)

func newIPSliceValue(val []net.IP, p *[]net.IP) *ipSliceValue {
	ipsv := new(ipSliceValue)
	ipsv.value = p
	*ipsv.value = val
	return ipsv
}

// Set converts, and assigns, the IP argument string representation as the []net.IP value of this flag.
// If Set is called on a flag that already has a []net.IP assigned, the newly converted values will be appended.
func (s *ipSliceValue) Set(val string) error {
	val = strings.TrimSpace(val)
	ip := net.ParseIP(val)
	if ip == nil {
		return errors.New("invalid string being converted to IP address")
	}

	if !s.changed {
		*s.value = []net.IP{}
	}
	*s.value = append(*s.value, ip)

	s.changed = true

	return nil
}

func (s *ipSliceValue) Get() any {
	return *s.value
}

// Type returns a string that uniquely represents this flag's type.
func (s *ipSliceValue) Type() string {
	return "ipSlice"
}

// String defines a "native" format for this net.IP slice flag value.
func (s *ipSliceValue) String() string {
	if s.value == nil {
		return "[]"
	}

	return fmt.Sprintf("%s", *s.value)
}

func (s *ipSliceValue) fromString(val string) (net.IP, error) {
	ip := net.ParseIP(strings.TrimSpace(val))
	if ip == nil {
		return nil, fmt.Errorf("invalid string being converted to IP address: %s", val)
	}
	return ip, nil
}

func (s *ipSliceValue) toString(val net.IP) string {
	return val.String()
}

func (s *ipSliceValue) Append(val string) error {
	ip, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, ip)
	return nil
}

func (s *ipSliceValue) Replace(val []string) error {
	out := make([]net.IP, len(val))
	for i, d := range val {
		ip, err := s.fromString(d)
		if err != nil {
			return err
		}
		out[i] = ip
	}
	*s.value = out
	return nil
}

func (s *ipSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

// GetIPSlice returns the []net.IP value of a flag with the given name
func (fs *FlagSet) GetIPSlice(name string) ([]net.IP, error) {
	val, err := fs.getFlagValue(name, "ipSlice")
	if err != nil {
		return []net.IP{}, err
	}
	return val.([]net.IP), nil
}

// MustGetIPSlice is like GetIPSlice, but panics on error.
func (fs *FlagSet) MustGetIPSlice(name string) []net.IP {
	val, err := fs.GetIPSlice(name)
	if err != nil {
		panic(err)
	}
	return val
}

// IPSliceVar defines a []net.IP flag with specified name, default value, and usage string.
// The argument p points to a []net.IP variable in which to store the value of the flag.
func (fs *FlagSet) IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string, opts ...Opt) {
	fs.Var(newIPSliceValue(value, p), name, usage, opts...)
}

// IPSliceVar defines a []net.IP flag with specified name, default value, and usage string.
// The argument p points to a []net.IP variable in which to store the value of the flag.
func IPSliceVar(p *[]net.IP, name string, value []net.IP, usage string, opts ...Opt) {
	CommandLine.IPSliceVar(p, name, value, usage, opts...)
}

// IPSlice defines a []net.IP flag with specified name, default value, and usage string.
// The return value is the address of a []net.IP variable that stores the value of the flag.
func (fs *FlagSet) IPSlice(name string, value []net.IP, usage string, opts ...Opt) *[]net.IP {
	var p []net.IP
	fs.IPSliceVar(&p, name, value, usage, opts...)
	return &p
}

// IPSlice defines a []net.IP flag with specified name, default value, and usage string.
// The return value is the address of a []net.IP variable that stores the value of the flag.
func IPSlice(name string, value []net.IP, usage string, opts ...Opt) *[]net.IP {
	return CommandLine.IPSlice(name, value, usage, opts...)
}
