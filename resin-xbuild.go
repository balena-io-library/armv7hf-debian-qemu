package main

import (
	"log"
	"os"
	"os/exec"
	filepath "path/filepath"
	"syscall"
)

func crossBuildStart() {
	if _, err := os.Stat("/bin/sh.real"); os.IsNotExist(err) {
		err = os.Link("/bin/sh", "/bin/sh.real")
		if err != nil {
			log.Fatal(err)
		}
	}

	err := os.Remove("/bin/sh")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Link("/usr/bin/resin-xbuild", "/bin/sh")
	if err != nil {
		log.Fatal(err)
	}
}

func crossBuildEnd() {
	err := os.Remove("/bin/sh")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Link("/bin/sh.real", "/bin/sh")
	if err != nil {
		log.Fatal(err)
	}
}

func runShell() error {
	var cmd *exec.Cmd

	if _, err := os.Stat("/usr/bin/qemu-arm-static"); err == nil {
		cmd = exec.Command("/usr/bin/qemu-arm-static", append([]string{"-execve", "-0", os.Args[0], "/bin/sh"}, os.Args[1:]...)...)
	} else {
		cmd = exec.Command("/usr/bin/qemu-aarch64-static", append([]string{"-execve", "-0", os.Args[0], "/bin/sh"}, os.Args[1:]...)...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	command := filepath.Base(os.Args[0])

	switch command {
	case "cross-build-start":
		crossBuildStart()
	case "cross-build-end":
		crossBuildEnd()
	default:
		code := 0
		crossBuildEnd()

		if err := runShell(); err != nil {
			code = 1
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					code = status.ExitStatus()
				}
			}
		}

		crossBuildStart()

		os.Exit(code)
	}
}
