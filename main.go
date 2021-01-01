package main

import "C"
import (
	"fmt"
	"log"
	"os"

	"minivim/terminal"
)

type editorConfig struct {
	cx        int
	cy        int
	rowOff    int
	colOff    int
	num_lines int
}

func main() {
	term, err := terminal.Init()
	if err != nil {
		log.Panicf("%v\n", err)
	}

	fmt.Fprintf(os.Stdout, "\x1b[?25l") // Hide cursor
	fmt.Fprintf(os.Stdout, "\x1b[H")    // Reset cursor position

	for i := 0; i < term.Rows; i++ {
		fmt.Fprint(os.Stdout, "\x1b[2K") // Clear line
		fmt.Fprint(os.Stdout, "~")
		fmt.Fprint(os.Stdout, "\r\n")
	}

	fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", 10, 10) // Move cursor
	fmt.Fprintf(os.Stdout, "\x1b[?25h")           // Show cursor

	for {
	}
	terminal.Close()

}
