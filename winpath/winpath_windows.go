// +build windows

package winpath

// Needs to
//   * use the registry editor directly to avoid possible PATH  truncation
//     ( https://stackoverflow.com/questions/9546324/adding-directory-to-path-environment-variable-in-windows )
//     ( https://superuser.com/questions/387619/overcoming-the-1024-character-limit-with-setx )
//   * explicitly send WM_SETTINGCHANGE
//     ( https://github.com/golang/go/issues/18680#issuecomment-275582179 )

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func add(p string) (bool, error) {
	cur, err := paths()
	if nil != err {
		return false, err
	}

	index := IndexOf(cur, p)
	// skip silently, successfully
	if index >= 0 {
		return false, nil
	}

	cur = append([]string{p}, cur...)
	err = write(cur)
	if nil != err {
		return false, err
	}

	return true, nil
}

func remove(p string) (bool, error) {
	cur, err := paths()
	if nil != err {
		return false, err
	}

	index := IndexOf(cur, p)
	// skip silently, successfully
	if index < 0 {
		return false, nil
	}

	var newpaths []string
	for i := range cur {
		if i != index {
			newpaths = append(newpaths, cur[i])
		}
	}

	err = write(newpaths)
	if nil != err {
		return false, err
	}

	return true, nil
}

func write(cur []string) error {
	// TODO --system to add to the system PATH rather than the user PATH

	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("Can't open HKCU Environment for writes: %s", err)
	}
	defer k.Close()

	err = k.SetStringValue(`Path`, strings.Join(cur, string(os.PathListSeparator)))
	if nil != err {
		return fmt.Errorf("Can't set HKCU Environment[Path]: %s", err)
	}

	err = k.Close()
	if nil != err {
		return err
	}

	if nil != sendmsg {
		sendmsg()
	} else {
		fmt.Fprintf(os.Stderr, "Warning: added PATH, but you must reboot for changes to take effect\n")
	}

	return nil
}

func paths() ([]string, error) {
	// This is the canonical reference, which is actually quite nice to have.
	// TBH, it's a mess to do this on *nix systems.
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("Can't open HKCU Environment for reads: %s", err)
	}
	defer k.Close()

	// This is case insensitive on Windows.
	// PATH, Path, path will all work.
	s, _, err := k.GetStringValue("Path")
	if err != nil {
		return nil, fmt.Errorf("Can't query HKCU Environment[Path]: %s", err)
	}

	// ";" on Windows
	return strings.Split(s, string(os.PathListSeparator)), nil
}
