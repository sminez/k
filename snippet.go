package k

import (
	"fmt"
	"strings"
)

// A snippet is an individual entry into the help system that can be
// searched for and pretty printed.
type snippet struct {
	file  string
	title string
	tags  string
	lines []string
}

// build a new snippet struct from a range of lines from a given helpfile
func newSnippet(fname string, lines []string) *snippet {
	s := snippet{
		file:  strings.TrimRight(fname, ".txt"),
		lines: lines,
	}

	for _, l := range lines {
		if len(l) > 0 {
			switch c := l[0]; c {
			case titleMarker:
				s.title = l[2:]
			case tagMarker:
				s.tags = l[2:]
			}
		}
	}
	// TODO: We should probably error if title/tags are empty or set twice
	return &s
}

// Stringify for ANSI escaped pretty printing
// TODO: This feels kinda dirty as an impl for String()...
func (s *snippet) String() string {
	ansiString := func(colorCode int, line string) string {
		return fmt.Sprintf("\033[1;%dm%s\033[0m\n", colorCode, line)
	}

	var builder strings.Builder
	for i, line := range s.lines {
		if len(line) > 0 {
			switch firstChar := line[0]; firstChar {
			case titleMarker:
				builder.WriteString(ansiString(yellow, line[2:]))
			case urlMarker:
				builder.WriteString(ansiString(blue, line[2:]))
			case commentMarker:
				builder.WriteString(ansiString(grey, line[2:]))
			case codeMarker:
				builder.WriteString(fmt.Sprintf("%s\n", line))
			}
		} else {
			if i != 0 {
				builder.WriteString("\n")
			}
		}
	}
	return builder.String()
}

// Format this snippet for passing on to fzf as a selection
func (s *snippet) fzfString(fileLen, tagLen int) string {
	pad := func(n int, s string) string {
		k := n - len(s)
		if k == 0 {
			return s
		}
		return s + strings.Repeat(" ", k)
	}

	var builder strings.Builder
	builder.WriteString(pad(fileLen, s.file))
	builder.WriteString(seperator)
	builder.WriteString(pad(tagLen, s.tags))
	builder.WriteString(seperator)
	builder.WriteString(s.title)

	return builder.String()
}
