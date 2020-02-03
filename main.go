// k - Answers to questions that you have already asked
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	// Default values for command line options
	helpfilePath             = "HELPFILE_PATH"
	defaultHelpfileDirectory = ".helpfiles" // under user homedir
	defaultPort              = 1992

	// Markup constants for snippet files
	endOfSectionMarker = "--"
	commentMarker      = '>'
	titleMarker        = '#'
	codeMarker         = '$'
	urlMarker          = '%'
	tagMarker          = '?'

	// UI and interaction with fzfo
	previewCmd      = "\"curl -s -X POST -d \"{}\" 127.0.0.1:%d/\""
	previewPosition = "--preview-window=up:70%"
	seperator       = " | " // Delimits sections in the lines passed to fzf
	yellow          = 33    // ANSI terminal escape sequence
	blue            = 34    // ANSI terminal escape sequence
	grey            = 37    // ANSI terminal escape sequence
)

var (
	// Flags
	serverPort      = flag.Int("port", defaultPort, "local port to listen to for preview requests")
	copyToClipboard = flag.Bool("clip", false, "copy the selected snippet to your clipboard instead of printing")
)

// A snippet is an individual entry into the help system that can be
// searched for and pretty printed.
type snippet struct {
	file  string
	title string
	tags  string
	lines []string
}

// A kServer handles preview requests from the spawned fzf process.
type kServer struct {
	snippetMap map[string]*snippet
	stubs      []string
	port       int
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

// Parallel read all helpfiles and parse into stubs for fzf
func newKServer(dirs []string, port int) *kServer {
	var fl, tl, k int
	var ss []*snippet
	var ts []string

	m := make(map[string]*snippet)
	ch := make(chan []*snippet)

	for _, d := range dirs {
		fs, err := ioutil.ReadDir(d)
		if err != nil {
			log.Fatal(err)
		}

		// NOTE: goro for each file
		for _, fh := range fs {
			f := fh.Name()
			if len(f) > fl {
				fl = len(f)
			}
			go getSnippetsFromFile(d, f, ch)
			k++
		}
	}

	for i := 0; i < k; i++ {
		ss = append(ss, <-ch...)
	}

	for _, s := range ss {
		if len(s.tags) > tl {
			tl = len(s.tags)
		}
	}

	for _, s := range ss {
		t := s.fzfString(fl-4, tl) // We trim off the '.txt' suffix
		ts = append(ts, t)
		m[t] = s
	}

	return &kServer{
		snippetMap: m,
		stubs:      ts,
		port:       port,
	}
}

// pretty print the selected snippet based on the stub received
func (k *kServer) ansiEscapedFromStub(w http.ResponseWriter, r *http.Request) {
	s, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Unable to read data from POST body: %s\n", err)
		return
	}

	fmt.Fprintf(w, "%s\n", k.snippetMap[string(s)])
}

// [NOTE: Runs in it's own goroutine]
// Start a local http server for fzf to call back to in order to render its preview window.
func (k *kServer) serveHTTP() {
	http.HandleFunc("/", k.ansiEscapedFromStub)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", k.port), nil)
}

// Kick off fzf as a subprocess and wait for the selection to come back.
func (k *kServer) runFzf() {
	s := fmt.Sprintf(
		"echo '%s' | fzf --preview %s %s",
		strings.Join(k.stubs, "\n"),
		fmt.Sprintf(previewCmd, k.port),
		previewPosition,
	)
	cmd := exec.Command("bash", "-c", s)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	t, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k.snippetMap[strings.TrimRight(string(t), "\n")])
}

// Determine where the user's helpfiles are located.
// $HELPFILE_PATH has the same semantics as $PATH
func locateHelpfileDirectories() []string {
	s, ok := os.LookupEnv(helpfilePath)
	if ok {
		return strings.Split(s, ":")
	}

	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return []string{path.Join(h, defaultHelpfileDirectory)}
}

func main() {
	flag.Parse()
	k := newKServer(locateHelpfileDirectories(), *serverPort)
	go k.serveHTTP()
	k.runFzf()
}
