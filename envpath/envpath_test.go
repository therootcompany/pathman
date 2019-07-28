package envpath

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestAddRemove(t *testing.T) {
	paths, err := Paths()
	if nil != err {
		t.Error(err)
		return
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

	var exists bool
	paths, err = Paths()
	for i := range paths {
		if "/tmp/delete/me" == paths[i] {
			exists = true
		}
	}
	if !exists {
		fmt.Println("len(paths):", len(paths))
		t.Error(fmt.Errorf("Paths: should have had the entry: /tmp/delete/me"))
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

	exists = false
	paths, err = Paths()
	for i := range paths {
		if "/tmp/delete/me" == paths[i] {
			exists = true
		}
	}
	if !exists {
		fmt.Println("len(paths):", len(paths))
		t.Error(fmt.Errorf("Paths: should have had the entry: /tmp/delete/me"))
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

	exists = false
	paths, err = Paths()
	for i := range paths {
		if "/tmp/delete/me" == paths[i] {
			exists = true
		}
	}
	if exists {
		fmt.Println("len(paths):", len(paths))
		t.Error(fmt.Errorf("Paths: should not have had the entry: /tmp/delete/me"))
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

func TestHome(t *testing.T) {
	home, _ := os.UserHomeDir()

	modified, err := Add(filepath.Join(home, "deleteme"))
	if nil != err {
		t.Error(err)
		return
	}
	if !modified {
		t.Error(fmt.Errorf("Add $HOME/deleteme: should have modified"))
		return
	}

	modified, err = Add(filepath.Join(home, "deleteme"))
	if nil != err {
		t.Error(err)
		return
	}
	if modified {
		t.Error(fmt.Errorf("Add $HOME/deleteme: should not have modified"))
		return
	}

	exists := false
	paths, err := Paths()
	for i := range paths {
		if "$HOME/deleteme" == paths[i] {
			exists = true
		}
	}
	if !exists {
		fmt.Println("len(paths):", len(paths))
		t.Error(fmt.Errorf("Paths: should have had the entry: $HOME/deleteme"))
		return
	}

	modified, err = Remove(filepath.Join(home, "deleteme"))
	if nil != err {
		t.Error(err)
		return
	}
	if !modified {
		t.Error(fmt.Errorf("Remove $HOME/deleteme: should have modified"))
		return
	}

	exists = false
	paths, err = Paths()
	for i := range paths {
		if "$HOME/deleteme" == paths[i] {
			exists = true
		}
	}
	if exists {
		fmt.Println("len(paths):", len(paths))
		t.Error(fmt.Errorf("Paths: should not have had the entry: $HOME/deleteme"))
		return
	}

	modified, err = Remove(filepath.Join(home, "deleteme"))
	if nil != err {
		t.Error(err)
		return
	}
	if modified {
		t.Error(fmt.Errorf("Remove $HOME/deleteme: should not have modified"))
		return
	}
}
