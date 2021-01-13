package screen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Screen struct {
	file     *os.File
	lines    []string
	viewport strings.Builder
	rows     int
	cols     int
	CursorX  int
	CursorY  int
}

func (s Screen) String() string {
	s.refresh()
	return s.viewport.String()
}

func New(f *os.File, rows, cols int) Screen {
	s := Screen{}

	s.file = f

	s.rows = rows
	s.cols = cols

	s.CursorX = 1
	s.CursorY = 1

	s.lines = make([]string, s.rows-1)

	scanner := bufio.NewScanner(s.file)

	i := 0

	for scanner.Scan() {
		s.lines[i] = scanner.Text()
		i++
	}

	return s

}

func (s *Screen) refresh() {
	fmt.Fprintf(&s.viewport, "\x1b[?25l") // Hide cursor
	fmt.Fprintf(&s.viewport, "\x1b[H")    // Reset cursor position
	for _, line := range s.lines {
		fmt.Fprintf(&s.viewport, "\x1b[2K") // Clear line
		if len(line) > 0 {
			fmt.Fprintf(&s.viewport, "%v\r\n", line)
		} else {
			fmt.Fprintf(&s.viewport, "~\r\n")
		}
	}

	fmt.Fprintf(&s.viewport, "\x1b[%d;%dH", s.CursorY, s.CursorX) // Move cursor
	fmt.Fprintf(&s.viewport, "\x1b[?25h")                         // Unhide cursor
}
