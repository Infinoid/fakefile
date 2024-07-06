package main

import (
	"bufio"

	"github.com/infinoid/fakefile"
)

// write a couple lines of text to the fake file
func writefile(ff *fakefile.Fakefile) {
	w := ff.Writer()
	defer w.Close()
	_, _ = w.Write([]byte("Hello world!\n"))
	_, _ = w.Write([]byte("So happy!\n"))
}

// read lines of text using bufio.Scanner
func readfile(ff *fakefile.Fakefile) []string {
	r := ff.Reader()
	defer r.Close()
	s := bufio.NewScanner(r)

	rv := make([]string, 0, 2)
	for s.Scan() {
		rv = append(rv, s.Text())
	}
	return rv
}

func main() {
	ff := fakefile.New()

	writefile(ff)

	lines := readfile(ff)
	for _, line := range lines {
		println(line)
	}
}
