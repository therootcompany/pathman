// +build !windows

package main

import (
	"git.rootprojects.org/root/pathman/envpath"
)

var pathstore = ""

func addPath(p string) (bool, error) {
	return envpath.Add(p)
}

func removePath(p string) (bool, error) {
	return envpath.Remove(p)
}

func listPaths() ([]string, error) {
	return envpath.Paths()
}

func indexOfPath(cur []string, p string) int {
	return envpath.IndexOf(cur, p)
}
