package screen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var refreshes int
var size int

// Screen stores information regarding the viewport and cursor.
type Screen struct {
	file     *os.File
	lines    []string
	viewport strings.Builder
	rows     int
	cols     int
	cursorX  int
	cursorY  int
	numLines int
}

// String returns the viewport in its current state, including terminal
// escape sequences.
func (s Screen) String() string {
	s.updateViewport()
	return s.viewport.String()
}

// New allocates a new Screen with viewport size rows x cols and displaying
// the contents of f.
func New(f *os.File, rows, cols int) Screen {
	s := Screen{}

	s.file = f

	s.rows = rows
	s.cols = cols

	s.cursorX = 1
	s.cursorY = 1

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
	// TODO: scrolling
	// will have to change for loop since it won't necessarily start at
	// line 1
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
	fmt.Fprintf(&s.viewport, "cursor: (%v,%v), ", s.cursorX, s.cursorY)
	fmt.Fprintf(&s.viewport, "row length: %v, ", len([]rune(s.lines[s.cursorY-1])))
	fmt.Fprintf(&s.viewport, "num length: %v, ", s.numLines)
	fmt.Fprintf(&s.viewport, "s.viewport size: %v, ", viewportSize)
	fmt.Fprintf(&s.viewport, "size: %v, ", size)
	fmt.Fprintf(&s.viewport, "# refreshes: %v", refreshes)

	fmt.Fprintf(&s.viewport, "\x1b[%d;%dH", s.cursorY, s.cursorX) // Move cursor
	fmt.Fprintf(&s.viewport, "\x1b[?25h")                         // Unhide cursor
}

func (s *Screen) fixBounds() {
	// We fix the bounds on cursorY before cursorX because we use cursorY
	// in the calculation of cursorX's bounds.
	if s.cursorY > s.numLines {
		s.cursorY = s.numLines
	}
	if s.cursorY < 1 {
		s.cursorY = 1
	}

	if s.cursorX < 1 {
		s.cursorX = 1
	}
	// We convert the string into a rune array so characters using
	// multiple code points only count as one character.
	if s.cursorX > len([]rune(s.lines[s.cursorY-1])) {
		s.cursorX = len([]rune(s.lines[s.cursorY-1]))
	}
}

// Left moves the cursor left one position, and updates the Screen's state
// accordingly.
func (s *Screen) Left() {
	s.cursorX -= 1
	s.fixBounds()
}

// Down moves the cursor down one position, and updates the Screen's state
// accordingly.
func (s *Screen) Down() {
	s.cursorY += 1
	s.fixBounds()
}

// Up moves the cursor up one position, and updates the Screen's state
// accordingly.
func (s *Screen) Up() {
	s.cursorY -= 1
	s.fixBounds()
}

// Right moves the cursor right one position, and updates the Screen's state
// accordingly.
func (s *Screen) Right() {
	s.cursorX += 1
	s.fixBounds()
}
