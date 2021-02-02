package screen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var refreshes int
var size int

type Screen struct {
	file     *os.File
	lines    []string
	viewport strings.Builder
	rows     int
	cols     int
	CursorX  int
	CursorY  int
	numLines int
}

func (s Screen) String() string {
	s.updateViewport()
	return s.viewport.String()
}

func New(f *os.File, rows, cols int) Screen {
	s := Screen{}

	s.file = f

	s.rows = rows
	s.cols = cols

	s.CursorX = 1
	s.CursorY = 1

	// We allocate one less line than the terminal has to make room for
	// the debug info.
	s.lines = make([]string, s.rows-1)

	scanner := bufio.NewScanner(s.file)

	s.numLines = 0

	for scanner.Scan() {
		s.lines[s.numLines] = scanner.Text()
		s.numLines++
	}

	return s

}

func (s *Screen) updateViewport() {
	// s.viewport = strings.Builder{}
	refreshes++
	size = len(s.viewport.String())

	fmt.Fprintf(&s.viewport, "\x1b[?25l") // Hide cursor
	fmt.Fprintf(&s.viewport, "\x1b[H")    // Reset cursor position

	// File contents
	for _, line := range s.lines {
		fmt.Fprintf(&s.viewport, "\x1b[2K") // Clear line
		if len(line) > 0 {
			fmt.Fprintf(&s.viewport, "%v\r\n", line)
		} else {
			fmt.Fprintf(&s.viewport, "~\r\n")
		}
	}

	// Exclude the next (debug) line in viewport size calculation
	viewportSize := len(s.viewport.String())

	// Debug
	fmt.Fprintf(&s.viewport, "\x1b[2K")

	fmt.Fprintf(&s.viewport, "screen: %vx%v, ", s.cols, s.rows)
	fmt.Fprintf(&s.viewport, "cursor: (%v,%v), ", s.CursorX, s.CursorY)
	fmt.Fprintf(&s.viewport, "row length: %v, ", len([]rune(s.lines[s.CursorY-1])))
	fmt.Fprintf(&s.viewport, "num length: %v, ", s.numLines)
	fmt.Fprintf(&s.viewport, "s.viewport size: %v, ", viewportSize)
	fmt.Fprintf(&s.viewport, "size: %v, ", size)
	fmt.Fprintf(&s.viewport, "# refreshes: %v", refreshes)

	fmt.Fprintf(&s.viewport, "\x1b[%d;%dH", s.CursorY, s.CursorX) // Move cursor
	fmt.Fprintf(&s.viewport, "\x1b[?25h")                         // Unhide cursor
}

func (s *Screen) fixBounds() {
	// We fix the bounds on CursorY before CursorX because we use CursorY
	// in the calculation of CursorX's bounds.
	if s.CursorY > s.numLines {
		s.CursorY = s.numLines
	}
	if s.CursorY < 1 {
		s.CursorY = 1
	}

	if s.CursorX < 1 {
		s.CursorX = 1
	}
	// We convert the string into a rune array so characters using
	// multiple code points only count as one character.
	if s.CursorX > len([]rune(s.lines[s.CursorY-1])) {
		s.CursorX = len([]rune(s.lines[s.CursorY-1]))
	}
}

func (s *Screen) Left() {
	s.CursorX -= 1
	s.fixBounds()
}

func (s *Screen) Down() {
	s.CursorY += 1
	s.fixBounds()
}

func (s *Screen) Up() {
	s.CursorY -= 1
	s.fixBounds()
}

func (s *Screen) Right() {
	s.CursorX += 1
	s.fixBounds()
}
