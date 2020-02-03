// k - Answers to questions that you have already asked
//
// TODO:
//   - Copy selection to clipboard (no ANSI colors)
//   - Copy only code selection to clipboard
//   - Multiple helpfile directories
//   - Pre-filter on certain tags?
//   - Dump all tags
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
	envVarHelpfileDir        = "HELPFILE_DIR"
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
	helpfileDirectory string
	snippetMap        map[string]*snippet
	stubs             []string
	httpServer        *http.Server
	port              int
}

// Stringify for ANSI escaped pretty printing
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
		diff := n - len(s)
		if n == 0 {
			return s
		}
		return s + strings.Repeat(" ", diff)
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
func newSnippet(file string, lines []string) *snippet {
	var title, tags string

	for _, line := range lines {
		if len(line) > 0 {
			switch firstChar := line[0]; firstChar {
			case titleMarker:
				title = line[2:]
			case tagMarker:
				tags = line[2:]
			}
		}
	}

	return &snippet{
		file:  strings.TrimRight(file, ".txt"),
		tags:  tags,
		title: title,
		lines: lines,
	}
}

// Runs in it's own goroutine
func getSnippetsFromFile(fname, absPath string, ch chan []*snippet) {
	var snippets []*snippet
	var currentLines []string

	file, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == endOfSectionMarker {
			snippets = append(snippets, newSnippet(fname, currentLines))
			currentLines = []string{}
		} else {
			currentLines = append(currentLines, line)
		}
	}

	ch <- snippets
}

// Parallel read all helpfiles and parse into stubs for fzf
func newKServer(helpfileDirectory string, port int) *kServer {
	snippetMap := make(map[string]*snippet)
	ch := make(chan []*snippet)
	var snippets []*snippet
	var fileLen, tagLen int
	var stubs []string

	helpfiles, err := ioutil.ReadDir(helpfileDirectory)
	if err != nil {
		log.Fatal(err)
	}

	// Kick off a goroutine for each file to process them in parallel
	for _, hf := range helpfiles {
		fname := hf.Name()
		noExt := strings.TrimRight(fname, ".txt")
		if len(noExt) > fileLen {
			fileLen = len(noExt)
		}
		absPath := path.Join(helpfileDirectory, fname)

		go getSnippetsFromFile(fname, absPath, ch)
	}

	for i := 0; i < len(helpfiles); i++ {
		ss := <-ch
		snippets = append(snippets, ss...)
		for _, s := range ss {
			if len(s.tags) > tagLen {
				tagLen = len(s.tags)
			}
		}
	}

	for _, s := range snippets {
		stub := s.fzfString(fileLen, tagLen)
		stubs = append(stubs, stub)
		snippetMap[stub] = s
	}

	return &kServer{
		helpfileDirectory: helpfileDirectory,
		snippetMap:        snippetMap,
		stubs:             stubs,
		port:              port,
	}
}

// pretty print the selected snippet based on the stub received
func (k *kServer) ansiEscapedFromStub(w http.ResponseWriter, r *http.Request) {
	stub, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Unable to read data from POST body: %s\n", err)
		return
	}

	fmt.Fprintf(w, "%s\n", k.snippetMap[string(stub)])
}

// [Runs in it's own goroutine]
// Start a local http server for fzf to call back to in order to render its preview window.
func (k *kServer) serveHTTP() {
	router := http.NewServeMux()
	router.HandleFunc("/", k.ansiEscapedFromStub)
	port := fmt.Sprintf(":%d", k.port)
	k.httpServer = &http.Server{
		Addr:    port,
		Handler: router,
	}

	k.httpServer.ListenAndServe()
}

// Kick off fzf as a subprocess
func (k *kServer) runFzf() {
	cmdStr := fmt.Sprintf(
		"echo '%s' | fzf --preview %s %s",
		strings.Join(k.stubs, "\n"),
		fmt.Sprintf(previewCmd, k.port),
		previewPosition,
	)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	stub, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(k.snippetMap[strings.TrimRight(string(stub), "\n")])
}

// Determine where the user's helpfiles are located
func locateHelpfileDirectory() string {
	helpfileDirectory, ok := os.LookupEnv(envVarHelpfileDir)
	if ok {
		return helpfileDirectory
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return path.Join(userHome, defaultHelpfileDirectory)
}

func main() {
	flag.Parse()
	k := newKServer(locateHelpfileDirectory(), *serverPort)
	go k.serveHTTP()
	k.runFzf()
}
