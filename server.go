package k

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// A Server handles preview requests from the spawned fzf process. It also holds
// all of the state we build up from parsing the helpfiles we find on init.
type Server struct {
	snippetMap map[string]*snippet
	stubs      []string
	port       int
}

// NewServer parallel reads all helpfiles that it finds before parsing into
// stubs for fzf.
// NOTE: The http server itself is created when calling ServeHTTP.
func NewServer(port int) *Server {
	var fl, tl, k int
	var ss []*snippet
	var ts []string

	m := make(map[string]*snippet)
	ch := make(chan []*snippet)
	ds := locateHelpfileDirectories()

	for _, d := range ds {
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

	return &Server{
		snippetMap: m,
		stubs:      ts,
		port:       port,
	}
}

// ServeHTTP starts a local http server for fzf to call back to in order to
// render its preview window. [NOTE: Runs in it's own goroutine]
func (k *Server) ServeHTTP() {
	http.HandleFunc("/", k.ansiEscapedFromStub)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", k.port), nil)
}

// pretty print the selected snippet based on the stub received
func (k *Server) ansiEscapedFromStub(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Unable to read data from POST body: %s\n", err)
		return
	}
	s := k.snippetMap[string(b)]
	if s != nil {
		fmt.Fprintf(w, "%s\n", s.ansiString())
	}
}

// RunFzf kicks off fzf as a subprocess and then waits for the selection to
// come back before rendering it out
func (k *Server) RunFzf(showPreview, copyToClipboard bool) {
	var p, s string

	if showPreview {
		p = fmt.Sprintf(
			"--preview %s %s",
			fmt.Sprintf(previewCmd, k.port),
			previewPosition,
		)
	}
	s = fmt.Sprintf(
		"echo '%s' | fzf %s",
		strings.Join(k.stubs, "\n"),
		p,
	)
	cmd := exec.Command("bash", "-c", s)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	t, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	snippet := k.snippetMap[strings.TrimRight(string(t), "\n")]

	if copyToClipboard {
		if err = copyStringToSystemClipboard(snippet.String()); err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("%s\n", snippet.ansiString())
	}
}