package file

import (
	"bufio"
	"os"
)

// Line stores information related to a single line of a file's
// contents.
type Line struct {
	Next   *Line
	Prev   *Line
	Text   string
	Number int
}

// Length returns the length of l.Text. We convert the string into a
// rune array so characters using multiple code points only count as one
// character.
func (l Line) Length() int {
	return len([]rune(l.Text))
}

// FileContents stores information related to the currently open file
type FileContents struct {
	Front    *Line
	Back     *Line
	Current  *Line
	NumLines int
}

// New allocates a new FileContents storing the contents of f.
func New(f *os.File) FileContents {
	file := FileContents{}

	file.Front = &Line{}
	file.Current = file.Front

	c := file.Front

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		file.NumLines++
		c.Text = scanner.Text()
		c.Number = file.NumLines
		c.Next = &Line{}
		c.Next.Prev = c
		c = c.Next
	}

	file.Back = c.Prev

	// Fencepost
	file.Back.Next = nil

	return file
}
