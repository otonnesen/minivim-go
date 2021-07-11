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

	fmt.Fprintf(&s.viewport, "\x1b[H")             // Reset cursor position
	fmt.Fprintf(&s.viewport, "\x1b[%dB", s.rows-1) // Move cursor to bottom line
	fmt.Fprintf(&s.viewport, s.printDebugLine())

	var x int
	if l := len(s.file.Current.Text); s.cursorX < l {
		x = s.cursorX
	} else {
		x = l
	}
	fmt.Fprintf(&s.viewport, "\x1b[%d;%dH", s.cursorY, x) // Move cursor
	fmt.Fprintf(&s.viewport, "\x1b[?25h")                 // Unhide cursor

}

func (s Screen) printDebugLine() string {
	var ret strings.Builder

	viewportSize := len(s.viewport.String())

	fmt.Fprintf(&ret, "\x1b[2K")
	fmt.Fprintf(&ret, "screen: %vx%v, ", s.cols, s.rows)
	fmt.Fprintf(&ret, "cursor: (%v,%v), ", s.cursorX, s.cursorY)
	fmt.Fprintf(&ret, "row length: %v, ", s.file.Current.Length())
	fmt.Fprintf(&ret, "num length: %v, ", s.file.NumLines)
	fmt.Fprintf(&ret, "s.viewport size: %v, ", viewportSize)
	fmt.Fprintf(&ret, "size: %v, ", size)
	fmt.Fprintf(&ret, "# refreshes: %v", refreshes)

	return ret.String()
}

// fixRowBounds checks if s.cursorY is outside the allowed bounds (1 <=
// s.cursorY <= # lines in file) and sets it to a valid value if it isn't.
func (s *Screen) fixRowBounds() {
	if s.cursorY > s.file.NumLines {
		s.cursorY = s.file.NumLines
		s.file.Current = s.file.Back
	}
	if s.cursorY < 1 {
		s.cursorY = 1
		s.file.Current = s.file.Front
	}
}

// fixColBounds checks if s.cursorX is outside the allowed boudns (1 <=
// s.CursorX <= current line length) and sets it to a valid value if it isn't.
func (s *Screen) fixColBounds() {
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
	// We call fixColBounds beforehand to fix s.cursorX in the case that
	// it's larger than the length of the current line.
	s.fixColBounds()
	s.cursorX -= 1
	s.fixColBounds()
}

// Down moves the cursor down one position, and updates the Screen's state
// accordingly.
func (s *Screen) Down() {
	s.cursorY += 1
	s.file.Current = s.file.Current.Next
	s.fixRowBounds()
}

// Up moves the cursor up one position, and updates the Screen's state
// accordingly.
func (s *Screen) Up() {
	s.cursorY -= 1
	s.file.Current = s.file.Current.Prev
	s.fixRowBounds()
}

// Right moves the cursor right one position, and updates the Screen's state
// accordingly.
func (s *Screen) Right() {
	// We skip this if cursorX is at the end of the line so we don't
	// overwrite the saved cursorX value unnecessarily
	if s.cursorX < s.file.Current.Length() {
		s.cursorX += 1
		s.fixColBounds()
	}
}
