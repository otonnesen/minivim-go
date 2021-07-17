package screen

import (
	"fmt"
	"log"
	"minivim/file"
	"os"
	"strings"
)

// Screen stores information regarding the viewport and cursor.
type Screen struct {
	// STATE
	file file.FileContents
	// The first line of the buffer
	bufferTop *file.Line
	// The last line of the buffer
	bufferBottom *file.Line
	viewport     strings.Builder
	rows         int
	cols         int
	cursorX      int
	cursorY      int

	// DEBUG
	fileLogger *log.Logger

	// FLAGS
	Number bool
}

// String returns the viewport in its current state, including terminal
// escape sequences.
func (s Screen) String() string {
	s.updateViewport()
	return s.viewport.String()
}

// New allocates a new Screen with viewport size rows x cols and displaying
// the contents of f.
func New(f *os.File, l *log.Logger, rows, cols int) Screen {
	s := Screen{}

	s.rows = rows
	s.cols = cols

	s.cursorX = 1
	s.cursorY = 1

	s.file = file.New(f)

	s.bufferTop = s.file.Front
	s.bufferBottom = s.bufferTop
	for i := 0; i < s.rows; i++ {
		if s.bufferBottom.Next == nil {
			break
		}
		s.bufferBottom = s.bufferBottom.Next
	}

	s.fileLogger = l

	s.Number = true

	return s

}

func (s *Screen) updateViewport() {

	fmt.Fprintf(&s.viewport, "\x1b[?25l") // Hide cursor
	fmt.Fprintf(&s.viewport, "\x1b[2J")   // Clear entire screen
	fmt.Fprintf(&s.viewport, "\x1b[H")    // Reset cursor position

	// File contents
	c := s.bufferTop
	sep := "\r\n"
	lineNumPadding := len(fmt.Sprintf("%d", s.file.NumLines))
	lineNumePrefix := ""
	for l := 0; l < s.rows; l++ {
		if l == s.rows-1 {
			sep = ""
		}
		if c != nil {
			if s.Number {
				lineNumePrefix = fmt.Sprintf("%*d ", lineNumPadding, c.Number)
			}
			fmt.Fprintf(&s.viewport, "%v%v%v", lineNumePrefix, c.Text, sep)
			c = c.Next
		} else {
			fmt.Fprintf(&s.viewport, "~%v", sep)
		}
	}

	fmt.Fprintf(&s.viewport, "\x1b[H")             // Reset cursor position
	fmt.Fprintf(&s.viewport, "\x1b[%dB", s.rows-1) // Move cursor to bottom line
	fmt.Fprintf(&s.viewport, s.debugLine())

	var x int
	if l := s.file.Current.Length(); s.cursorX <= l {
		x = s.cursorX
	} else if s.file.Current.Length() == 0 {
		x = 1
	} else {
		x = l
	}
	x += s.lineNumPadding()
	fmt.Fprintf(&s.viewport, "\x1b[%d;%dH", s.cursorY, x) // Move cursor
	fmt.Fprintf(&s.viewport, "\x1b[?25h")                 // Unhide cursor

}

// lineNumPadding returns the number of characters each line is padded by.
// If the Number flag is set, then the line number is prepended to each line of
// the file's contents. For example, if our file is 100 lines long, then
// lineNumPadding will return 4, since each line will look like:
//   1 This is line 1
//   2 This is line 2
// ...
//  99 This is line 99
// 100 This is line 100
func (s Screen) lineNumPadding() int {
	l := 0
	if s.Number {
		l = 1 + len(fmt.Sprintf("%d", s.file.NumLines))
	}
	return l
}

