// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
)

// Opt configures a Flag when it is defined. Options are applied in the order
// they are passed to the flag definition functions.
type Opt func(f *Flag) error

func applyFlagOptions(f *Flag, options ...Opt) error {
	for _, option := range options {
		if err := option(f); err != nil {
			return err
		}
	}
	return nil
}

// OptAddNegative automatically adds a --no-<flag> option for boolean flags.
func OptAddNegative() Opt {
	return func(f *Flag) error {
		f.AddNegative = true
		return nil
	}
}

// OptShorthand sets the one-letter shorthand for the flag.
func OptShorthand(shorthand rune) Opt {
	return func(f *Flag) error {
		f.Shorthand = shorthand
		return nil
	}
}

// OptShorthandStr sets the one-letter shorthand for the flag from a string.
// It panics if the string is not a single rune.
func OptShorthandStr(shorthand string) Opt {
	r, err := shorthandStrToRune(shorthand)
	if err != nil {
		panic(err)
	}

	return OptShorthand(r)
}

// OptShorthandOnly restricts the flag to its shorthand form, so the long name
// is not accepted.
func OptShorthandOnly() Opt {
	return func(f *Flag) error {
		f.ShorthandOnly = true
		return nil
	}
}

// OptUsage sets the help message shown for the flag.
func OptUsage(help string) Opt {
	return func(f *Flag) error {
		f.Usage = help
		return nil
	}
}

// OptUsageType sets the type name displayed for the flag in the help message.
func OptUsageType(usageType string) Opt {
	return func(f *Flag) error {
		f.UsageType = usageType
		return OptDisableUnquoteUsage()(f)
	}
}

// OptDisableUnquoteUsage disables unquoting and extraction of the type from
// the usage string.
func OptDisableUnquoteUsage() Opt {
	return func(f *Flag) error {
		f.DisableUnquoteUsage = true
		return nil
	}
}

// OptDisablePrintDefault disables printing of the default value in the usage
// message.
func OptDisablePrintDefault() Opt {
	return func(f *Flag) error {
		f.DisablePrintDefault = true
		return nil
	}
}

// OptDefValue sets the default value (as text) shown in the usage message.
func OptDefValue(defValue string) Opt {
	return func(f *Flag) error {
		f.DefValue = defValue
		return nil
	}
}

// OptDeprecated indicated that a flag is deprecated in your program. It will
// continue to function but will not show up in help or usage messages. Using
// this flag will also print the given usageMessage.
func OptDeprecated(msg string) Opt {
	return func(f *Flag) error {
		if msg == "" {
			return fmt.Errorf("deprecated message for flag %q must be set", f.Name)
		}

		f.Deprecated = msg
		return OptHidden()(f)
	}
}

// OptHidden hides the flag from help and usage text.
func OptHidden() Opt {
	return func(f *Flag) error {
		f.Hidden = true
		return nil
	}
}

// OptRequired marks the flag as required, so parsing fails if it is not set.
func OptRequired() Opt {
	return func(f *Flag) error {
		f.Required = true
		return nil
	}
}

// OptShorthandDeprecated marks the flag's shorthand as deprecated and prints
// msg as the replacement to use.
func OptShorthandDeprecated(msg string) Opt {
	return func(f *Flag) error {
		if msg == "" {
			return fmt.Errorf("shorthand deprecated message for flag %q must be set", f.Name)
		}

		f.ShorthandDeprecated = msg
		return nil
	}
}

// OptGroup assigns the flag to a named group.
func OptGroup(group string) Opt {
	return func(f *Flag) error {
		f.Group = group
		return nil
	}
}

// OptAnnotation attaches an application-specific annotation to the flag.
func OptAnnotation(key string, value []string) Opt {
	return func(f *Flag) error {
		f.SetAnnotation(key, value)
		return nil
	}
}
