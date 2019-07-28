//go:generate go run -mod=vendor git.rootprojects.org/root/go-gitver

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GitRev is the git commit hash of the build
var GitRev = "000000000"

// GitVersion is the git description converted to semver
var GitVersion = "v0.5.2-pre+dirty"

// GitTimestamp is the timestamp of the latest commit
var GitTimestamp = time.Now().Format(time.RFC3339)

func usage() {
	fmt.Fprintf(os.Stdout, "Usage: pathman <action> [path]\n")
	fmt.Fprintf(os.Stdout, "\tex: pathman list\n")
	fmt.Fprintf(os.Stdout, "\tex: pathman add ~/.local/bin\n")
	fmt.Fprintf(os.Stdout, "\tex: pathman remove ~/.local/bin\n")
	fmt.Fprintf(os.Stdout, "\tex: pathman version\n")
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
	if 3 == len(os.Args) {
		entry = os.Args[2]
	}

	home, _ := os.UserHomeDir()
	if "" != entry && '~' == entry[0] {
		// Let windows users not to have to type %USERPROFILE% or \Users\me every time
		entry = strings.Replace(entry, "~", home, 1)
	}
	switch action {
	default:
		usage()
		os.Exit(1)
	case "version":
		fmt.Printf("pathman %s (%s) %s\n", GitVersion, GitRev, GitTimestamp)
		os.Exit(0)
		return
	case "list":
		if 2 != len(os.Args) {
			usage()
			os.Exit(1)
		}
		list()
	case "add":
		checkShell()
		add(entry)
	case "remove":
		checkShell()
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
	home, _ := os.UserHomeDir()
	for i := range managedpaths {
		pathsmap[managedpaths[i]] = true
	}

	// Paths in the environment which are not managed
	var hasExtras bool
	paths := Paths()
	for i := range paths {
		// TODO normalize
		path := paths[i]
		path1 := ""
		path2 := ""
		if strings.HasPrefix(path, home) {
			path1 = "$HOME" + strings.TrimPrefix(path, home)
			path2 = "%USERPROFILE%" + strings.TrimPrefix(path, home)
		}
		if !pathsmap[path] && !pathsmap[path1] && !pathsmap[path2] {
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
		fmt.Fprintf(os.Stderr, "%sfailed to add %q to PATH: %s", pathstore, entry, err)
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

		msg += " To set the PATH immediately, update the current session:\n\n\t" + Remove(newpaths) + "\n"
	}

	fmt.Println(msg + "\n")
}

// warns if this is an unknown / untested shell
func checkShell() {
	// https://superuser.com/a/69190/73857
	// https://github.com/rust-lang-nursery/rustup.rs/issues/686#issuecomment-253982841
	// exec source $HOME/.profile
	shellexe := filepath.Base(os.Getenv("SHELL"))
	shell := strings.TrimSuffix(shellexe, ".exe")
	switch shell {
	case ".":
		shell = ""
		fallthrough
	case "":
		if strings.HasSuffix(os.Getenv("COMSPEC"), "\\cmd.exe") {
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
			"%q isn't a recognized shell. Please open an issue at https://git.rootprojects.org/root/pathman/issues?q=%s\n",
			shellexe,
			shellexe,
		)
	}
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
		return fmt.Sprintf(`PATH %s;%%PATH%%`, strings.Replace(p, "%", "%%", -1))
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

func isCmdExe() bool {
	return "" == os.Getenv("SHELL") && strings.HasSuffix(strings.ToLower(os.Getenv("COMSPEC")), "\\cmd.exe")
}
