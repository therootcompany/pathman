package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stdout, "Usage: envpath <action> [path]\n")
	fmt.Fprintf(os.Stdout, "\tex: envpath list\n")
	fmt.Fprintf(os.Stdout, "\tex: envpath add ~/.local/bin\n")
	fmt.Fprintf(os.Stdout, "\tex: envpath remove ~/.local/bin\n")
}

func main() {
	var action string
	var entry string

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
		return
	} else if len(os.Args) > 3 {
		usage()
		os.Exit(1)
		return
	}

	action = os.Args[1]
	if 2 == len(os.Args) {
		entry = os.Args[2]
	}

	// https://superuser.com/a/69190/73857
	// https://github.com/rust-lang-nursery/rustup.rs/issues/686#issuecomment-253982841
	// exec source $HOME/.profile
	shell := os.Getenv("SHELL")
	switch shell {
	case "":
		if strings.HasSuffix(os.Getenv("COMSPEC"), "/cmd.exe") {
			shell = "cmd"
		}
	case "fish":
		// ignore
	case "zsh":
		// ignore
	case "bash":
		// ignore
	default:
		// warn and try anyway
		fmt.Fprintf(
			os.Stderr,
			"%q isn't a recognized shell. Please open an issue at https://git.rootprojects.org/envpath/issues?q=%s",
			shell,
			shell,
		)
	}

	switch action {
	case "list":
		if 2 == len(os.Args) {
			usage()
			os.Exit(1)
		}
		list()
	case "add":
		add(entry)
	case "remove":
		remove(entry)
	}
}

func list() {
	managedpaths, err := listPaths()
	if nil != err {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	fmt.Println("pathman-managed PATH entries:\n")
	for i := range managedpaths {
		fmt.Println("\t" + managedpaths[i])
	}
	if 0 == len(managedpaths) {
		fmt.Println("\t(none)")
	}
	fmt.Println("")

	fmt.Println("other PATH entries:\n")
	// All managed paths
	pathsmap := map[string]bool{}
	for i := range managedpaths {
		// TODO normalize
		pathsmap[managedpaths[i]] = true
	}

	// Paths in the environment which are not managed
	var hasExtras bool
	envpaths := Paths()
	for i := range envpaths {
		// TODO normalize
		path := envpaths[i]
		if !pathsmap[path] {
			hasExtras = true
			fmt.Println("\t" + path)
		}
	}
	if !hasExtras {
		fmt.Println("\t(none)")
	}
	fmt.Println("")
}

func add(entry string) {
	// TODO noramlize away $HOME, %USERPROFILE%, etc
	abspath, err := filepath.Abs(entry)
	stat, err := os.Stat(entry)
	if nil != err {
		fmt.Fprintf(os.Stderr, "warning: couldn't access %q: %s\n", abspath, err)
	} else if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "warning: %q is not a directory", abspath)
	}

	modified, err := addPath(entry)
	if nil != err {
		fmt.Fprintf(os.Stderr, "failed to add %q to PATH: %s", entry, err)
		os.Exit(1)
	}

	var msg string
	if modified {
		msg = "Saved PATH changes."
	} else {
		msg = "PATH not changed."
	}

	paths := Paths()
	index := indexOfPath(Paths(), entry)
	if -1 == index {
		// TODO is os.PathListSeparator correct in MINGW / git bash?
		// generally this has no effect, but just in case this is included in a library with children processes
		paths = append([]string{entry}, paths...)
		err = os.Setenv(`PATH`, strings.Join(paths, string(os.PathListSeparator)))
		if nil != err {
			// ignore and carry on, as this is optional
			fmt.Fprintf(os.Stderr, "%s", err)
		}

		msg += " To set the PATH immediately, update the current session:\n\n\t" + Add(entry) + "\n"
	}

	fmt.Println(msg + "\n")
}

func remove(entry string) {
	modified, err := removePath(entry)
	if nil != err {
		fmt.Fprintf(os.Stderr, "failed to add %q to PATH: %s", entry, err)
		os.Exit(1)
	}

	var msg string
	if modified {
		msg = "Saved PATH changes."
	} else {
		msg = "PATH not changed."
	}

	paths := Paths()
	index := indexOfPath(Paths(), entry)
	if index >= 0 {
		newpaths := []string{}
		for i := range paths {
			if i != index {
				newpaths = append(newpaths, paths[i])
			}
		}
		// TODO is os.PathListSeparator correct in MINGW / git bash?
		// generally this has no effect, but just in case this is included in a library with children processes
		err = os.Setenv(`PATH`, strings.Join(newpaths, string(os.PathListSeparator)))
		if nil != err {
			// ignore and carry on, as this is optional
			fmt.Fprintf(os.Stderr, "%s", err)
		}

		msg += " To set the PATH immediately, update the current session:\n\n\t" + Remove(entry) + "\n"
	}

	fmt.Println(msg + "\n")
}

// Paths returns path entries in the current environment
func Paths() []string {
	cur := os.Getenv("PATH")
	if "" == cur {
		// unlikely, but possible... so whatever
		return nil
	}

	if isCmdExe() {
		//return strings.Split(cur, string(os.PathListSeparator))
		return strings.Split(cur, ";")
	}
	return strings.Split(cur, string(os.PathListSeparator))
}

// Add returns a string which can be used to add the given
// path entry to the current shell session
func Add(p string) string {
	if isCmdExe() {
		return fmt.Sprintf(`PATH %s;%PATH%`, p)
	}
	return fmt.Sprintf(`export PATH="%s:$PATH"`, p)
}

// Remove returns a string which can be used to remove the given
// path entry from the current shell session
func Remove(entries []string) string {
	if isCmdExe() {
		return fmt.Sprintf(`PATH %s`, strings.Join(entries, ";"))
	}
	return fmt.Sprintf(`export PATH="%s"`, strings.Join(entries, ":"))
}

func isCmdExe() {
	return "" == os.Getenv("SHELL") && strings.Contains(strings.ToLower(os.Getenv("COMSPEC")), "/cmd.exe")
}
