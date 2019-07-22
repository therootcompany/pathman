// +build windows

package main

import (
	"git.rootprojects.org/root/pathman/envpath"
)

func addPath(p string) (bool, error) {
	return envpath.Add(p)
}

func removePath(p string) (bool, error) {
	return envpath.Remove(p)
}

func listPaths() ([]string, error) {
	return envpath.List()
}

func indexOfPath(cur []string, p string) int {
	return envpath.IndexOf(cur, p)
}
