// Try to keep the gross stuff in here

package terminal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// Term stores terminal information, including the terminal's original config
// to be restored after exiting.
type Term struct {
	Rows     int
	Cols     int
	origTios syscall.Termios
}

// New sets the terminal to non-canonical mode and allocates a new Term.
func New() (Term, error) {
	var config Term

	origTios, err := enableRawMode()
	if err != nil {
		return Term{}, err
	}
	config.origTios = origTios

	rows, cols := getWinsize()
	config.Rows = rows
	config.Cols = cols

	return config, err
}

// Close returns the terminal to canonical mode and clears the screen.
// TODO: save and restore contents of screen prior to opening
func (config Term) Close() error {
	fmt.Fprint(os.Stdout, "\x1b[2J") // Clear entire screen
	fmt.Fprint(os.Stdout, "\x1b[H")  // Reset cursor position
	return disableRawMode(config)
}

// Returns terminal window size as (rows, cols)
func getWinsize() (int, int) {
	// ioctl winsize struct definition
	// https://refspecs.linuxfoundation.org/LSB_3.0.0/LSB-Core-generic/LSB-Core-generic/libc-ddefs.html
	ws := struct {
		ws_row    uint16
		ws_col    uint16
		ws_xpixel uint16
		ws_ypixel uint16
	}{}

	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(),
		syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
	if e != 0 {
		log.Fatalf("Syscall error: %v\n", e)
	}

	return int(ws.ws_row), int(ws.ws_col)
}

// Resets terminal to default
func disableRawMode(config Term) error {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(),
		syscall.TCSETS, uintptr(unsafe.Pointer(&config.origTios)))
	if r != 0 {
		return errors.New(fmt.Sprintf("%v", e))
	}

	return nil
}

// Sets terminal to non-canonical mode
func enableRawMode() (syscall.Termios, error) {
	var origTios syscall.Termios
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(),
		syscall.TCGETS, uintptr(unsafe.Pointer(&origTios)))
	if r != 0 {
		return syscall.Termios{}, errors.New(fmt.Sprintf("%v", e))
	}

	// Enable raw mode
	raw := origTios

	raw.Lflag &^= syscall.ECHO | syscall.ICANON |
		syscall.IEXTEN | syscall.ISIG
	raw.Iflag &^= syscall.IXON | syscall.ICRNL

	raw.Oflag &^= syscall.OPOST

	raw.Iflag &^= syscall.BRKINT | syscall.INPCK |
		syscall.ISTRIP

	raw.Cflag &^= syscall.CS8

	r, _, e = syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(),
		syscall.TCSETS, uintptr(unsafe.Pointer(&raw)))
	if r != 0 {
		return syscall.Termios{}, errors.New(fmt.Sprintf("%v", e))
	}

	return origTios, nil
}
