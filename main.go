package main

import "C"
import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"
)

var orig_tios syscall.Termios

func main() {
	enableRawMode()

	rows, cols := getWinsize()
	cols = cols

	fmt.Fprintf(os.Stdout, "\x1b[?25l") // Hide cursor
	fmt.Fprintf(os.Stdout, "\x1b[H")    // Reset cursor position

	for i := 0; i < rows; i++ {
		fmt.Fprint(os.Stdout, "\x1b[2K") // Clear line
		fmt.Fprint(os.Stdout, "~")
		fmt.Fprint(os.Stdout, "\r\n")
	}

	fmt.Fprintf(os.Stdout, "\x1b[%d;%dH", 10, 10) // Move cursor
	fmt.Fprintf(os.Stdout, "\x1b[?25h")           // Show cursor

	for {
	}
	disableRawMode()

}

func getWinsize() (int, int) {
	// ioctl winsize struct definition
	// https://refspecs.linuxfoundation.org/LSB_3.0.0/LSB-Core-generic/LSB-Core-generic/libc-ddefs.html
	ws := struct {
		ws_row    uint16
		ws_col    uint16
		ws_xpixel uint16
		ws_ypixel uint16
	}{}

	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&ws)))
	if e != 0 {
		log.Fatalf("Syscall error: %v\n", e)
	}

	return int(ws.ws_row), int(ws.ws_col)
}

func disableRawMode() {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&orig_tios)))
	if r != 0 {
		log.Fatalf("Syscall error: %v\n", e)
	}
}

func enableRawMode() {
	r, _, e := syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TCGETS, uintptr(unsafe.Pointer(&orig_tios)))
	if r != 0 {
		log.Fatalf("Syscall error: %v\n", e)
	}

	// Enable raw mode
	raw := orig_tios

	raw.Lflag &^= syscall.ECHO | syscall.ICANON |
		syscall.IEXTEN | syscall.ISIG
	raw.Iflag &^= syscall.IXON | syscall.ICRNL

	raw.Oflag &^= syscall.OPOST

	raw.Iflag &^= syscall.BRKINT | syscall.INPCK |
		syscall.ISTRIP

	raw.Cflag &^= syscall.CS8

	r, _, e = syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), syscall.TCSETS, uintptr(unsafe.Pointer(&raw)))
	if r != 0 {
		log.Fatalf("Syscall error: %v\n", e)
	}
}
