package main

import "C"
import (
	"bufio"
	"fmt"
	"log"
	"os"

	"minivim/screen"
	"minivim/terminal"
)

type ExitCode int

const (
	Break ExitCode = iota
	Continue
)

type editorConfig struct {
	screen screen.Screen
	term   terminal.Term
}

var E editorConfig

func main() {
	term, err := terminal.New()
	if err != nil {
		log.Panicf("%v\n", err)
	}
	E.term = term

	f, err := os.Open("./test")
	if err != nil {
		log.Panicf("%v\n", err)
	}

	screen := screen.New(f, E.term.Rows, E.term.Cols)
	E.screen = screen

	for {
		fmt.Fprintf(os.Stdout, E.screen.String())
		if processKey() == Break {
			break
		}
	}
	term.Close()

}

// Handles key presses.
func processKey() ExitCode {
	c := readKey()
	switch c {
	case ctrlKey('q'):
		return Break
	case 'h':
		E.screen.Left()
	case 'j':
		E.screen.Down()
	case 'k':
		E.screen.Up()
	case 'l':
		E.screen.Right()
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
