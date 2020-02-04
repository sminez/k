package k

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

// [NOTE: Runs in it's own goroutine]
// Open up the given file and parse each of the snippets we find. Results are
// send back on a channel once all snippets are found so that the contents of
// each helpfile render together when passed on to fzf. This also means that
// we don't need to do anything clever around signalling to determine when we
// are done parsing each file: we just wait for each file to be parsed.
func getSnippetsFromFile(dir, fname string, ch chan []*snippet) {
	var ss []*snippet
	var ls []string

	f, err := os.Open(path.Join(dir, fname))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	s := bufio.NewScanner(f)

	for s.Scan() {
		l := s.Text()
		if l == endOfSectionMarker {
			ss = append(ss, newSnippet(fname, ls))
			ls = []string{}
		} else {
			ls = append(ls, l)
		}
	}
	ch <- ss
}

// $HELPFILE_PATH has the same semantics as $PATH
func locateHelpfileDirectories() []string {
	var ds []string

	s, ok := os.LookupEnv(helpfilePath)
	if ok {
		ds = strings.Split(s, ":")
	}

	// Add ~/.helpfiles if it exists
	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	p := path.Join(h, defaultHelpfileDirectory)

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return ds
	}

	return append(ds, p)
}

func copyStringToSystemClipboard(s string) error {
	cmd := exec.Command(clipboardProg)

	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if _, err = in.Write([]byte(s)); err != nil {
		return err
	}

	if err = in.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}
