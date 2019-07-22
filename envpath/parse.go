package envpath

import (
	"fmt"
	"strings"
)

type Warning struct {
	LineNumber int
	Line       string
	Message    string
}

// Parse will return a list of paths from an export file
func Parse(envname string, b []byte) ([]string, []Warning) {
	s := string(b)
	s = strings.Replace(s, "\r\n", "\n", -1)

	badlines := []Warning{}
	newlines := []string{}
	entries := make(map[string]bool)
	lines := strings.Split(s, "\n")
	for i := range lines {
		line := strings.TrimPrefix(strings.TrimSpace(lines[i]), "export ")
		if "" == line {
			continue
		}
		if "# Generated for envman. Do not edit." == line {
			continue
		}

		if '#' == line[0] {
			badlines = append(badlines, Warning{
				LineNumber: i,
				Line:       line,
				Message:    "comment",
			})
			continue
		}

		index := strings.Index(line, "=")
		if index < 1 {
			badlines = append(badlines, Warning{
				LineNumber: i,
				Line:       line,
				Message:    "invalid assignment",
			})
			continue
		}

		env := line[:index]
		if env != envname {
			badlines = append(badlines, Warning{
				LineNumber: i,
				Line:       line,
				Message:    fmt.Sprintf("wrong name (%s != %s)", env, envname),
			})
			continue
		}

		val := line[index+1:]
		if len(val) < 2 || '"' != val[0] || '"' != val[len(val)-1] {
			badlines = append(badlines, Warning{
				LineNumber: i,
				Line:       line,
				Message:    "value not quoted",
			})
			continue
		}
		val = val[1 : len(val)-1]

		if strings.Contains(val, `"`) {
			badlines = append(badlines, Warning{
				LineNumber: i,
				Line:       line,
				Message:    "invalid quotes",
			})

			continue
		}

		// TODO normalize $HOME
		if entries[val] {
			badlines = append(badlines, Warning{
				LineNumber: i,
				Line:       line,
				Message:    "duplicate entry",
			})
			continue
		}
		entries[val] = true

		newlines = append(newlines, val)
	}

	return newlines, badlines
}
