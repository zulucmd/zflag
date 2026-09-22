// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag

import (
	"fmt"
)

// FlagUsageFormatter is a function type that prints the usage for a single Flag.
// This should be returning two strings, one that is considered the "left" hand side,
// and one that is considered the "right" hand side.
// Once the left and right are determined for all flags, the length of the text is
// determined, and each is appropriated cut based the terminal's width, and some space
// is added between left and right.
type FlagUsageFormatter func(*Flag) (string, string)

func defaultUsageFormatter(flag *Flag) (string, string) {
	left := "  "
	if flag.Shorthand != 0 && flag.ShorthandDeprecated == "" {
		left += fmt.Sprintf("-%c", flag.Shorthand)
		if !flag.ShorthandOnly {
			left += ", "
		}
	} else {
		left += "    "
	}
	left += "--"
	if _, isBoolFlag := flag.Value.(BoolFlag); isBoolFlag && flag.AddNegative {
		left += "[no-]"
	}
	left += flag.Name

	varname, usage := UnquoteUsage(flag)
	if _, isBoolFlag := flag.Value.(BoolFlag); isBoolFlag && !flag.AddNegative {
		left += "[=true|false]"
	} else if varname != "" {
		left += " " + varname
	}

	right := usage
	if flag.Required {
		right += " (required)"
	}

	if !flag.DisablePrintDefault && !flag.DefaultIsZeroValue() {
		if v, ok := flag.Value.(Typed); ok && v.Type() == "string" {
			right += fmt.Sprintf(" (default %q)", flag.DefValue)
		} else {
			right += fmt.Sprintf(" (default %s)", flag.DefValue)
		}
	}
	if len(flag.Deprecated) != 0 {
		right += fmt.Sprintf(" (DEPRECATED: %s)", flag.Deprecated)
	}

	return left, right
}
