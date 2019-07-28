package envpath

import (
	"fmt"
	"testing"
)

func TestAddRemove(t *testing.T) {
	paths, err := Paths()
	if nil != err {
		t.Error(err)
		return
	}
	for i := range paths {
		fmt.Println(paths[i])
	}

	modified, err := Remove("/tmp/doesnt/exist")
	if nil != err {
		t.Error(err)
		return
	}
	if modified {
		t.Error(fmt.Errorf("Remove /tmp/doesnt/exist: should not have modified"))
		return
	}

	modified, err = Add("/tmp/delete/me")
	if nil != err {
		t.Error(err)
		return
	}
	if !modified {
		t.Error(fmt.Errorf("Add /tmp/delete/me: should have modified"))
		return
	}

	paths, err = Paths()
	if 1 != len(paths) || "/tmp/delete/me" != paths[0] {
		fmt.Println("len(paths):", len(paths))
		t.Error(fmt.Errorf("Paths: should have had exactly one entry: /tmp/delete/me"))
		return
	}

	modified, err = Add("/tmp/delete/me")
	if nil != err {
		t.Error(err)
		return
	}
	if modified {
		t.Error(fmt.Errorf("Add /tmp/delete/me: should not have modified"))
		return
	}

	paths, err = Paths()
	if 1 != len(paths) || "/tmp/delete/me" != paths[0] {
		t.Error(fmt.Errorf("Paths: should have had exactly one entry: /tmp/delete/me"))
		return
	}

	modified, err = Remove("/tmp/doesnt/exist")
	if nil != err {
		t.Error(err)
		return
	}
	if modified {
		t.Error(fmt.Errorf("Remove /tmp/doesnt/exist: should not have modified"))
		return
	}

	modified, err = Remove("/tmp/delete/me")
	if nil != err {
		t.Error(err)
		return
	}
	if !modified {
		t.Error(fmt.Errorf("Remove /tmp/delete/me: should have modified"))
		return
	}

	paths, err = Paths()
	if 0 != len(paths) {
		t.Error(fmt.Errorf("Paths: should have had no entries"))
		return
	}

	modified, err = Remove("/tmp/delete/me")
	if nil != err {
		t.Error(err)
		return
	}
	if modified {
		t.Error(fmt.Errorf("Remove /tmp/delete/me: should not have modified"))
		return
	}
}
