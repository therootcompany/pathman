package envpath

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type envConfig struct {
	shell      string
	shellDesc  string
	home       string
	rcFile     string
	rcScript   string
	loadFile   string
	loadScript string
}

var confs []*envConfig

func init() {
	home, err := os.UserHomeDir()
	if nil != err {
		panic(err) // Must get home directory
	}
	home = filepath.ToSlash(home)

	confs = []*envConfig{
		&envConfig{
			home:       home,
			shell:      "bash",
			shellDesc:  "bourne-compatible shell (bash)",
			rcFile:     ".bashrc",
			rcScript:   "[ -s \"$HOME/.config/envman/load.sh\" ] && source \"$HOME/.config/envman/load.sh\"\n",
			loadFile:   ".config/envman/load.sh",
			loadScript: "for x in ~/.config/envman/*.env; do\n\tsource \"$x\"\ndone\n",
		},
		&envConfig{
			home:       home,
			shell:      "zsh",
			shellDesc:  "bourne-compatible shell (zsh)",
			rcFile:     ".zshrc",
			rcScript:   "[ -s \"$HOME/.config/envman/load.sh\" ] && source \"$HOME/.config/envman/load.sh\"\n",
			loadFile:   ".config/envman/load.sh",
			loadScript: "for x in ~/.config/envman/*.env; do\n\tsource \"$x\"\ndone\n",
		},
		&envConfig{
			home:       home,
			shell:      "fish",
			shellDesc:  "fish shell",
			rcFile:     ".config/fish/config.fish",
			rcScript:   "test -s \"$HOME/.config/envman/load.fish\"; and source \"$HOME/.config/envman/load.fish\"\n",
			loadFile:   ".config/envman/load.fish",
			loadScript: "for x in ~/.config/envman/*.env\n\tsource \"$x\"\nend\n",
		},
	}
}

func initializeShells(home string) error {
	envmand := filepath.Join(home, ".config/envman")
	err := os.MkdirAll(envmand, 0755)
	if nil != err {
		return err
	}

	var hasRC bool
	var nativeMatch *envConfig
	shell := strings.TrimSuffix(filepath.Base(os.Getenv("SHELL")), ".exe")
	for i := range confs {
		c := confs[i]

		if shell == c.shell {
			nativeMatch = c
		}

		_, err := os.Stat(filepath.Join(home, c.rcFile))
		if nil != err {
			continue
		}
		hasRC = true
	}

	// ensure rc
	if !hasRC {
		if nil == nativeMatch {
			return fmt.Errorf(
				"%q is not a recognized shell and found none of .bashrc, .zshrc, .config/fish/config.fish",
				os.Getenv("SHELL"),
			)
		}

		// touch the rc file
		f, err := os.OpenFile(filepath.Join(home, nativeMatch.rcFile), os.O_CREATE|os.O_WRONLY, 0644)
		if nil != err {
			return err
		}
		if err := f.Close(); nil != err {
			return err
		}
	}

	// MacOS is special. It *requires* .bash_profile in order to read .bashrc
	if "darwin" == runtime.GOOS && "bash" == shell {
		if err := ensureBashProfile(home); nil != err {
			return err
		}
	}

	//
	// Bash (sh, dash, zsh, ksh)
	//
	// http://www.joshstaiger.org/archives/2005/07/bash_profile_vs.html
	for i := range confs {
		c := confs[i]
		err := c.maybeInitializeShell()
		if nil != err {
			return err
		}
	}

	return nil
}

func (c *envConfig) maybeInitializeShell() error {
	if _, err := os.Stat(filepath.Join(c.home, c.rcFile)); nil != err {
		if !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		return nil
	}

	changed, err := c.initializeShell()
	if nil != err {
		return err
	}

	if changed {
		fmt.Printf(
			"Detected %s shell and updated ~/%s\n",
			c.shellDesc,
			strings.TrimPrefix(c.rcFile, c.home),
		)
	}

	return nil
}

func (c *envConfig) initializeShell() (bool, error) {
	if err := c.ensurePathsLoader(); err != nil {
		return false, err
	}

	// Get current config
	// ex: ~/.bashrc
	// ex: ~/.config/fish/config.fish
	b, err := ioutil.ReadFile(filepath.Join(c.home, c.rcFile))
	if nil != err {
		return false, err
	}

	// For Windows, just in case
	s := strings.Replace(string(b), "\r\n", "\n", -1)

	// Looking to see if loader script has been added to rc file
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := range lines {
		line := lines[i]
		if line == strings.TrimSpace(c.rcScript) {
			// indicate that it was not neccesary to change the rc file
			return false, nil
		}
	}

	// Open rc file to append and write
	f, err := os.OpenFile(filepath.Join(c.home, c.rcFile), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}

	// Generate our script
	script := fmt.Sprintf("# Generated for envman. Do not edit.\n%s\n", c.rcScript)

	// If there's not a newline before our template,
	// include it in the template. We want nice things.
	n := len(lines)
	if "" != strings.TrimSpace(lines[n-1]) {
		script = "\n" + script
	}

	// Write and close the rc file
	if _, err := f.Write([]byte(script)); err != nil {
		return false, err
	}
	if err := f.Close(); err != nil {
		return true, err
	}

	// indicate that we have changed the rc file
	return true, nil
}

func (c *envConfig) ensurePathsLoader() error {
	loadFile := filepath.Join(c.home, c.loadFile)

	if _, err := os.Stat(loadFile); nil != err {
		// Write the loop file. For example:
		// $HOME/.config/envman/load.sh
		// $HOME/.config/envman/load.fish
		// TODO maybe don't write every time
		if err := ioutil.WriteFile(
			loadFile,
			[]byte(fmt.Sprintf("# Generated for envman. Do not edit.\n%s\n", c.loadScript)),
			os.FileMode(0755),
		); nil != err {
			return err
		}
		fmt.Printf("Created %s\n", "~/"+c.loadFile)
	}
	return nil
}

// I think this issue only affects darwin users with bash as the default shell
func ensureBashProfile(home string) error {
	profileFile := filepath.Join(home, ".bash_profile")

	// touch the profile file
	f, err := os.OpenFile(profileFile, os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		return err
	}
	if err := f.Close(); nil != err {
		return err
	}

	b, err := ioutil.ReadFile(profileFile)
	if !bytes.Contains(b, []byte(".bashrc")) {
		f, err := os.OpenFile(profileFile, os.O_APPEND|os.O_WRONLY, 0644)
		if nil != err {
			return err
		}
		sourceBashRC := "[ -s \"$HOME/.bashrc\" ] && source \"$HOME/.bashrc\"\n"
		b := []byte(fmt.Sprintf("# Generated for MacOS bash. Do not edit.\n%s\n", sourceBashRC))
		_, err = f.Write(b)
		if nil != err {
			return err
		}
		fmt.Printf("Updated ~/.bash_profile to source ~/.bashrc\n")
	}

	return nil
}
