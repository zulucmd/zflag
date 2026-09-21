// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"testing"

	"github.com/zulucmd/zflag/v2"
)

func TestFormatter(t *testing.T) {
	t.Parallel()

	f := zflag.NewFlagSet("test", zflag.ContinueOnError)
	f.FlagUsageFormatter = func(_ *zflag.Flag) (string, string) {
		return "--not-uis", "not-usage"
	}
	f.String("uis", "asom", "testing `varname` and usage", zflag.OptDeprecated("some msg"))
	f.Lookup("uis").Hidden = false
	actual := f.FlagUsages()
	expected := "--not-uis   not-usage\n"

	if actual != expected {
		t.Fatalf("expected %q but got %q", expected, actual)
	}
}