func (s Screen) debugLine() string {
	var ret strings.Builder
	var log strings.Builder

	viewportSize := len(s.viewport.String())

	fmt.Fprintf(&ret, "\x1b[2K")

	fmt.Fprintf(&log, "screen: %vx%v, ", s.cols, s.rows)
	fmt.Fprintf(&log, "cursor: (%v,%v), ", s.cursorX, s.cursorY)
	fmt.Fprintf(&log, "row length: %v, ", s.file.Current.Length())
	fmt.Fprintf(&log, "line number: %v, ", s.file.Current.Number)
	fmt.Fprintf(&log, "bufferTop: %v, ", s.bufferTop.Number)
	fmt.Fprintf(&log, "bufferBottom: %v, ", s.bufferBottom.Number)
	fmt.Fprintf(&log, "lineNumPadding: %v, ", s.lineNumPadding())
	fmt.Fprintf(&log, "s.viewport size: %v ", viewportSize)

	fmt.Fprintf(&ret, "%v", log.String())

	s.fileLogger.Printf("%v", log.String())

	return ret.String()
}

// ScrollDown scrolls l rows down.
func (s *Screen) ScrollDown(l int) {
	for i := 0; i < l; i++ {
		if s.bufferTop.Next == nil {
			break
		}
		s.bufferTop = s.bufferTop.Next
		if s.bufferBottom.Next != nil {
			s.bufferBottom = s.bufferBottom.Next
		}
		s.cursorY--
		s.fixRowBounds()
	}
}

// ScrollUp scrolls l rows up.
func (s *Screen) ScrollUp(l int) {
	for i := 0; i < l; i++ {
		if s.bufferTop.Prev == nil {
			break
		}
		s.bufferTop = s.bufferTop.Prev
		if s.bufferTop.Number+s.rows-1 <= s.file.NumLines {
			s.bufferBottom = s.bufferBottom.Prev
		}
		s.cursorY++
		s.fixRowBounds()
	}
}

// fixRowBounds checks if s.cursorY is outside the allowed bounds (1 <=
// s.cursorY <= # lines in file) and sets it to a valid value if it isn't.
func (s *Screen) fixRowBounds() {
	if s.cursorY > s.rows-1 {
		s.cursorY = s.rows - 1
		s.file.Current = s.bufferBottom
	}
	if s.cursorY < 1 {
		s.cursorY = 1
		s.file.Current = s.bufferTop
	}
}

// fixColBounds checks if s.cursorX is outside the allowed bounds (1 <=
// s.CursorX <= current line length) and sets it to a valid value if it isn't.
func (s *Screen) fixColBounds() {
	if s.cursorX < 1 {
		s.cursorX = 1
	}
	if l := s.file.Current.Length(); s.cursorX > l {
		s.cursorX = l
	}
	if s.file.Current.Length() == 0 {
		s.cursorX = 1
	}
}

// Left moves the cursor left one position, and updates the Screen's state
// accordingly.
func (s *Screen) Left() {
	// We call fixColBounds beforehand to fix s.cursorX in the case that
	// it's larger than the length of the current line.
	s.fixColBounds()
	s.cursorX--
	s.fixColBounds()
}

// Down moves the cursor down one position, and updates the Screen's state
// accordingly.
func (s *Screen) Down() {
	// If s.Current is the last line in the file, we do nothing.
	if s.file.Current.Next == nil {
		return
	}
	// If s.cursorY is s.rows-1 (due to the debug line), then we're on the
	// last line of the buffer and need to scroll up.
	if s.cursorY == s.rows-1 {
		s.ScrollDown(1)
	}
	s.cursorY++
	if s.file.Current.Next != nil {
		s.file.Current = s.file.Current.Next
	}
	s.fixRowBounds()
}

// Up moves the cursor up one position, and updates the Screen's state
// accordingly.
func (s *Screen) Up() {
	// If s.cursorY is 1, then we're on the first line of the buffer and
	// need to scroll up.
	if s.cursorY == 1 {
		s.ScrollUp(1)
	}
	s.cursorY--
	s.file.Current = s.file.Current.Prev
	s.fixRowBounds()

}

// Right moves the cursor right one position, and updates the Screen's state
// accordingly.
// TODO: handle lines longer than the screen width.
func (s *Screen) Right() {
	// We skip this if cursorX is at the end of the line so we don't
	// overwrite the saved cursorX value unnecessarily.
	if s.cursorX < s.file.Current.Length() {
		s.cursorX++
		s.fixColBounds()
	}
}
