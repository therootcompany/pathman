// +build windows

package main

import (
	"git.rootprojects.org/root/pathman/winpath"
)

var pathstore = "[winpath] "

func addPath(p string) (bool, error) {
	return winpath.Add(p)
}

func removePath(p string) (bool, error) {
	return winpath.Remove(p)
}

func listPaths() ([]string, error) {
	return winpath.Paths()
}

func indexOfPath(cur []string, p string) int {
	return winpath.IndexOf(cur, p)
}
