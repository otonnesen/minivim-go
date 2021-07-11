package screen

import (
	"fmt"
	"minivim/file"
	"os"
	"strings"
)

var refreshes int
var size int

// Screen stores information regarding the viewport and cursor.
type Screen struct {
	file     file.File
	viewport strings.Builder
	rows     int
	cols     int
	cursorX  int
	cursorY  int
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

	s.rows = rows
	s.cols = cols

	s.cursorX = 1
	s.cursorY = 1

	s.file = file.New(f)

	return s

}

func (s *Screen) updateViewport() {
	// TODO: scrolling
	// will have to change for loop since it won't necessarily start
	// at line 1
	refreshes++
	size = len(s.viewport.String())

	fmt.Fprintf(&s.viewport, "\x1b[?25l") // Hide cursor
	fmt.Fprintf(&s.viewport, "\x1b[2J")   // Clear entire screen
	fmt.Fprintf(&s.viewport, "\x1b[H")    // Reset cursor position

	// File contents
	c := s.file.Front
	sep := "\r\n"
	for l := 0; l < s.rows; l++ {
		if l == s.rows-1 {
			sep = ""
		}
		if c != nil {
			fmt.Fprintf(&s.viewport, "%v%v", c.Text, sep)
			c = c.Next
		} else {
			fmt.Fprintf(&s.viewport, "~%v", sep)
		}
	}

	// Exclude the next (debug) line in viewport size calculation
	viewportSize := len(s.viewport.String())

	// Debug
	fmt.Fprintf(&s.viewport, "\x1b[2K")

	fmt.Fprintf(&s.viewport, "screen: %vx%v, ", s.cols, s.rows)
	fmt.Fprintf(&s.viewport, "cursor: (%v,%v), ", s.cursorX, s.cursorY)
	fmt.Fprintf(&s.viewport, "row length: %v, ", s.file.Current.Length())
	fmt.Fprintf(&s.viewport, "num length: %v, ", s.file.NumLines)
	fmt.Fprintf(&s.viewport, "s.viewport size: %v, ", viewportSize)
	fmt.Fprintf(&s.viewport, "size: %v, ", size)
	fmt.Fprintf(&s.viewport, "# refreshes: %v", refreshes)

	fmt.Fprintf(&s.viewport, "\x1b[%d;%dH", s.cursorY, s.cursorX) // Move cursor
	fmt.Fprintf(&s.viewport, "\x1b[?25h")                         // Unhide cursor
}

func (s *Screen) fixBounds() {
	// We fix the bounds on cursorY before cursorX because we use cursorY
	// in the calculation of cursorX's bounds.
	if s.cursorY > s.file.NumLines {
		s.cursorY = s.file.NumLines
		s.file.Current = s.file.Back
	}
	if s.cursorY < 1 {
		s.cursorY = 1
		s.file.Current = s.file.Front
	}

	if s.cursorX < 1 {
		s.cursorX = 1
	}
	if l := s.file.Current.Length(); s.cursorX > l {
		s.cursorX = l
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
	s.file.Current = s.file.Current.Next
	s.fixBounds()
}

// Up moves the cursor up one position, and updates the Screen's state
// accordingly.
func (s *Screen) Up() {
	s.cursorY -= 1
	s.file.Current = s.file.Current.Prev
	s.fixBounds()
}

// Right moves the cursor right one position, and updates the Screen's state
// accordingly.
func (s *Screen) Right() {
	s.cursorX += 1
	s.fixBounds()
}
