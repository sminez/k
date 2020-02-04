// +build linux

package k

var (
	clipboardProg     = "xclip"
	clipboardProgArgs = []string{"-in", "-selection", "clipboard"}
)
