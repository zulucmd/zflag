// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
	"strings"
)

func getFlagWithDashes(name string) string {
	dash := "--"
	if len(name) == 1 {
		dash = "-"
	}

	return dash + name
}

// UnknownFlagError is returned when a flag is passed that has not been defined
// in the FlagSet.
type UnknownFlagError struct {
	name string
}

var _ error = (*UnknownFlagError)(nil)

// NewUnknownFlagError returns an error reporting that the named flag is unknown.
func NewUnknownFlagError(name string) error {
	return UnknownFlagError{name: name}
}

// Error implements the error interface.
func (e UnknownFlagError) Error() string {
	return fmt.Sprintf("unknown flag: %s", getFlagWithDashes(e.name))
}

// MissingFlagsError collects the names of required flags that were not set
// during parsing.
type MissingFlagsError []string

var _ error = (*MissingFlagsError)(nil)

// AddMissingFlag records the name of a required flag that was not set.
func (e *MissingFlagsError) AddMissingFlag(f *Flag) {
	*e = append(*e, getFlagWithDashes(f.Name))
}

// Error implements the error interface.
func (e MissingFlagsError) Error() string {
	flagNames := make([]string, 0, len(e))
	for _, s := range e {
		flagNames = append(flagNames, fmt.Sprintf("%q", s))
	}

	return fmt.Sprintf(`required flag(s) %s not set`, strings.Join(flagNames, `, `))
}

// InvalidArgumentError is returned when a flag value cannot be parsed or
// validated.
type InvalidArgumentError struct {
	flagName string
	value    any
	err      error
}

var _ error = (*InvalidArgumentError)(nil)

// NewInvalidArgumentError returns an error describing why the given value is
// invalid for the flag.
func NewInvalidArgumentError(err error, f *Flag, value any) error {
	var flagName string
	if f.Shorthand != 0 && f.ShorthandDeprecated == "" {
		flagName = fmt.Sprintf("-%c", f.Shorthand)
		if !f.ShorthandOnly {
			flagName = fmt.Sprintf("%s, --%s", flagName, f.Name)
		}
	} else {
		flagName = getFlagWithDashes(f.Name)
	}

	return InvalidArgumentError{
		flagName: flagName,
		value:    value,
		err:      err,
	}
}

// Error implements the error interface.
func (e InvalidArgumentError) Error() string {
	return fmt.Sprintf("invalid argument %q for %q flag: %s", e.value, e.flagName, e.err)
}

// Unwrap returns the underlying error so callers can inspect it with errors.Is
// and errors.As.
func (e InvalidArgumentError) Unwrap() error {
	return e.err
}
