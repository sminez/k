// k - Answers to questions that you have already asked
package main

import (
	"flag"

	"github.com/sminez/k"
)

var (
	serverPort      = flag.Int("port", k.DefaultPort, "local port to listen to for preview requests")
	copyToClipboard = flag.Bool("clip", false, "copy the selected snippet to your clipboard instead of printing")
	showPreview     = flag.Bool("preview", true, "Enable ANSI color previews while searching within fzf")
)

func main() {
	flag.Parse()
	s := k.NewServer(*serverPort)

	if *showPreview {
		go s.Listen()
	}

	s.RunFzf(*showPreview, *copyToClipboard)
}
