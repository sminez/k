package k

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

// A Server handles preview requests from the spawned fzf process. It also holds
// all of the state we build up from parsing the helpfiles we find on init.
type Server struct {
	tcpListener net.Listener
	snippetMap  map[string]*snippet
	stubs       []string
	port        int
}

// NewServer parallel reads all helpfiles that it finds before parsing into
// stubs for fzf. Provided there were no errors in parsing the directories
// we found, it will also bind a TCP listener to <port> ready to accept
// incoming connections.
func NewServer(port int) *Server {
	var fl, tl, k int
	var ss, fs []*snippet
	var ts []string

	m := make(map[string]*snippet)
	ch := make(chan []*snippet)
	ds := locateHelpfileDirectories()

	for _, d := range ds {
		fs, err := ioutil.ReadDir(d)
		if err != nil {
			log.Fatal(err)
		}

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
		fs = <-ch

		for _, s := range fs {
			if len(s.tags) > tl {
				tl = len(s.tags)
			}
		}
		ss = append(ss, fs...)
	}

	for _, s := range ss {
		t := s.fzfString(fl-4, tl) // We trim off the '.txt' suffix
		ts = append(ts, t)
		m[t] = s
	}

	l, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	return &Server{
		tcpListener: l,
		snippetMap:  m,
		stubs:       ts,
		port:        port,
	}
}

// BindAndListen starts a local tcp server for fzf to call back to in order to
// render its preview window. [NOTE: Runs in it's own goroutine]
func (k *Server) Listen() {
	for {
		// silently dropping failed incoming connections
		conn, _ := k.tcpListener.Accept()
		go k.ansiEscapedFromStub(conn)
	}
}

// pretty print the selected snippet based on the stub received
func (k *Server) ansiEscapedFromStub(conn net.Conn) {
	var resp string

	r, _ := bufio.NewReader(conn).ReadString('\n')
	defer conn.Close()

	r = strings.TrimRight(r, "\n")
	s, ok := k.snippetMap[r]

	if ok {
		resp = s.ansiString()
	} else {
		resp = fmt.Sprintf("-- ERROR --\n\nunable to find snippet for '%s'", r)
	}

	conn.Write([]byte(resp))
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
