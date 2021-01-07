package main

import "C"
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"minivim/terminal"
)

type ExitCode int

const (
	Break ExitCode = iota
	Continue
)

type editorConfig struct {
	cx        int
	cy        int
	rowOff    int
	colOff    int
	num_lines int
	term      terminal.TermConfig
}

var E editorConfig

func main() {
	term, err := terminal.Init()
	if err != nil {
		log.Panicf("%v\n", err)
	}
	E.term = term

	for {
		refreshScreen()
		if processKey() == Break {
			break
		}
	}
	term.Close()

}

// Re-draws the screen to reflect changes in the editor's internal state.
func refreshScreen() {
	fmt.Fprintf(os.Stdout, "\x1b[?25l") // Hide cursor
	fmt.Fprintf(os.Stdout, "\x1b[H")    // Reset cursor position

	fmt.Fprintf(os.Stdout, draw())

	fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", E.cx, E.cy) // Move cursor
	fmt.Fprintf(os.Stdout, "\x1b[?25h")               // Show cursor
}

// Handles key presses.
func processKey() ExitCode {
	c := readKey()
	switch c {
	case ctrlKey('q'):
		return Break
	}
	return Continue
}

// Returns the byte sequence corresponding to pressing Ctrl + `c`.
func ctrlKey(c byte) byte {
	return c & 0x1f
}

// Get next byte of input from stdin
func readKey() byte {
	in := bufio.NewReader(os.Stdin)
	c, err := in.ReadByte()
	if err != nil {
		log.Panicf("%v\n", err)
	}
	return c
}

// Returns file contents as a string, including escape sequences for
// terminal output.
func draw() string {
	var b strings.Builder
	for i := 0; i < E.term.Rows; i++ {
		fmt.Fprint(&b, "\x1b[2K") // Clear line
		fmt.Fprint(&b, "~")
		fmt.Fprint(&b, "\r\n")
	}
	return b.String()
}
