// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zflag_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/zulucmd/zflag/v2"
)

const expectedOutput = `      --long-form                 Some description
      --long-form2                Some description
                                    with multiline
  -s, --long-name                 Some description
      --[no-]long-name-negated    Some description
  -l, --[no-]long-name-negated2   Some description with
                                    multiline
  -t, --long-name2                Some description with
                                    multiline
`

func setUpZFlagSet(buf io.Writer) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ExitOnError)
	f.Bool("long-form", false, "Some description")
	f.Bool("long-form2", false, "Some description\n  with multiline")
	f.Bool("long-name", false, "Some description", zflag.OptShorthand('s'))
	f.Bool("long-name2", false, "Some description with\n  multiline", zflag.OptShorthand('t'))
	f.Bool("long-name-negated", false, "Some description", zflag.OptAddNegative())
	f.Bool("long-name-negated2", false, "Some description with\n  multiline", zflag.OptAddNegative(), zflag.OptShorthand('l'))
	f.SetOutput(buf)
	return f
}

func TestPrintUsage(t *testing.T) {
	var buf bytes.Buffer
	f := setUpZFlagSet(&buf)
	f.PrintDefaults()
	res := buf.String()
	if res != expectedOutput {
		t.Errorf("Expected \n%s \nActual \n%s", expectedOutput, res)
	}
}

func setUpZFlagSet2(buf io.Writer) *zflag.FlagSet {
	f := zflag.NewFlagSet("test", zflag.ExitOnError)
	f.Bool("long-form", false, "Some description", zflag.OptAddNegative())
	f.Bool("long-form2", false, "Some description\n  with multiline", zflag.OptAddNegative())
	f.Bool("long-name", false, "Some description", zflag.OptShorthand('s'))
	f.Bool("long-name2", false, "Some description with\n  multiline", zflag.OptShorthand('t'))
	f.String("some-very-long-arg", "test", "Some very long description having break the limit", zflag.OptShorthand('l'))
	f.String("other-very-long-arg", "long-default-value", "Some very long description having break the limit", zflag.OptShorthand('o'))
	f.String("some-very-long-arg2", "very long default value", "Some very long description\nwith line break\nmultiple")
	f.SetOutput(buf)
	return f
}

const expectedOutput2 = `      --[no-]long-form               Some description
      --[no-]long-form2              Some description
                                       with multiline
  -s, --long-name                    Some description
  -t, --long-name2                   Some description with
                                       multiline
  -o, --other-very-long-arg string   Some very long description having
                                     break the limit (default
                                     "long-default-value")
  -l, --some-very-long-arg string    Some very long description having
                                     break the limit (default "test")
      --some-very-long-arg2 string   Some very long description
                                     with line break
                                     multiple (default "very long default
                                     value")
`

func TestPrintUsage_2(t *testing.T) {
	var buf bytes.Buffer
	f := setUpZFlagSet2(&buf)
	res := f.FlagUsagesWrapped(80)
	if res != expectedOutput2 {
		t.Errorf("Expected \n%q \nActual \n%q", expectedOutput2, res)
	}
}

// A word wider than the description column must sit on its own line.
// Wrapping has to resume for the words after it (#501).
func TestFlagUsagesWrappedContinuesAfterUnbreakableWord(t *testing.T) {
	f := zflag.NewFlagSet("example", zflag.ContinueOnError)
	f.String("mount", "", "a mount specification e.g. 'type=bind,source=/opt,destination=/hostopt'. The rest of this description is never wrapped.")
	got := f.FlagUsagesWrapped(60)
	want := `      --mount string   a mount specification e.g.
                       'type=bind,source=/opt,destination=/hostopt'.
                       The rest of this description is
                       never wrapped.
`
	if got != want {
		t.Errorf("FlagUsagesWrapped(60)\nwant:\n%q\ngot:\n%q", want, got)
	}
}
