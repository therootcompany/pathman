// +build windows

package winpath

import "testing"

func TestShow(t *testing.T) {
	paths, err := Paths()
	if nil != err {
		t.Error(err)
	}

	if len(paths) < 1 {
		t.Error("should have paths")
	}
}
