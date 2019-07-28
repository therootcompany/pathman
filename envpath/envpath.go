package envpath

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Paths parses the PATH.env file and returns a slice of valid paths
func Paths() ([]string, error) {
	home, err := os.UserHomeDir()
	if nil != err {
		return nil, err
	}
	home = filepath.ToSlash(home)

	_, paths, err := getEnv(home, "PATH")
	if nil != err {
		return nil, err
	}

	// ":" on *nix
	return paths, nil
}

// Add adds a path entry to the PATH env file
func Add(entry string) (bool, error) {
	home, err := os.UserHomeDir()
	if nil != err {
		return false, err
	}
	home = filepath.ToSlash(home)

	pathentry, err := normalizePathEntry(home, entry)
	if nil != err {
		return false, err
	}

	err = initializeShells(home)
	if nil != err {
		return false, err
	}

	fullpath, paths, err := getEnv(home, "PATH")
	if nil != err {
		return false, err
	}

	index := IndexOf(paths, pathentry)
	if index >= 0 {
		return false, nil
	}

	paths = append(paths, pathentry)
	err = writeEnv(fullpath, paths)
	if nil != err {
		return false, err
	}

	fmt.Println("Wrote " + fullpath)
	return true, nil
}

// Remove adds a path entry to the PATH env file
func Remove(entry string) (bool, error) {
	home, err := os.UserHomeDir()
	if nil != err {
		return false, err
	}
	home = filepath.ToSlash(home)

	pathentry, err := normalizePathEntry(home, entry)
	if nil != err {
		return false, err
	}

	err = initializeShells(home)
	if nil != err {
		return false, err
	}

	fullpath, oldpaths, err := getEnv(home, "PATH")
	if nil != err {
		return false, err
	}

	index := IndexOf(oldpaths, pathentry)
	if index < 0 {
		return false, nil
	}

	paths := []string{}
	for i := range oldpaths {
		if index != i {
			paths = append(paths, oldpaths[i])
		}
	}

	err = writeEnv(fullpath, paths)
	if nil != err {
		return false, err
	}

	fmt.Println("Wrote " + fullpath)
	return true, nil
}

func getEnv(home string, env string) (string, []string, error) {
	envmand := filepath.Join(home, ".config/envman")
	err := os.MkdirAll(envmand, 0755)
	if nil != err {
		return "", nil, err
	}

	nodes, err := ioutil.ReadDir(envmand)
	if nil != err {
		return "", nil, err
	}

	//filename := fmt.Sprintf("00-%s.env", env)
	filename := fmt.Sprintf("%s.env", env)
	for i := range nodes {
		name := nodes[i].Name()
		if fmt.Sprintf("%s.env", env) == name || strings.HasSuffix(name, fmt.Sprintf("-%s.env", env)) {
			filename = name
			break
		}
	}

	fullpath := filepath.Join(envmand, filename)
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_RDONLY, 0644)
	if nil != err {
		return "", nil, err
	}

	b, err := ioutil.ReadAll(f)
	f.Close()
	if nil != err {
		return "", nil, err
	}

	paths, warnings := Parse(b, env)
	for i := range warnings {
		w := warnings[i]
		fmt.Printf("warning: dropped %q from %s:%d: %s\n", w.Line, filename, w.LineNumber, w.Message)
	}

	pathlines := []string{}
	for i := range paths {
		pathname := strings.TrimSuffix(paths[i], ":$PATH")
		if strings.HasPrefix(pathname, "$PATH:") {
			fixed := strings.TrimPrefix(pathname, "$PATH:")
			fmt.Fprintf(os.Stderr, "warning: re-arranging $PATH:%s to %s:$PATH\n", fixed, fixed)
			pathname = fixed
		}
		pathlines = append(pathlines, pathname)
	}

	if len(warnings) > 0 {
		err := writeEnv(fullpath, pathlines)
		if nil != err {
			return "", nil, err
		}
	}

	return fullpath, pathlines, nil
}

func writeEnv(fullpath string, paths []string) error {
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if nil != err {
		return err
	}

	_, err = f.Write([]byte("# Generated for envman. Do not edit.\n"))
	if nil != err {
		return err
	}

	for i := range paths {
		_, err := f.Write([]byte(fmt.Sprintf("export PATH=\"%s:$PATH\"\n", paths[i])))
		if nil != err {
			return err
		}
	}

	return f.Close()
}

// IndexOf searches the given path list for first occurence
// of the given path entry and returns the index, or -1
func IndexOf(paths []string, p string) int {
	home, err := os.UserHomeDir()
	if nil != err {
		panic(err)
	}

	p, _ = normalizePathEntry(home, p)
	index := -1
	for i := range paths {
		entry, _ := normalizePathEntry(home, paths[i])
		if p == entry {
			index = i
			break
		}
	}
	return index
}

func normalizePathEntry(home, pathentry string) (string, error) {
	var err error

	// We add the slashes so that we don't get false matches
	// ex: foo should match foo/bar, but should NOT match foobar
	home, err = filepath.Abs(home)
	if nil != err {
		// I'm not sure how it's possible to get an error with Abs...
		return "", err
	}
	home += "/"
	pathentry, err = filepath.Abs(pathentry)
	if nil != err {
		return "", err
	}
	pathentry += "/"

	// Next we make the path relative to / or ~/
	// ex: /Users/me/.local/bin/ => .local/bin/
	if strings.HasPrefix(pathentry, home) {
		pathentry = "$HOME/" + strings.TrimPrefix(pathentry, home)
	}

	return strings.TrimSuffix(pathentry, "/"), nil
}
