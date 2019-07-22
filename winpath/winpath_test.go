package winpath

import (
	"fmt"
	"os"
	"testing"
)

func TestNormalize(t *testing.T) {
	home, _ := os.UserHomeDir()

	absexp := ""
	homeexp := "%USERPROFILE%" + string(os.PathSeparator) + "foo"
	abspath, homepath := NormalizePathEntry(home + string(os.PathSeparator) + "foo")

	if absexp != abspath {
		t.Error(fmt.Errorf("Expected %q, but got %q", absexp, abspath))
	}

	if homeexp != homepath {
		t.Error(fmt.Errorf("Expected %q, but got %q", homeexp, homepath))
	}
}
