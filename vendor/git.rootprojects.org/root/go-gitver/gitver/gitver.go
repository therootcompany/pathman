package gitver

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var exactVer *regexp.Regexp
var gitVer *regexp.Regexp

func init() {
	// exactly vX.Y.Z (go-compatible semver)
	exactVer = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

	// vX.Y.Z-n-g0000000 git post-release, semver prerelease
	// vX.Y.Z-dirty git post-release, semver prerelease
	gitVer = regexp.MustCompile(`^(v\d+\.\d+)\.(\d+)(-(\d+))?(-(g[0-9a-f]+))?(-(dirty))?`)
}

type Versions struct {
	Timestamp time.Time
	Version   string
	Rev       string
}

func ExecAndParse() (*Versions, error) {
	desc, err := gitDesc()
	if nil != err {
		return nil, err
	}
	rev, err := gitRev()
	if nil != err {
		return nil, err
	}
	ver, err := semVer(desc)
	if nil != err {
		return nil, err
	}
	ts, err := gitTimestamp(desc)
	if nil != err {
		ts = time.Now()
	}

	return &Versions{
		Timestamp: ts,
		Version:   ver,
		Rev:       rev,
	}, nil
}

func gitDesc() (string, error) {
	args := strings.Split("git describe --tags --dirty --always", " ")
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if nil != err {
		// Don't panic, just carry on
		//out = []byte("v0.0.0-0-g0000000")
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func gitRev() (string, error) {
	args := strings.Split("git rev-parse HEAD", " ")
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if nil != err {
		return "", fmt.Errorf("\nUnexpected Error\n\n"+
			"Please open an issue at https://git.rootprojects.org/root/go-gitver/issues/new \n"+
			"Please include the following:\n\n"+
			"Command: %s\n"+
			"Output: %s\n"+
			"Error: %s\n"+
			"\nPlease and Thank You.\n\n", strings.Join(args, " "), out, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func semVer(desc string) (string, error) {
	if exactVer.MatchString(desc) {
		// v1.0.0
		return desc, nil
	}

	if !gitVer.MatchString(desc) {
		return "", nil
	}

	// (v1.0).(0)(-(1))(-(g0000000))(-(dirty))
	vers := gitVer.FindStringSubmatch(desc)
	patch, err := strconv.Atoi(vers[2])
	if nil != err {
		return "", fmt.Errorf("\nUnexpected Error\n\n"+
			"Please open an issue at https://git.rootprojects.org/root/go-gitver/issues/new \n"+
			"Please include the following:\n\n"+
			"git description: %s\n"+
			"RegExp: %#v\n"+
			"Error: %s\n"+
			"\nPlease and Thank You.\n\n", desc, gitVer, err)
	}

	// v1.0.1-pre1
	// v1.0.1-pre1+g0000000
	// v1.0.1-pre0+dirty
	// v1.0.1-pre0+g0000000-dirty
	if "" == vers[4] {
		vers[4] = "0"
	}
	ver := fmt.Sprintf("%s.%d-pre%s", vers[1], patch+1, vers[4])
	if "" != vers[6] || "dirty" == vers[8] {
		ver += "+"
		if "" != vers[6] {
			ver += vers[6]
			if "" != vers[8] {
				ver += "-"
			}
		}
		ver += vers[8]
	}

	return ver, nil
}

func gitTimestamp(desc string) (time.Time, error) {
	args := []string{
		"git",
		"show", desc,
		"--format=%cd",
		"--date=format:%Y-%m-%dT%H:%M:%SZ%z",
		"--no-patch",
	}
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if nil != err {
		// a dirty desc was probably used
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, strings.TrimSpace(string(out)))
}
