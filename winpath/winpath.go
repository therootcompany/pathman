// Package winpath is useful for managing PATH as part of the Environment
// in the Windows HKey Local User registry. It returns an error for most
// operations on non-Windows systems.
package winpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrWrongPlatform indicates that this was not built for Windows
var ErrWrongPlatform = fmt.Errorf("method not implemented on this platform")

// sendmsg uses a syscall to broadcast the registry change so that
// new shells will get the new PATH immediately, without a reboot
var sendmsg func()

// Paths returns all PATHs according to the Windows HKLU registry
// (or nil on non-windows platforms)
func Paths() ([]string, error) {
	return paths()
}

// Add will rewrite the Windows registry HKLU Environment,
// prepending the given directory path to the user's PATH.
// It will return whether the PATH was modified and an
// error if it should have been modified, but wasn't.
func Add(p string) (bool, error) {
	return add(p)
}

// Remove will rewrite the Windows registry HKLU Environment
// without the given directory path.
// It will return whether the PATH was modified and an
// error if it should have been modified, but wasn't.
func Remove(p string) (bool, error) {
	return remove(p)
}

// NormalizePathEntry will return the given directory path relative
// from its absolute path to the %USERPROFILE% (home) directory.
func NormalizePathEntry(pathentry string) (string, string) {
	home, err := os.UserHomeDir()
	if nil != err {
		fmt.Fprintf(os.Stderr, "Couldn't get HOME directory. That's an unrecoverable hard fail.")
		panic(err)
	}

	sep := string(os.PathSeparator)
	absentry, _ := filepath.Abs(pathentry)
	home, _ = filepath.Abs(home)

	var homeentry string
	if strings.HasPrefix(strings.ToLower(absentry)+sep, strings.ToLower(home)+sep) {
		// %USERPROFILE% is allowed, but only for user PATH
		// https://superuser.com/a/442163/73857
		homeentry = `%USERPROFILE%` + absentry[len(home):]
	}

	if absentry == pathentry {
		absentry = ""
	}
	if homeentry == pathentry {
		homeentry = ""
	}

	return absentry, homeentry
}

// IndexOf searches the given path list for first occurence
// of the given path entry and returns the index, or -1
func IndexOf(paths []string, p string) int {
	abspath, homepath := NormalizePathEntry(p)

	index := -1
	for i := range paths {
		if strings.ToLower(p) == strings.ToLower(paths[i]) {
			index = i
			break
		}
		if strings.ToLower(abspath) == strings.ToLower(paths[i]) {
			index = i
			break
		}
		if strings.ToLower(homepath) == strings.ToLower(paths[i]) {
			index = i
			break
		}
	}

	return index
}
