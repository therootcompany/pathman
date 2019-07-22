package envpath

import (
	"fmt"
	"strings"
	"testing"
)

const file = `# Generated for envman. Do not edit.
PATH="/foo"


# ignore
# ignore

PATH="/foo"
PATH="/foo:$PATH"
PATH="/foo:$PATH"
PATH="/foo:"$PATH"
PATH="/foo:""$PATH"
PATH=""

PATH=

JUNK=""
JUNK=
=""
=

whatever


PATH="/boo:$PATH"
PATH=""

`

var paths = []string{
	`PATH="/foo"`,
	`PATH="/foo:$PATH"`,
	`PATH=""`,
	`PATH="/boo:$PATH"`,
}

func TestParse(t *testing.T) {
	newlines, warnings := Parse("PATH", []byte(file))
	newfile := `PATH="` + strings.Join(newlines, "\"\n\tPATH=\"") + `"`
	expfile := strings.Join(paths, "\n\t")
	if newfile != expfile {
		t.Errorf("\nExpected:\n\t%s\nGot:\n\t%s", expfile, newfile)
	}
	for i := range warnings {
		w := warnings[i]
		fmt.Printf("warning dropping %q from line %d: %s\n", w.Message, w.LineNumber, w.Line)
	}
}
