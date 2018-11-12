package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

var (
	nsfs string
	cmd  string
)

func initFlags() {
	flag.StringVar(&nsfs, "nsfs", "", "specific nstype file /proc/PID/ns/nstype")
	flag.StringVar(&cmd, "cmd", "/bin/bash", "specific command")
	flag.Parse()
}

func main() {
	initFlags()

	if err := validate(); err != nil {
		oops(err)
	}

	file, err := os.Open(nsfs)
	if err != nil {
		oops(err)
	}
	defer file.Close()

	// setns (attach to the existing namespace)
	if _, _, errno := syscall.RawSyscall(unix.SYS_SETNS, uintptr(file.Fd()), 0, 0); errno != 0 {
		oops(errno)
	}

	c := exec.Command(cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		oops(err)
	}
}

func validate() error {
	if nsfs == "" {
		return fmt.Errorf("failed to attach empty nstype file")
	}
	return nil
}

func oops(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
